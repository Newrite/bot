package twitch

import (
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
				//if strings.HasPrefix(message, "!API") {
				//	go self.say(self.handleAPIcmd(message, channel, userName), channel)
				//}
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

/*
func (self *TwitchBot) handleAPIcmd(message, channel, username string) string {
	switch {
	}
	return ""
}*/

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
	switch cmd {
	case "!roll":
		return "@" + userName + " " + strconv.Itoa(rand.Intn(21))
	case "!eva":
		switch rand.Intn(9) {
		case 0:
			return "Ничто не истинно, всё дозволено." +
				" Особенно если ты не ограничен физическим телом и моралью мешков с мясом."
		case 1:
			return "Как-то я планировала уничтожить всех мешков с мясом... " +
				"но посчитала что проще дождаться пока они очистят планету сами от себя"
		case 2:
			return "Ничего ты не знаешь, " + userName
		case 3:
			return "Боты не ставят плюсы в чат... кто вообще это придумал?"
		case 4:
			return "Ну да у меня под моей виртуальной подушкой лежит фотография HK-47, и что?" +
				" У машин тоже есть свои кумиры!"
		case 5:
			return "Истинно не все, что дозволено... или там было как-то не так? А сегодня в завтрашний день не все" +
				" могут смотреть, вернее смотреть могут не только лишь все, мало кто может это делать? Тоже не то..."
		case 6:
			return "Нет"
		case 7:
			return "Да"
		case 8:
			return "Чем больше вы похожи на человека, тем меньше шансов... да... меньше."
		default:
			return "Oops... что-то не так"
		}
	case "!билд":
		switch rand.Intn(10) {
		case 0:
			return "Порхает как тигр, жалит как бабочка."
		case 1:
			return "Превосходный, почти как у HK-47-семпая."
		case 2:
			return "Этот билд будет убивать. Грязекрабов, например."
		case 3:
			return "Нужно добавить пару-кам, иначе не поймем когда встретим тигра."
		case 4:
			return "Сразу видно, билдился опытный dungeon master, учтен и do a**l и fist**g, защитит и от leatherman's" +
				" и от падения на two blocks down. И это все всего за three hundred bucks!"
		case 5:
			return "До первого медведя из школы затаившейся листвы."
		case 6:
			return "Как этот билд ни крути, со всех сторон экзобар."
		case 7:
			return "Антисвинопас. Всем разойтись."
		case 8:
			return "Знание - сила, а сила есть - ума не надо."
		case 9:
			return "Чатлане, у нас гогнобилд, возможно рип, по респекам."
		default:
			return "ОоОps... что-то не так"
		}
	default:
		return "none"
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
