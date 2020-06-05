package bots

import (
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
	if line == "PING :tmi.Twitch.tv" {
		fmt.Println("PING :tmi.Twitch.tv")
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
	var userName, channel, message string = bt.handleLine(line)
	bt.writeLog(userName, channel, message)
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
		bt.say(bt.resurrected(message, channel), channel)
	}
	message = strings.ToLower(message)
	if strings.HasPrefix(message, TwPrefix) {
		if _, ok := bt.Settings[channel]; ok {
			if bt.Settings[channel].CMDStatus {
				msgSl := strings.Fields(message)
				go bt.say(checkCMD(userName, channel, msgSl[0], "TW", message), channel)
			}
		}
	}
	time.Sleep(10 * time.Millisecond)
	return nil
}
func (bt *BotTwitch) writeLog(userName, channel, message string) {
	if message != "" && !strings.Contains(userName, bt.BotName+".tmi.Twitch.tv 353") &&
		!strings.Contains(userName, bt.BotName+".tmi.Twitch.tv 366") {
		_, err := bt.FileChannelLog[channel].WriteString("[" + timeStamp() + "] Канал:" + channel +
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
		log.Infof("Ник: %s Канал: %s Сообщение: %s\n", userName, channel, message)
		//fmt.Print("[" + timeStamp() + "] Канал:" + channel +
		//	"\tНик:" + userName + "\tСообщение:" + message + "\n")
	}
}

func (bt *BotTwitch) checkReact(channel, message string) {
	if _, ok := bt.Settings[channel]; ok {
		if bt.Settings[channel].ReactStatus &&
			(time.Now().Unix()-bt.Settings[channel].ReactRate.Unix() >=
				int64(bt.Settings[channel].ReactTime)) {
			for key := range reactTW {
				if strings.Contains(message, key) {
					bt.Settings[channel].ReactRate = time.Now()
					go bt.saveSettings(channel)
					bt.say(reactTW[key], channel)
					break
				}
			}
		}
	}
}

func (bt *BotTwitch) handleMasterCmd(message, channel string) {
	switch {
	case strings.HasPrefix(message, "!Ada, switch"):
		switch bt.Settings[channel].Status {
		case true:
			bt.Settings[channel].Status = false
			bt.say("Засыпаю...", channel)
			bt.saveSettings(channel)
			return
		case false:
			bt.Settings[channel].Status = true
			bt.say("Проснулись, улыбнулись!", channel)
			bt.saveSettings(channel)
			return
		}
	case strings.HasPrefix(message, "!Ada, switch reactTW"):
		switch bt.Settings[channel].ReactStatus {
		case true:
			bt.Settings[channel].ReactStatus = false
			bt.say("Больше никаких приветов?", channel)
			bt.saveSettings(channel)
			return
		case false:
			bt.Settings[channel].ReactStatus = true
			bt.say("00101!", channel)
			bt.saveSettings(channel)
			return
		}
	case strings.HasPrefix(message, "!Ada, switch cmd"):
		switch bt.Settings[channel].CMDStatus {
		case true:
			bt.Settings[channel].CMDStatus = false
			bt.say("No !roll's for you", channel)
			bt.saveSettings(channel)
			return
		case false:
			bt.Settings[channel].CMDStatus = true
			bt.say("!roll?", channel)
			bt.saveSettings(channel)
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
		}()+" Last reactTW time: "+bt.Settings[channel].ReactRate.Format(TimeFormatReact)+
			" React rate time: "+strconv.Itoa(bt.Settings[channel].ReactTime), channel)
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
				bt.Settings[channel].ReactTime, _ = strconv.Atoi(tempstr[4])
				go bt.saveSettings(channel)
				bt.say("Частота реакции установлена на раз в "+
					strconv.Itoa(bt.Settings[channel].ReactTime)+" секунд.", channel)
				return
			}
		}
	case strings.HasPrefix(message, "!Ada, set points"):
		tempstr := strings.Fields(message)
		if len(tempstr) < 6 {
			bt.say("Некорректный ввод", channel)
		} else {
			_, err := strconv.Atoi(tempstr[5])
			if err != nil {
				bt.say("Некорректный ввод", channel)
			} else {
				tempstr[3] = strings.TrimPrefix(tempstr[3], "@")
				tempstr[3] = strings.ToLower(tempstr[3])
				for _, viewer := range bt.Viewers[channel].Viewers {
					if viewer.Name == tempstr[3] {
						viewer.Points, _ = strconv.Atoi(tempstr[5])
						go bt.saveViewersData(channel)
						bt.say("Поинты "+viewer.Name+" установлены в "+tempstr[5], channel)
						return
					}
				}
				bt.say("Не нашла зрителя в базе, попробуйте позже ", channel)
				go bt.saveViewersData(channel)
				return
			}
		}
	case strings.HasPrefix(message, "!Ada, show points"):
		tempstr := strings.Fields(message)
		if len(tempstr) < 4 {
			bt.say("Некорректный ввод", channel)
		} else {
			tempstr[3] = strings.TrimPrefix(tempstr[3], "@")
			tempstr[3] = strings.ToLower(tempstr[3])
			for _, viewer := range bt.Viewers[channel].Viewers {
				if viewer.Name == tempstr[3] {
					go bt.saveViewersData(channel)
					bt.say("Поинты "+viewer.Name+" "+strconv.Itoa(viewer.Points), channel)
					return
				}
			}
			bt.say("Не нашла зрителя в базе, попробуйте позже ", channel)
			go bt.saveViewersData(channel)
			return
		}
	}
}

func (bt *BotTwitch) handleLine(line string) (user, channel, message string) {
	var temp int
	for _, sym := range line {
		if sym == '!' {
			break
		}
		if sym != ':' {
			user += string(sym)
		}
	}
	for _, sym := range line {
		if sym == '#' {
			temp = 1
			continue
		}
		if temp == 1 && sym == ' ' {
			break
		}
		if temp == 1 {
			channel += string(sym)
		}
	}
	temp = 0
	for _, sym := range line {
		if sym == ':' {
			temp += 1
			continue
		}
		if temp == 2 && sym == '\n' {
			break
		}
		if temp == 2 {
			message += string(sym)
		}
	}
	return user, channel, message
}
