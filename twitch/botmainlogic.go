package twitch

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

func (bt *BotTwitch) initBot() {
	botFile, err := ioutil.ReadFile("BotData.json")
	if err != nil {
		fmt.Println("Ошибка чтения данных бота (BotData.Json),"+
			" должно находиться в корневой папке с исполняемым файлом: ", err)
	}
	err = json.Unmarshal(botFile, bt)
	if err != nil {
		fmt.Println("Ошибка конвертирования структуры из файла в структуру бота: ", err)
	}
}

func (bt *BotTwitch) initSettings() {
	bt.Settings = make(map[string]*botSettings)
	for _, channel := range bt.Channels {
		bt.Settings[channel] = &botSettings{
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
		err = json.Unmarshal(channelSettingsJsonFile, bt.Settings[channel])
		if err != nil {
			fmt.Println("Ошибка конвертирования структуры из файла в структуру настроек: ", err)
		}
		if bt.handleApiRequest("", channel, "", "!evaismod") == "true" {
			bt.Settings[channel].IsModerator = true
		} else {
			bt.Settings[channel].IsModerator = false
		}
		bt.saveSettings(channel)
	}
}

func (bt *BotTwitch) saveSettings(channel string) {
	channelSettingsJson, err := json.MarshalIndent(*bt.Settings[channel], "", " ")
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

func (bt *BotTwitch) initViewersData() {
	bt.Viewers = make(map[string]*viewersData)
	for _, channel := range bt.Channels {
		bt.Viewers[channel] = &viewersData{}
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
		err = json.Unmarshal(channelViwerFileJson, &bt.Viewers[channel].Viewers)
		if err != nil {
			fmt.Println("Ошибка конвертирования структуры из файла в структуру зрителей: ", err)
		}
		bt.saveViewersData(channel)
	}
}

func (bt *BotTwitch) saveViewersData(channel string) {
	bt.handleApiRequest("", channel, "", "requestallviewers")
	channelViewerJson, err := json.MarshalIndent(&bt.Viewers[channel].Viewers, "", " ")
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

func (bt *BotTwitch) openChannelLog() {
	bt.FileChannelLog = make(map[string]*os.File)
	for _, channel := range bt.Channels {
		var err error
		err = os.MkdirAll("logs/"+channel+" Channel", 0777)
		if err != nil && !strings.Contains(err.Error(), "Cannot create a file when that file already exists.") {
			fmt.Println("Не удалось создать директорию для канала:", err)
			err = nil
		}
		bt.FileChannelLog[channel], err = os.OpenFile(
			"logs/"+channel+" Channel/"+channel+" Log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			fmt.Println("Не удалось создать \\ открыть файл:", err)
		}
	}
}

func (bt *BotTwitch) evalute() {
	for {
		var cmd, message, channel string = "", "", ""
		fmt.Scan(&cmd)
		switch cmd {
		case "!ES":
			fmt.Scan(&message)
			message = strings.Replace(message, "!", " ", -1)
			fmt.Scan(&channel)
			bt.say(message, channel)
		}
		time.Sleep(1 * time.Second)
	}
}

func (bt *BotTwitch) Start() {
	var err error
	bt.initBot()
	go bt.initApiConfig()
	go bt.evalute()
	go bt.initViewersData()
	for {
		bt.connect()
		err = bt.joinChannels()
		if err != nil {
			err = nil
			continue
		}
		bt.ReadChannels = textproto.NewReader(bufio.NewReader(bt.Connection))
		err = bt.listenChannels()
		if err != nil {
			err = nil
			time.Sleep(10 * time.Second)
			continue
		} else {
			break
		}
	}
	defer bt.Connection.Close()
}

func (bt *BotTwitch) connect() {
	var err error
	bt.Connection, err = net.Dial("tcp", bt.Server+":"+bt.Port)
	if err != nil {
		fmt.Println("Ошибка попытки соединения: ", err)
		time.Sleep(10 * time.Second)
		bt.connect()
	}
}

func (bt *BotTwitch) joinChannels() error {
	var err error
	_, err = bt.Connection.Write([]byte("PASS " + bt.OAuth + "\r\n"))
	_, err = bt.Connection.Write([]byte("NICK " + bt.BotName + "\r\n"))
	if err != nil {
		fmt.Println("Ошибка во время отправки логина: ", err)
		time.Sleep(10 * time.Second)
		return err
	}
	for _, channel := range bt.Channels {
		_, err := bt.Connection.Write([]byte("JOIN #" + channel + "\r\n"))
		if err != nil {
			fmt.Println("Ошибка во время входа в чат-комнату: ", err)
			return err
		}
	}
	return nil
}

func (bt *BotTwitch) listenChannels() error {
	var err error
	bt.openChannelLog()
	bt.initSettings()
	for _, channelFile := range bt.FileChannelLog {
		defer channelFile.Close()
	}
	go bt.sendBlindRepeatMessage()
	//go bt.startPubSub()
	for {
		if err = bt.handleChat(); err != nil {
			return err
		}
	}
}

func (bt *BotTwitch) say(msg, channel string) {
	if msg == "" {
		fmt.Println("Сообщение пустое")
		return
	}
	_, err := bt.Connection.Write([]byte("PRIVMSG #" + channel + " :" + msg + "\r\n"))
	if err != nil {
		fmt.Println("Ошибка отправки сообщения: ", err)
	}
	bt.FileChannelLog[channel].WriteString("[" + timeStamp() + "] Канал:" + channel +
		" Ник:" + bt.BotName + "\tСообщение:" + msg + "\n")
	fmt.Println("[" + timeStamp() + "] Канал:" + channel + "\tНик:" + bt.BotName + "\tСообщение:" + msg + "\n")
}
