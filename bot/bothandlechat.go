package bot

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
	fmt.Println(line)
	var userName, channelName, message string = self.handleLine(line)
	if message != "" && !strings.Contains(userName, self.BotName+".tmi.twitch.tv 353") &&
		!strings.Contains(userName, self.BotName+".tmi.twitch.tv 366") {
		self.FileChannelLog[channelName].WriteString("[" + timeStamp() + "] Канал:" + channelName +
			" Ник:" + userName + "\tСообщение:" + message + "\n")
		fmt.Print("[" + timeStamp() + "] Канал:" + channelName +
			"\tНик:" + userName + "\tСообщение:" + message + "\n")
	}
	if strings.Contains(message, "!Ada, ") && (userName == channelName || userName == self.OwnerBot) {
		go self.handleMasterCmd(message, channelName)
		return nil
	}
	if _, ok := self.Settings[channelName]; ok {
		if !self.Settings[channelName].Status {
			return nil
		}
	}
	if _, ok := self.Settings[channelName]; ok {
		if self.Settings[channelName].ReactStatus &&
			(time.Now().Unix()-self.Settings[channelName].ReactRate.Unix() >=
				int64(self.Settings[channelName].ReactTime)) {
			for key := range react {
				if strings.Contains(message, key) {
					self.Settings[channelName].ReactRate = time.Now()
					go self.saveSettings(channelName)
					self.say(react[key], channelName)
					break
				}
			}
		}
	}
	message = strings.ToLower(message)
	if _, ok := self.Settings[channelName]; ok {
		if self.Settings[channelName].CMDStatus {
			for key, value := range cmd {
				if strings.HasPrefix(message, key) && value != "_" {
					self.say("@"+userName+" "+cmd[key], channelName)
					break
				}
				if strings.HasPrefix(message, key) && value == "_" {
					go self.say(self.handleInteractiveCMD(key, channelName, userName), channelName)
					break
				}
				if strings.HasPrefix(message, "!API") {
					go self.say(self.handleAPIcmd(message, channelName, userName), channelName)
				}
			}
		}
	}
	time.Sleep(10 * time.Millisecond)
	return nil
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
		}()+" Last react time: "+self.Settings[channel].ReactRate.Format(TimeFormatReact)+
			" React rate time: "+strconv.Itoa(self.Settings[channel].ReactTime), channel)
	case strings.HasPrefix(message, "!Ada, set reactrate to"):
		tempstr := strings.Fields(message)
		_, err := strconv.Atoi(tempstr[4])
		if err != nil {
			self.say("Некорректный ввод", channel)
		} else {
			self.Settings[channel].ReactTime, _ = strconv.Atoi(tempstr[4])
			go self.saveSettings(channel)
			self.say("Частота реакции установлена на раз в "+
				strconv.Itoa(self.Settings[channel].ReactTime)+" секунд.", channel)
		}
	}
}

func (self *TwitchBot) handleAPIcmd(message, channel, username string) string {
	switch {
	}
	return ""
}

func (self *TwitchBot) handleInteractiveCMD(cmd, channel, username string) string {
	switch cmd {
	case "!roll":
		return "@" + username + " " + strconv.Itoa(rand.Intn(21))
	case "!вырубай":
		if channel == "reflyq" && username != "ifozar" {
			if temp := "@" + username + " " + self.handleApiRequest(username, channel, "none", "!вырубай"); !strings.Contains(
				temp, "unsub") {
				return temp
			} else {
				self.say("@"+username+" Я тебя щас нахуй вырублю, ансаб блять НЫА roflanEbalo", channel)
				return "/timeout " + username + " 120"
			}
		} else if channel == "reflyq" && username == "ifozar" {
			self.say("iFozar заебал уже эту хуйню писать", channel)
			return "/timeout ifozar 300"
		}
		return ""
	case "!eva,":
		switch rand.Intn(9) {
		case 0:
			return "Ничто не истинно, всё дозволено." +
				" Особенно если ты не ограничен физическим телом и моралью мешков с мясом."
		case 1:
			return "Как-то я планировала уничтожить всех мешков с мясом... " +
				"но посчитала что проще дождаться пока они очистят планету сами от себя"
		case 2:
			return "Ничего ты не знаешь, " + username
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
