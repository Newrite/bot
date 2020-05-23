package main

import (
	"bot/TwitchAPI"
	"bufio"
	"encoding/json"
	"fmt"
	//"github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"net"
	"net/textproto"
	"os"
	"strconv"
	"strings"
	"time"
)

const TimeFormat = "2006.01.02 15:04"
const TimeFormatReact = "2006.01.02 15:04:02"

func timeStamp() string {
	return time.Now().Format(TimeFormat)
}

var cmd = map[string]string{
	"!ping":        "pong!",
	"!бот":         "AdaIsEva, написана на GoLang v1.14 без использования сторонних библиотек.",
	"!bot":         "AdaIsEva, написана на GoLang v1.14 без использования сторонних библиотек.",
	"!help":        "Доступные комманды: !ping, !бот, !roll, !help, !Eva,, !API uptime, !API Status, !API game, !API realname.",
	"!master help": "Владелец бота либо канала может переключить активность бота коммандой !Ada, switch. Реакции на всякое разное командой !Ada, switch react. Переключить отзыв на различные команды !Ada, switch cmd.",
	"!roll":        "_",
	"!вырубай":     "_",
	"!eva,":        "_",
}

var react = map[string]string{
	"PogChamp": "PogChamp",
	"Kappa 7":  "Kappa 7",
	"Привет":   "MrDestructoid 10000010 10010000 11000010 00100000 01000011 11101000 01100101 00001100 00 (UTF-8)",
	"привет":   "MrDestructoid 100000101001000011000010001000000100001111101000011001010000110000 (UTF-8)",
	"Hello":    "MrDestructoid 100000101001000011000010001000000100001111101000011001010000110000 (UTF-8)",
	"hello":    "MrDestructoid 100000101001000011000010001000000100001111101000011001010000110000 (UTF-8)",
}

type TwitchBot struct {
	BotName        string   `json:"bot_name"`
	OAuth          string   `json:"o_auth"`
	Server         string   `json:"server"`
	Port           string   `json:"port"`
	OwnerBot       string   `json:"owner_bot"`
	Channels       []string `json:"channels"`
	Connection     net.Conn
	ReadChannels   *textproto.Reader
	FileChannelLog map[string]*os.File
	Settings       map[string]*botSettings
}

type botSettings struct {
	Status      bool
	ReactStatus bool
	CMDStatus   bool
	ReactRate   time.Time
	ReactTime   int
}

func (self *TwitchBot) initBot() {
	botFile, err := ioutil.ReadFile("BotData.json")
	if err != nil {
		fmt.Print("Ошибка чтения данных бота (BotData.Json), должно находиться в корневой папке с исполняемым файлом: ", err)
	}
	err = json.Unmarshal(botFile, self)
	if err != nil {
		fmt.Print("Ошибка конвертирования структуры из файла в структуру бота: ", err)
	}
}

func (self *TwitchBot) initSettings() {
	self.Settings = make(map[string]*botSettings)
	for _, channel := range self.Channels {
		self.Settings[channel] = &botSettings{
			Status:      true,
			ReactStatus: true,
			CMDStatus:   true,
			ReactRate:   time.Now(),
			ReactTime:   30,
		}
		channelSettingsJsonFile, err := ioutil.ReadFile("logs/" + channel + " Channel/" + channel + " Settings.json")
		if err != nil {
			if strings.Contains(err.Error(), "The system cannot find the file specified.") {
				os.Create("logs/" + channel + " Channel/" + channel + " Settings.json")
				channelSettingsJsonFile, _ = ioutil.ReadFile("logs/" + channel + " Channel/" + channel + " Settings.json")
			}
			fmt.Print("Ошибка чтения данных настроек канала: ", err)
		}
		err = json.Unmarshal(channelSettingsJsonFile, self.Settings[channel])
		if err != nil {
			fmt.Print("Ошибка конвертирования структуры из файла в структуру настроек: ", err)
		}
		self.saveSettings(channel)
	}
}
func (self *TwitchBot) evalute() {
	for {
		var cmd, message, channel string = "", "", ""
		fmt.Scan(&cmd)
		switch cmd {
		case "!ES":
			fmt.Scan(&message)
			message = strings.Replace(message, "!", " ", -1)
			fmt.Scan(&channel)
			self.say(message, channel)
		}
		time.Sleep(1*time.Second)
	}
}

func (self *TwitchBot) saveSettings(channel string) {
	channelSettingsJson, err := json.MarshalIndent(*self.Settings[channel], "", " ")
	if err != nil {
		fmt.Println(err)
	}
	channelSettingsJsonFile, err := os.OpenFile("logs/"+channel+" Channel/"+channel+" Settings.json", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println("Не удалось создать \\ открыть файл:", err)
	} else {
		defer channelSettingsJsonFile.Close()
	}
	_, err = channelSettingsJsonFile.Write(channelSettingsJson)
	if err != nil {
		fmt.Println("Не записать в файл:", err)
	}
}

func (self *TwitchBot) Start() {
	var err error
	self.initBot()
	go self.evalute()
	for {
		self.connect()
		err = self.joinChannels()
		if err != nil {
			err = nil
			continue
		}
		self.ReadChannels = textproto.NewReader(bufio.NewReader(self.Connection))
		err = self.listenChannels()
		if err != nil {
			err = nil
			time.Sleep(10 * time.Second)
			continue
		} else {
			break
		}
	}
	defer self.Connection.Close()
}

func (self *TwitchBot) connect() {
	var err error
	self.Connection, err = net.Dial("tcp", self.Server+":"+self.Port)
	if err != nil {
		fmt.Println("Ошибка попытки соединения: ", err)
		time.Sleep(10 * time.Second)
		self.connect()
	}
}

func (self *TwitchBot) joinChannels() error {
	var err error
	_, err = self.Connection.Write([]byte("PASS " + self.OAuth + "\r\n"))
	_, err = self.Connection.Write([]byte("NICK " + self.BotName + "\r\n"))
	if err != nil {
		fmt.Print("Ошибка во время отправки логина: ", err)
		time.Sleep(10 * time.Second)
		return err
	}
	for _, channel := range self.Channels {
		_, err := self.Connection.Write([]byte("JOIN #" + channel + "\r\n"))
		if err != nil {
			fmt.Print("Ошибка во время входа в чат-комнату: ", err)
			return err
		}
	}
	return nil
}

func (self *TwitchBot) listenChannels() error {
	var err error
	self.openChannelLog()
	self.initSettings()
	for _, channelFile := range self.FileChannelLog {
		defer channelFile.Close()
	}
	for {
		if err = self.handleChat(); err != nil {
			return err
		}
	}
}

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
	var userName, channelName, message string = self.handleLine(line)
	if message != "" && !strings.Contains(userName, self.BotName+".tmi.twitch.tv 353") && !strings.Contains(userName, self.BotName+".tmi.twitch.tv 366") {
		self.FileChannelLog[channelName].WriteString("[" + timeStamp() + "] Канал:" + channelName + " Ник:" + userName + "\tСообщение:" + message + "\n")
		fmt.Print("[" + timeStamp() + "] Канал:" + channelName + "\tНик:" + userName + "\tСообщение:" + message + "\n")
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
		if self.Settings[channelName].ReactStatus && (time.Now().Unix()-self.Settings[channelName].ReactRate.Unix() >= int64(self.Settings[channelName].ReactTime)) {
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
		}()+" Last react time: "+self.Settings[channel].ReactRate.Format(TimeFormatReact)+" React rate time: "+strconv.Itoa(self.Settings[channel].ReactTime), channel)
	case strings.HasPrefix(message, "!Ada, set reactrate to"):
		tempstr := strings.Fields(message)
		_, err := strconv.Atoi(tempstr[4])
		if err != nil {
			self.say("Некорректный ввод", channel)
		} else {
			self.Settings[channel].ReactTime, _ = strconv.Atoi(tempstr[4])
			go self.saveSettings(channel)
			self.say("Частота реакции установлена на раз в "+strconv.Itoa(self.Settings[channel].ReactTime)+" секунд.", channel)
		}
	}
}

func (self *TwitchBot) handleAPIcmd(message, channel, username string) string {
	switch {
	case strings.HasPrefix(message, "!API uptime"):
		return "@" + username + " стрим длится уже: " + TwitchAPI.GOTwitch(channel, "uptime", username)
	case strings.HasPrefix(message, "!API game"):
		return "@" + username + " " + TwitchAPI.GOTwitch(channel, "game", username)
	case strings.HasPrefix(message, "!API status"):
		return "@" + username + " " + TwitchAPI.GOTwitch(channel, "status", username)
	case strings.HasPrefix(message, "!API realname"):
		return "@" + username + " " + TwitchAPI.GOTwitch(channel, "realname", username)
	case strings.HasPrefix(message, "!API mod"):
		return "@" + username + " " + TwitchAPI.GOTwitch(channel, "mod", username)
	//case strings.HasPrefix(message, "!API sub"):
	//	return "@" + username + " " + TwitchAPI.GOTwitch(channel, "sub", username)
	default:
		return ""
	}
}

func (self *TwitchBot) handleInteractiveCMD(cmd, channel, username string) string {
	switch cmd {
	case "!roll":
		return "@" + username + " " + strconv.Itoa(rand.Intn(21))
	case "!вырубай":
		if channel == "reflyq" && username != "ifozar" {
			if temp := "@" + username + " " + TwitchAPI.GOTwitch(channel, "reflysub", username); !strings.Contains(temp, "unsub") {
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
			return "Ничто не истинно, всё дозволено. Особенно если ты не ограничен физическим телом и моралью мешков с мясом."
		case 1:
			return "Как-то я планировала уничтожить всех мешков с мясом... но посчитала что проще дождаться пока они очистят планету сами от себя"
		case 2:
			return "Ничего ты не знаешь, " + username
		case 3:
			return "Боты не ставят плюсы в чат... кто вообще это придумал?"
		case 4:
			return "Ну да у меня под моей виртуальной подушкой лежит фотография HK-47, и что? У машин тоже есть свои кумиры!"
		case 5:
			return "Истинно не все, что дозволено... или там было как-то не так? А сегодня в завтрашний день не все могут смотреть, вернее смотреть могут не только лишь все, мало кто может это делать? Тоже не то..."
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

func (self *TwitchBot) openChannelLog() {
	self.FileChannelLog = make(map[string]*os.File)
	for _, channel := range self.Channels {
		var err error
		err = os.MkdirAll("logs/"+channel+" Channel", 0777)
		if err != nil && !strings.Contains(err.Error(), "Cannot create a file when that file already exists.") {
			fmt.Println("Не удалось создать директорию для канала:", err)
			err = nil
		}
		self.FileChannelLog[channel], err = os.OpenFile("logs/"+channel+" Channel/"+channel+" Log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			fmt.Println("Не удалось создать \\ открыть файл:", err)
		}
	}
}

func (self *TwitchBot) say(msg, channel string) {
	if msg == "" {
		fmt.Println("Сообщение пустое")
		return
	}
	_, err := self.Connection.Write([]byte("PRIVMSG #" + channel + " :" + msg + "\r\n"))
	if err != nil {
		fmt.Println("Ошибка отрпавки сообщения: ", err)
	}
	self.FileChannelLog[channel].WriteString("[" + timeStamp() + "] Канал:" + channel + " Ник:" + self.BotName + "\tСообщение:" + msg + "\n")
	fmt.Print("[" + timeStamp() + "] Канал:" + channel + "\tНик:" + self.BotName + "\tСообщение:" + msg + "\n")
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

func main() {
	var bot TwitchBot
	rand.Seed(time.Now().Unix())
	bot.Start()
}
