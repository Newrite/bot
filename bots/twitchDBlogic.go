package bots

import (
	"bot/resource"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

func (bt *BotTwitch) createChannelSettingsTable() {
	_, err := resource.SingleDB().Exec(`
create table if not exists TWITCH_CHANNEL_CONFIG
(
    CHANNEL_NAME    TEXT    not null primary key,
    STATUS          INTEGER not null,
    REACT_STATUS    INTEGER not null,
    CMD_STATUS      INTEGER not null,
    REACT_RATE      INTEGER not null,
    LAST_REACT_TIME INTEGER not null,
    IS_MODERATOR    INTEGER not null
);`)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "Exec",
			"error":    err,
		}).Errorln("Ошибка создания таблицы в бд.")
	}
}

func (bt *BotTwitch) createChannelCMDTable(channel string) {
	_, err := resource.SingleDB().Exec(`
create table if not exists ` + channel + `_CMD_LIST
(
    COMMAND    TEXT    not null,
    ANSWER          TEXT not null
);`)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "Exec",
			"error":    err,
		}).Errorln("Ошибка создания таблицы команд в бд.")
	}
}

func (bt *BotTwitch) checkChannelCMDinTable(channel, cmd string) bool {
	rows, err := resource.SingleDB().Query(`SELECT COMMAND FROM ` + channel + `_CMD_LIST`)
	if rows == nil || err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "Query",
			"error":    err,
		}).Errorln("Ошибка запроса команды с бд.")
	}
	defer rows.Close()
	var command string
	for rows.Next() {
		err = rows.Scan(&command)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "bots",
				"function": "Scan",
				"error":    err,
			}).Errorln("Ошибка скан запроса.")
		}
		if command == cmd {
			return true
		}
	}
	return false
}

func (bt *BotTwitch) channelCMDfromTable(channel, cmd string) (string, error) {
	rows, err := resource.SingleDB().Query(`SELECT ANSWER FROM `+channel+`_CMD_LIST WHERE COMMAND = $1`, cmd)
	if rows == nil || err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "Query",
			"error":    err,
		}).Errorln("Ошибка запроса команды с бд.")
		return "", err
	}
	defer rows.Close()
	var answer string
	for rows.Next() {
		err = rows.Scan(&answer)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "bots",
				"function": "Scan",
				"error":    err,
			}).Errorln("Ошибка скан запроса.")
			return "", err
		}
		return answer, nil
	}
	return "", nil
}

func (bt *BotTwitch) CMDlistFromChannel(channel string) string {
	rows, err := resource.SingleDB().Query(`SELECT COMMAND FROM ` + channel + `_CMD_LIST`)
	if rows == nil || err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "Query",
			"error":    err,
		}).Errorln("Ошибка запроса команды с бд.")
		return err.Error()
	}
	defer rows.Close()
	var answer string
	var a string
	for rows.Next() {
		err = rows.Scan(&a)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "bots",
				"function": "Scan",
				"error":    err,
			}).Errorln("Ошибка скан запроса.")
			return err.Error()
		}
		answer += a + ", "
	}
	return strings.TrimSuffix(answer, `, `)
}

func (bt *BotTwitch) checkChannelCMDinList(channel, cmd string) bool {
	for _, list := range CMDList {
		for _, platform := range list.Platform {
			if platform == GG || platform == TW || platform == "all" {
				for _, ch := range list.Channels {
					if ch == channel || ch == "all" {
						for _, cc := range list.Command {
							if cc == cmd {
								return true
							}
						}
					}
				}
			}
		}
	}
	return false
}

func (bt *BotTwitch) addChannelCMD(command, answer, channel string) string {
	result, err := resource.SingleDB().Exec(`
			INSERT INTO `+channel+`_CMD_LIST
				(COMMAND, ANSWER)
			VALUES
				($1, $2)`,
		command, answer)
	if err != nil || result == nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "Exec",
			"error":    err,
		}).Errorln("Ошибка добавления значения бд.")
		return err.Error()
	}
	return "Создана команда: " + command + ". С ответом: " + answer + ". На канале: " + channel
}

func (bt *BotTwitch) updateChannelCMD(command, answer, channel string) string {
	_, err := resource.SingleDB().Exec(`
		UPDATE `+channel+`_CMD_LIST
		SET 
			ANSWER = $1
		WHERE COMMAND = $2`,
		answer,
		command)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "Exec",
			"error":    err,
		}).Errorln("Ошибка изменение данных в бд.")
		return err.Error()
	}
	return "Обновлена команда: " + command + ". Новый ответ: " + answer + ". На канале: " + channel
}

func (bt *BotTwitch) deleteCMDfromChannel(command, channel string) string {
	_, err := resource.SingleDB().Exec(`
		DELETE FROM `+channel+`_CMD_LIST
		WHERE COMMAND = $1`,
		command)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "Exec",
			"error":    err,
		}).Errorln("Ошибка изменение данных в бд.")
		return err.Error()
	}
	return "Удалена команда: " + command + ". На канале: " + channel
}

func (bt *BotTwitch) initChannelSettings() {
	bt.createChannelSettingsTable()
	for _, channel := range bt.Channels {
		bt.Settings[channel] = &botSettings{
			Status:        true,
			ReactStatus:   true,
			CMDStatus:     true,
			ReactRate:     30,
			LastReactTime: time.Now().Unix(),
			IsModerator:   false,
		}
		bt.createChannelCMDTable(channel)
		result, err := resource.SingleDB().Exec(`
			INSERT OR IGNORE INTO TWITCH_CHANNEL_CONFIG
				(CHANNEL_NAME, STATUS, REACT_STATUS, CMD_STATUS, REACT_RATE, LAST_REACT_TIME, IS_MODERATOR)
			VALUES
				($1, $2, $3, $4, $5, $6, $7)`,
			channel, bt.Settings[channel].Status, bt.Settings[channel].ReactStatus, bt.Settings[channel].CMDStatus,
			bt.Settings[channel].ReactRate, bt.Settings[channel].LastReactTime, bt.Settings[channel].IsModerator)
		if err != nil || result == nil {
			log.WithFields(log.Fields{
				"package":  "bots",
				"function": "Exec",
				"error":    err,
			}).Errorln("Ошибка добавления значения бд.")
		}
	}
	bt.loadChannelSettings()
}

func (bt *BotTwitch) loadChannelSettings() {
	for _, channel := range bt.Channels {
		rows, err := resource.SingleDB().Query(`SELECT * FROM TWITCH_CHANNEL_CONFIG WHERE CHANNEL_NAME = $1`, channel)
		if rows == nil || err != nil {
			log.WithFields(log.Fields{
				"package":  "bots",
				"function": "Query",
				"error":    err,
			}).Errorln("Ошибка запроса с бд.")
		}
		for rows.Next() {
			err = rows.Scan(&channel, &bt.Settings[channel].Status, &bt.Settings[channel].ReactStatus,
				&bt.Settings[channel].CMDStatus, &bt.Settings[channel].ReactRate,
				&bt.Settings[channel].LastReactTime, &bt.Settings[channel].IsModerator)
			if err != nil {
				log.WithFields(log.Fields{
					"package":  "bots",
					"function": "Scan",
					"error":    err,
				}).Errorln("Ошибка скан запроса.")
			}
			err = rows.Close()
			if err != nil {
				log.WithFields(log.Fields{
					"package":  "bots",
					"function": "Close",
					"error":    err,
				}).Errorln("Ошибка закрытия rows.")
			}
		}
		if bt.handleApiRequest("", channel, "", "!evaismod") == "true" {
			bt.Settings[channel].IsModerator = true
		} else {
			bt.Settings[channel].IsModerator = false
		}
		bt.saveChannelSettings(channel)
	}
}

func (bt *BotTwitch) saveChannelSettings(channel string) {
	_, err := resource.SingleDB().Exec(`
		UPDATE TWITCH_CHANNEL_CONFIG
		SET 
			STATUS = $1,
			REACT_STATUS = $2,
			CMD_STATUS = $3,
			REACT_RATE = $4,
			LAST_REACT_TIME = $5,
			IS_MODERATOR = $6
		WHERE CHANNEL_NAME = $7`,
		bt.Settings[channel].Status,
		bt.Settings[channel].ReactStatus,
		bt.Settings[channel].CMDStatus,
		bt.Settings[channel].ReactRate,
		bt.Settings[channel].LastReactTime,
		bt.Settings[channel].IsModerator,
		channel)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "Exec",
			"error":    err,
		}).Errorln("Ошибка изменение данных в бд.")
	}
}
