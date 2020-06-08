package bots

import (
	"bot/resource"
	log "github.com/sirupsen/logrus"
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

func (bt *BotTwitch) initChannelSettings() {
	for _, channel := range bt.Channels {
		bt.Settings[channel] = &botSettings{
			Status:        true,
			ReactStatus:   true,
			CMDStatus:     true,
			ReactRate:     30,
			LastReactTime: time.Now().Unix(),
			IsModerator:   false,
		}
		bt.createChannelSettingsTable()
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
