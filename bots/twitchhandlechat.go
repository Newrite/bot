package bots

import (
	"bot/resource"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

func (bt *BotTwitch) handleChat() error {
	line, err := bt.ReadChannels.ReadLine()
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "ReadChannels.ReadLine",
			"file":     "twitchhandlechat.go",
			"body":     "handleChat",
			"error":    err,
		}).Errorln("Ошибка чтения строки.")
		return err
	}
	if line == `PING :tmi.twitch.tv` {
		fmt.Println("PING :tmi.twitch.tv in handleChat")
		_, err := bt.Connection.Write([]byte("PONG\r\n"))
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "bots",
				"function": "Connection.Write",
				"file":     "twitchhandlechat.go",
				"body":     "handleChat",
				"error":    err,
			}).Errorln("Ошибка отправки сообщения.")
			return err
		}
		return nil
	}
	var userName, channel, message, rewardID string = bt.handleLine(line)
	if rewardID != "" {
		go bt.handleRewards(message, userName, channel, rewardID)
	}
	if time.Since(time.Unix(bt.uptime, 0)) > 7*time.Second {
		bt.writeLog(userName, channel, message)
	}
	if strings.Contains(message, "!Ada, ") && (userName == channel || userName == bt.OwnerBot) {
		go bt.handleMasterCmd(message, channel)
		return nil
	}
	if _, ok := bt.Settings[channel]; ok {
		if !bt.Settings[channel].Status {
			return nil
		}
	}
	bt.checkReact(channel, message)
	if strings.Contains(message, "револьвер выстреливает!") && userName == "moobot" {
		go bt.say(bt.resurrected(message, channel), channel)
	}
	lowMessage := strings.ToLower(message)
	if strings.HasPrefix(lowMessage, TwPrefix) {
		if _, ok := bt.Settings[channel]; ok {
			if bt.Settings[channel].CMDStatus {
				msgSl := strings.Fields(lowMessage)
				go bt.say(checkCMD(userName, channel, msgSl[0], TW, lowMessage, message, ""), channel)
			}
		}
	}
	if strings.HasPrefix(lowMessage, TwPrefix) {
		if _, ok := bt.Settings[channel]; ok {
			if bt.Settings[channel].CMDStatus {
				msgSl := strings.Fields(lowMessage)
				go bt.say(handleCMDfromDB(userName, channel, strings.TrimPrefix(msgSl[0], TwPrefix)), channel)
			}
		}
	}
	time.Sleep(10 * time.Millisecond)
	return nil
}

func (bt *BotTwitch) writeLog(userName, channel, message string) {
	if message != "" &&
		!strings.Contains(userName, `tmi.twitch.tv 353`) &&
		!strings.Contains(userName, `tmi.twitch.tv 366`) &&
		!strings.Contains(userName, `tmi.twitch.tv 001`) &&
		!strings.Contains(userName, `tmi.twitch.tv 002`) &&
		!strings.Contains(userName, `tmi.twitch.tv 003`) &&
		!strings.Contains(userName, `tmi.twitch.tv 004`) &&
		!strings.Contains(userName, `tmi.twitch.tv 376`) &&
		!strings.Contains(userName, `tmi.twitch.tv 372`) &&
		!strings.Contains(userName, `tmi.twitch.tv 372`) {
		_, err := bt.FileChannelLog[channel].WriteString("[" + timeStamp() + "] [TWITCH] Канал:" + channel +
			" Ник:" + userName + "\tСообщение:" + message + "\n")
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "bots",
				"function": "FileChannelLog[channel].WriteString",
				"file":     "twitchhandlechat.go",
				"body":     "writeLog",
				"error":    err,
			}).Errorln("Ошибка записи лога.")
		}
		if !strings.HasPrefix(message, TwPrefix) {
			bt.MarkovChain += " " + resource.ReadTxt(message)
		}
		fmt.Print("[" + timeStamp() + "] [TWITCH] Канал: " + channel + " " +
			"Ник: " + userName + "\tСообщение: " + message + "\n")
	}
}

func (bt *BotTwitch) checkReact(channel, message string) {
	if _, ok := bt.Settings[channel]; ok {
		if bt.Settings[channel].ReactStatus &&
			(time.Now().Unix()-bt.Settings[channel].LastReactTime >=
				int64(bt.Settings[channel].ReactRate)) {
			for key := range reactTW {
				if strings.Contains(message, key) {
					bt.Settings[channel].LastReactTime = time.Now().Unix()
					go bt.saveChannelSettings(channel)
					bt.say(reactTW[key], channel)
					break
				}
			}
		}
	}
}

func (bt *BotTwitch) handleRewards(message, userName, channel, rewardID string) {
	switch channel {
	case "reflyq":
		switch rewardID {
		case "fa297b45-75cc-4ef2-ba49-841b0fa86ec1":
			msgSlice := strings.Fields(message)
			//switch  {
			//case bt.handleApiRequest(userName,channel, message, "userstate") == "mod"
			//case len(msgSlice) < 1:
			//	bt.say("Ошибка, пустое сообщение.", channel)
			//	return
			//}
			if len(msgSlice) < 1 {
				bt.say("Ошибка, пустое сообщение.", channel)
				return
			}
			bt.say(msgSlice[0]+" заткнули", channel)
			bt.say("/timeout "+msgSlice[0]+" 300 заткнули", channel)
			time.Sleep(300 * time.Second)
			bt.say("/untimeout "+msgSlice[0], channel)
		}
	}
}

func (bt *BotTwitch) handleMasterCmd(message, channel string) {
	switch {
	case strings.HasPrefix(message, "!Ada, switch react"):
		switch bt.Settings[channel].ReactStatus {
		case true:
			bt.Settings[channel].ReactStatus = false
			bt.say("Больше никаких приветов?", channel)
			bt.saveChannelSettings(channel)
			return
		case false:
			bt.Settings[channel].ReactStatus = true
			bt.say("00101!", channel)
			bt.saveChannelSettings(channel)
			return
		}
	case strings.HasPrefix(message, "!Ada, switch"):
		switch bt.Settings[channel].Status {
		case true:
			bt.Settings[channel].Status = false
			bt.say("Засыпаю...", channel)
			bt.saveChannelSettings(channel)
			return
		case false:
			bt.Settings[channel].Status = true
			bt.say("Проснулись, улыбнулись!", channel)
			bt.saveChannelSettings(channel)
			return
		}
	case strings.HasPrefix(message, "!Ada, switch cmd"):
		switch bt.Settings[channel].CMDStatus {
		case true:
			bt.Settings[channel].CMDStatus = false
			bt.say("No !roll's for you", channel)
			bt.saveChannelSettings(channel)
			return
		case false:
			bt.Settings[channel].CMDStatus = true
			bt.say("!roll?", channel)
			bt.saveChannelSettings(channel)
			return
		}
	case strings.HasPrefix(message, "!Ada, show settings"):
		bt.say("Status: "+func() string {
			if bt.Settings[channel].Status {
				return "True"
			} else {
				return "False"
			}
		}()+" ReactStatus: "+func() string {
			if bt.Settings[channel].ReactStatus {
				return "True"
			} else {
				return "False"
			}
		}()+" CMD Status: "+func() string {
			if bt.Settings[channel].CMDStatus {
				return "True"
			} else {
				return "False"
			}
		}()+" Moderator Status: "+func() string {
			if bt.Settings[channel].IsModerator {
				return "True"
			} else {
				return "False"
			}
		}()+" Last react time: "+time.Unix(bt.Settings[channel].LastReactTime, 0).String()+
			" React rate time: "+strconv.Itoa(bt.Settings[channel].ReactRate), channel)
		return
	case strings.HasPrefix(message, "!Ada, set reactrate to"):
		tempstr := strings.Fields(message)
		if len(tempstr) < 5 {
			bt.say("Некорректный ввод", channel)
		} else {
			_, err := strconv.Atoi(tempstr[4])
			if err != nil {
				bt.say("Некорректный ввод", channel)
			} else {
				bt.Settings[channel].ReactRate, _ = strconv.Atoi(tempstr[4])
				go bt.saveChannelSettings(channel)
				bt.say("Частота реакции установлена на раз в "+
					strconv.Itoa(bt.Settings[channel].ReactRate)+" секунд.", channel)
				return
			}
		}
	}
}

func (bt *BotTwitch) handleLine(line string) (user, channel, message, rewardID string) {
	var msgID int
	lineSlice := strings.Fields(strings.Replace(line, ";", " ", -1))
	for id, lin := range lineSlice {
		//fmt.Println("ID:", id, " Field:", lin)
		if lin == "PRIVMSG" && msgID == 0 {
			msgID = id + 2
		}
		if id == msgID && msgID != 0 {
			message = strings.TrimPrefix(lin, `:`)
		}
		if id > msgID && msgID != 0 {
			message += " " + lin
		}
		if strings.Contains(lin, "custom-reward-id=") {
			rewardID = strings.TrimPrefix(lin, "custom-reward-id=")
		}
		if strings.Contains(lin, "display-name=") {
			user = strings.ToLower(strings.TrimPrefix(lin, "display-name="))
		}
		if strings.HasPrefix(lin, `#`) {
			channel = strings.TrimPrefix(lin, `#`)
		}
	}
	//fmt.Println("Ревард:", rewardID)
	return user, channel, message, rewardID
}
