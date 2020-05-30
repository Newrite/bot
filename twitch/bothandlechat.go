package twitch

import (
	"bot/reso"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func (self *TwitchBot) handleChat() error {
	line, err := self.ReadChannels.ReadLine()
	if err != nil {
		fmt.Println("Ошибка во время чтения строки: ", err)
		return err
	}
	if line == "PING :tmi.twitch.tv" {
		fmt.Println("PING :tmi.twitch.tv")
		self.Connection.Write([]byte("PONG\r\n"))
		return nil
	}
	var userName, channel, message string = self.handleLine(line)
	self.writeLog(userName, channel, message)
	if strings.Contains(message, "!Ada, ") && (userName == channel || userName == self.OwnerBot) {
		go self.handleMasterCmd(message, channel)
		return nil
	}
	if _, ok := self.Settings[channel]; ok {
		if !self.Settings[channel].Status {
			return nil
		}
	}
	self.checkReact(channel, message)
	message = strings.ToLower(message)
	self.checkCMD(channel, userName, message)
	time.Sleep(10 * time.Millisecond)
	return nil
}
func (self *TwitchBot) writeLog(userName, channel, message string) {
	if message != "" && !strings.Contains(userName, self.BotName+".tmi.twitch.tv 353") &&
		!strings.Contains(userName, self.BotName+".tmi.twitch.tv 366") {
		self.FileChannelLog[channel].WriteString("[" + timeStamp() + "] Канал:" + channel +
			" Ник:" + userName + "\tСообщение:" + message + "\n")
		fmt.Print("[" + timeStamp() + "] Канал:" + channel +
			"\tНик:" + userName + "\tСообщение:" + message + "\n")
	}
}

func (self *TwitchBot) checkReact(channel, message string) {
	if _, ok := self.Settings[channel]; ok {
		if self.Settings[channel].ReactStatus &&
			(time.Now().Unix()-self.Settings[channel].ReactRate.Unix() >=
				int64(self.Settings[channel].ReactTime)) {
			for key := range react {
				if strings.Contains(message, key) {
					self.Settings[channel].ReactRate = time.Now()
					go self.saveSettings(channel)
					self.say(react[key], channel)
					break
				}
			}
		}
	}
}

func (self *TwitchBot) checkCMD(channel, userName, message string) {
	if _, ok := self.Settings[channel]; ok {
		if self.Settings[channel].CMDStatus {
			for key, value := range cmd {
				if strings.HasPrefix(message, key) && value != "_" {
					self.say("@"+userName+" "+cmd[key], channel)
					break
				}
				if strings.HasPrefix(message, key) && value == "_" {
					go self.say(self.handleInteractiveCMD(key, channel, userName, message), channel)
					break
				}
			}
		}
	}
}

func (self *TwitchBot) handleMasterCmd(message, channel string) {
	switch {
	case strings.HasPrefix(message, "!Ada, switch"):
		switch self.Settings[channel].Status {
		case true:
			self.Settings[channel].Status = false
			self.say("Засыпаю...", channel)
			self.saveSettings(channel)
			return
		case false:
			self.Settings[channel].Status = true
			self.say("Проснулись, улыбнулись!", channel)
			self.saveSettings(channel)
			return
		}
	case strings.HasPrefix(message, "!Ada, switch react"):
		switch self.Settings[channel].ReactStatus {
		case true:
			self.Settings[channel].ReactStatus = false
			self.say("Больше никаких приветов?", channel)
			self.saveSettings(channel)
			return
		case false:
			self.Settings[channel].ReactStatus = true
			self.say("00101!", channel)
			self.saveSettings(channel)
			return
		}
	case strings.HasPrefix(message, "!Ada, switch cmd"):
		switch self.Settings[channel].CMDStatus {
		case true:
			self.Settings[channel].CMDStatus = false
			self.say("No !roll's for you", channel)
			self.saveSettings(channel)
			return
		case false:
			self.Settings[channel].CMDStatus = true
			self.say("!roll?", channel)
			self.saveSettings(channel)
			return
		}
	case strings.HasPrefix(message, "!Ada, show settings"):
		self.say("Status: "+func() string {
			if self.Settings[channel].Status {
				return "True"
			} else {
				return "False"
			}
		}()+" ReactStatus: "+func() string {
			if self.Settings[channel].ReactStatus {
				return "True"
			} else {
				return "False"
			}
		}()+" CMD Status: "+func() string {
			if self.Settings[channel].CMDStatus {
				return "True"
			} else {
				return "False"
			}
		}()+" Moderator Status: "+func() string {
			if self.Settings[channel].IsModerator {
				return "True"
			} else {
				return "False"
			}
		}()+" Last react time: "+self.Settings[channel].ReactRate.Format(TimeFormatReact)+
			" React rate time: "+strconv.Itoa(self.Settings[channel].ReactTime), channel)
		return
	case strings.HasPrefix(message, "!Ada, set reactrate to"):
		tempstr := strings.Fields(message)
		if len(tempstr) < 5 {
			self.say("Некорректный ввод", channel)
		} else {
			_, err := strconv.Atoi(tempstr[4])
			if err != nil {
				self.say("Некорректный ввод", channel)
			} else {
				self.Settings[channel].ReactTime, _ = strconv.Atoi(tempstr[4])
				go self.saveSettings(channel)
				self.say("Частота реакции установлена на раз в "+
					strconv.Itoa(self.Settings[channel].ReactTime)+" секунд.", channel)
				return
			}
		}
	case strings.HasPrefix(message, "!Ada, set points"):
		tempstr := strings.Fields(message)
		if len(tempstr) < 6 {
			self.say("Некорректный ввод", channel)
		} else {
			_, err := strconv.Atoi(tempstr[5])
			if err != nil {
				self.say("Некорректный ввод", channel)
			} else {
				tempstr[3] = strings.TrimPrefix(tempstr[3], "@")
				tempstr[3] = strings.ToLower(tempstr[3])
				for _, viewer := range self.Viewers[channel].Viewers {
					if viewer.Name == tempstr[3] {
						viewer.Points, _ = strconv.Atoi(tempstr[5])
						go self.saveViewersData(channel)
						self.say("Поинты "+viewer.Name+" установлены в "+tempstr[5], channel)
						return
					}
				}
				self.say("Не нашла зрителя в базе, попробуйте позже ", channel)
				go self.saveViewersData(channel)
				return
			}
		}
	case strings.HasPrefix(message, "!Ada, show points"):
		tempstr := strings.Fields(message)
		if len(tempstr) < 4 {
			self.say("Некорректный ввод", channel)
		} else {
			tempstr[3] = strings.TrimPrefix(tempstr[3], "@")
			tempstr[3] = strings.ToLower(tempstr[3])
			for _, viewer := range self.Viewers[channel].Viewers {
				if viewer.Name == tempstr[3] {
					go self.saveViewersData(channel)
					self.say("Поинты "+viewer.Name+" "+strconv.Itoa(viewer.Points), channel)
					return
				}
			}
			self.say("Не нашла зрителя в базе, попробуйте позже ", channel)
			go self.saveViewersData(channel)
			return
		}
	}
}

func (self *TwitchBot) handleInteractiveCMD(cmd, channel, userName, message string) string {
	switch channel {
	case "blindwalkerboy":
		if tempAnswer := self.handleBlindCMD(userName, message, cmd); tempAnswer != "none" {
			return tempAnswer
		}
	case "reflyq":
		if tempAnswer := self.handleReflyqCMD(userName, message, cmd); tempAnswer != "none" {
			return tempAnswer
		}
	}
	if cmd == "револьвер выстреливает!" && userName == "moobot" {
		return self.resurrected(message, channel)
	}
	switch cmd {
	case "!roll":
		return "@" + userName + " " + strconv.Itoa(rand.Intn(21))
	case "!uptime":
		return "@" + userName + " " + self.handleApiRequest(userName, channel, message, "uptime")
	case "!eva":
		return reso.EvaAnswers[rand.Intn(16)]
	case "!билд":
		return reso.BuildAnswers[rand.Intn(16)]
	default:
		return ""
	}
}

func (self *TwitchBot) handleLine(line string) (user, channel, message string) {
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
