package bot

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/textproto"
	"os"
	"strings"
	"time"
)

func (self *TwitchBot) initBot() {
	botFile, err := ioutil.ReadFile("BotData.json")
	if err != nil {
		fmt.Println("Ошибка чтения данных бота (BotData.Json),"+
			" должно находиться в корневой папке с исполняемым файлом: ", err)
	}
	err = json.Unmarshal(botFile, self)
	if err != nil {
		fmt.Println("Ошибка конвертирования структуры из файла в структуру бота: ", err)
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
			IsModerator: false,
		}
		channelSettingsJsonFile, err := ioutil.ReadFile(
			"logs/" + channel + " Channel/" + channel + " Settings.json")
		if err != nil {
			if strings.Contains(err.Error(), "The system cannot find the file specified.") {
				os.Create("logs/" + channel + " Channel/" + channel + " Settings.json")
				channelSettingsJsonFile, _ = ioutil.ReadFile(
					"logs/" + channel + " Channel/" + channel + " Settings.json")
			}
			fmt.Println("Ошибка чтения данных настроек канала: ", err)
		}
		err = json.Unmarshal(channelSettingsJsonFile, self.Settings[channel])
		if err != nil {
			fmt.Println("Ошибка конвертирования структуры из файла в структуру настроек: ", err)
		}
		if self.handleApiRequest("", channel, "", "!evaismod") == "true" {
			self.Settings[channel].IsModerator = true
		} else {
			self.Settings[channel].IsModerator = false
		}
		self.saveSettings(channel)
	}
}

func (self *TwitchBot) saveSettings(channel string) {
	channelSettingsJson, err := json.MarshalIndent(*self.Settings[channel], "", " ")
	if err != nil {
		fmt.Println(err)
	}
	channelSettingsJsonFile, err := os.OpenFile(
		"logs/"+channel+" Channel/"+channel+" Settings.json", os.O_WRONLY|os.O_CREATE, 0600)
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

func (self *TwitchBot) initViewersData() {
	self.Viewers = make(map[string]*viewersData)
	for _, channel := range self.Channels {
		self.Viewers[channel] = &viewersData{}
		channelViwerFileJson, err := ioutil.ReadFile(
			"logs/" + channel + " Channel/" + channel + " ViewersData.json")
		if err != nil {
			if strings.Contains(err.Error(), "The system cannot find the file specified.") {
				os.Create("logs/" + channel + " Channel/" + channel + " ViewersData.json")
				channelViwerFileJson, _ = ioutil.ReadFile(
					"logs/" + channel + " Channel/" + channel + " ViewersData.json")
			}
			fmt.Println("Ошибка чтения данных зрителей канала: ", err)
		}
		err = json.Unmarshal(channelViwerFileJson, &self.Viewers[channel].Viewers)
		if err != nil {
			fmt.Println("Ошибка конвертирования структуры из файла в структуру зрителей: ", err)
		}
		self.saveViewersData(channel)
	}
}

func (self *TwitchBot) saveViewersData(channel string) {
	self.handleApiRequest("", channel, "", "requestallviewers")
	channelViewerJson, err := json.MarshalIndent(&self.Viewers[channel].Viewers, "", " ")
	if err != nil {
		fmt.Println(err)
	}
	channelViewersFileJson, err := os.OpenFile(
		"logs/"+channel+" Channel/"+channel+" ViewersData.json", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println("Не удалось создать \\ открыть файл:", err)
	} else {
		defer channelViewersFileJson.Close()
	}
	_, err = channelViewersFileJson.Write(channelViewerJson)
	if err != nil {
		fmt.Println("Не записать в файл:", err)
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
		self.FileChannelLog[channel], err = os.OpenFile(
			"logs/"+channel+" Channel/"+channel+" Log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			fmt.Println("Не удалось создать \\ открыть файл:", err)
		}
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
		time.Sleep(1 * time.Second)
	}
}

func (self *TwitchBot) Start() {
	var err error
	self.initBot()
	self.initApiConfig()
	go self.evalute()
	self.initViewersData()
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
		fmt.Println("Ошибка во время отправки логина: ", err)
		time.Sleep(10 * time.Second)
		return err
	}
	for _, channel := range self.Channels {
		_, err := self.Connection.Write([]byte("JOIN #" + channel + "\r\n"))
		if err != nil {
			fmt.Println("Ошибка во время входа в чат-комнату: ", err)
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

func (self *TwitchBot) say(msg, channel string) {
	if msg == "" {
		fmt.Println("Сообщение пустое")
		return
	}
	_, err := self.Connection.Write([]byte("PRIVMSG #" + channel + " :" + msg + "\r\n"))
	if err != nil {
		fmt.Println("Ошибка отрпавки сообщения: ", err)
	}
	self.FileChannelLog[channel].WriteString("[" + timeStamp() + "] Канал:" + channel +
		" Ник:" + self.BotName + "\tСообщение:" + msg + "\n")
	fmt.Println("[" + timeStamp() + "] Канал:" + channel + "\tНик:" + self.BotName + "\tСообщение:" + msg + "\n")
}
