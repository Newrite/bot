package bots

import (
	"bufio"
	"encoding/json"
	log "github.com/sirupsen/logrus"
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
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "ioutil.ReadFile",
			"file":     "twitchlogic.go",
			"body":     "initBot",
			"error":    err,
		}).Fatalln("Ошибка чтения данных бота (BotData.Json), " +
			"должно находиться в корневой папке с исполняемым файлом.")
	}
	err = json.Unmarshal(botFile, bt)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "json.Unmarshal",
			"file":     "twitchlogic.go",
			"body":     "initBot",
			"error":    err,
		}).Fatalln("Ошибка конвертирования структуры из файла в структуру бота.")
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
				log.WithFields(log.Fields{
					"package":  "bots",
					"function": "ioutil.ReadFile",
					"file":     "twitchlogic.go",
					"body":     "initSettings",
					"error":    err,
				}).Info("Ошибка открытия настроек.")
				err = nil
				_, err = os.Create("logs/" + channel + " Channel/" + channel + " Settings.json")
				if err != nil {
					log.WithFields(log.Fields{
						"package":  "bots",
						"function": "os.Create",
						"file":     "twitchlogic.go",
						"body":     "initSettings",
						"error":    err,
					}).Errorln("Ошибка создания файла для лога чата.")
				}
				channelSettingsJsonFile, err = ioutil.ReadFile(
					"logs/" + channel + " Channel/" + channel + " Settings.json")
				if err != nil {
					log.WithFields(log.Fields{
						"package":  "bots",
						"function": "ioutil.ReadFile",
						"file":     "twitchlogic.go",
						"body":     "initSettings",
						"error":    err,
					}).Errorln("Ошибка открытия настроек.")
				}
				err = nil
			}
			log.WithFields(log.Fields{
				"package":  "bots",
				"function": "ioutil.ReadFile",
				"file":     "twitchlogic.go",
				"body":     "initSettings",
				"error":    err,
			}).Errorln("Ошибка чтения данных настроек канала.")
		}
		err = json.Unmarshal(channelSettingsJsonFile, bt.Settings[channel])
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "bots",
				"function": "json.Unmarshal",
				"file":     "twitchlogic.go",
				"body":     "initSettings",
				"error":    err,
			}).Errorln("Ошибка конвертирования структуры из файла в структуру настроек.")
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
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "json.MarshalIndent",
			"file":     "twitchlogic.go",
			"body":     "saveSettings",
			"error":    err,
		}).Errorln("Ошибка маршал настроек в жисон.")
	}
	channelSettingsJsonFile, err := os.OpenFile(
		"logs/"+channel+" Channel/"+channel+" Settings.json", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "os.OpenFile",
			"file":     "twitchlogic.go",
			"body":     "saveSettings",
			"error":    err,
		}).Errorln("Не удалось создать \\ открыть файл.")
	} else {
		defer channelSettingsJsonFile.Close()
	}
	if channelSettingsJsonFile != nil {
		_, err = channelSettingsJsonFile.Write(channelSettingsJson)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "bots",
				"function": "channelSettingsJsonFile.Write",
				"file":     "twitchlogic.go",
				"body":     "saveSettings",
				"error":    err,
			}).Errorln("Не записать в файл.")
		}
	} else {
		log.WithFields(log.Fields{
			"package":                 "bots",
			"function":                "channelSettingsJsonFile == nil",
			"file":                    "twitchlogic.go",
			"body":                    "saveSettings",
			"error":                   err,
			"channelSettingsJsonFile": channelSettingsJsonFile,
		}).Errorln("Файл пуст.")
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
				log.WithFields(log.Fields{
					"package":  "bots",
					"function": "ioutil.ReadFile",
					"file":     "twitchlogic.go",
					"body":     "initViewersData",
					"error":    err,
				}).Errorln("Ошибка открытия json зрителей.")
				err = nil
				_, err = os.Create("logs/" + channel + " Channel/" + channel + " ViewersData.json")
				if err != nil {
					log.WithFields(log.Fields{
						"package":  "bots",
						"function": "os.Create",
						"file":     "twitchlogic.go",
						"body":     "initSettings",
						"error":    err,
					}).Errorln("Ошибка создания файла для зрителей.")
					err = nil
				}
				channelViwerFileJson, err = ioutil.ReadFile(
					"logs/" + channel + " Channel/" + channel + " ViewersData.json")
				if err != nil {
					log.WithFields(log.Fields{
						"package":  "bots",
						"function": "ioutil.ReadFile",
						"file":     "twitchlogic.go",
						"body":     "initSettings",
						"error":    err,
					}).Errorln("Ошибка открытия json зрителей после создания.")
					err = nil
				}
			}
			log.WithFields(log.Fields{
				"package":  "bots",
				"function": "ioutil.ReadFile",
				"file":     "twitchlogic.go",
				"body":     "initSettings",
				"error":    err,
			}).Errorln("Ошибка чтения данных зрителей канала в конце.")
		}
		err = json.Unmarshal(channelViwerFileJson, &bt.Viewers[channel].Viewers)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "bots",
				"function": "json.Unmarshal",
				"file":     "twitchlogic.go",
				"body":     "initSettings",
				"error":    err,
			}).Errorln("Ошибка конвертирования структуры из файла в структуру зрителей.")
		}
		bt.saveViewersData(channel)
	}
}

func (bt *BotTwitch) saveViewersData(channel string) {
	bt.handleApiRequest("", channel, "", "requestallviewers")
	channelViewerJson, err := json.MarshalIndent(&bt.Viewers[channel].Viewers, "", " ")
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "json.MarshalIndent",
			"file":     "twitchlogic.go",
			"body":     "saveViewersData",
			"error":    err,
		}).Errorln("Ошибка маршал настроек в жисон зрителей.")
	}
	channelViewersFileJson, err := os.OpenFile(
		"logs/"+channel+" Channel/"+channel+" ViewersData.json", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "os.OpenFile",
			"file":     "twitchlogic.go",
			"body":     "saveViewersData",
			"error":    err,
		}).Errorln("Не удалось создать \\ открыть файл.")
	} else {
		defer channelViewersFileJson.Close()
	}
	if channelViewersFileJson != nil {
		_, err = channelViewersFileJson.Write(channelViewerJson)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "bots",
				"function": "channelSettingsJsonFile.Write",
				"file":     "twitchlogic.go",
				"body":     "saveSettings",
				"error":    err,
			}).Errorln("Не записать в файл.")
		}
	} else {
		log.WithFields(log.Fields{
			"package":                "bots",
			"function":               "channelViewersFileJson == nil",
			"file":                   "twitchlogic.go",
			"body":                   "saveSettings",
			"error":                  err,
			"channelViewersFileJson": channelViewersFileJson,
		}).Errorln("Файл пуст.")
	}
}

func (bt *BotTwitch) openChannelLog() {
	bt.FileChannelLog = make(map[string]*os.File)
	for _, channel := range bt.Channels {
		var err error
		err = os.MkdirAll("logs/"+channel+" Channel", 0777)
		if err != nil && !strings.Contains(err.Error(), "Cannot create a file when that file already exists.") {
			log.WithFields(log.Fields{
				"package":  "bots",
				"function": "os.MkdirAll",
				"file":     "twitchlogic.go",
				"body":     "openChannelLog",
				"error":    err,
			}).Errorln("Не удалось создать директорию для канала.")
			err = nil
		} else if err != nil {
			log.WithFields(log.Fields{
				"package":  "bots",
				"function": "os.MkdirAll",
				"file":     "twitchlogic.go",
				"body":     "openChannelLog",
				"error":    err,
			}).Infoln("Не удалось создать директорию для канала.")
			err = nil
		}
		bt.FileChannelLog[channel], err = os.OpenFile(
			"logs/"+channel+" Channel/"+channel+" Log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "bots",
				"function": "os.OpenFile",
				"file":     "twitchlogic.go",
				"body":     "openChannelLog",
				"error":    err,
			}).Errorln("Не удалось создать \\ открыть файл.")
		}
	}
}

func (bt *BotTwitch) Start() {
	var err error
	bt.initBot()
	go bt.initApiConfig()
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
			log.WithFields(log.Fields{
				"package":  "bots",
				"function": "listenChannels",
				"file":     "twitchlogic.go",
				"body":     "Start",
				"error":    err,
			}).Errorln("Ошибка прослушки чата.")
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
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "net.Dial",
			"file":     "twitchlogic.go",
			"body":     "connect",
			"error":    err,
		}).Errorln("Ошибка попытки соединения.")
		time.Sleep(10 * time.Second)
		bt.connect()
	}
}

func (bt *BotTwitch) joinChannels() error {
	var err error
	_, err = bt.Connection.Write([]byte("PASS " + bt.OAuth + "\r\n"))
	_, err = bt.Connection.Write([]byte("NICK " + bt.BotName + "\r\n"))
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "Connection.Write",
			"file":     "twitchlogic.go",
			"body":     "joinChannels",
			"error":    err,
		}).Errorln("Ошибка во время отправки логина.")
		time.Sleep(10 * time.Second)
		return err
	}
	for _, channel := range bt.Channels {
		_, err := bt.Connection.Write([]byte("JOIN #" + channel + "\r\n"))
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "bots",
				"function": "Connection.Write",
				"file":     "twitchlogic.go",
				"body":     "joinChannels",
				"error":    err,
				"channel":  channel,
			}).Errorln("Ошибка во время входа в чат-комнату.")
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
			log.WithFields(log.Fields{
				"package":  "bots",
				"function": "handleChat",
				"file":     "twitchlogic.go",
				"body":     "listenChannels",
				"error":    err,
			}).Errorln("Ошибка обработки чата.")
			return err
		}
	}
}

func (bt *BotTwitch) say(msg, channel string) {
	if msg == "" {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "msg == \"\"",
			"file":     "twitchlogic.go",
			"body":     "say",
			"error":    nil,
		}).Infoln("Пустое сообщение.")
		return
	}
	_, err := bt.Connection.Write([]byte("PRIVMSG #" + channel + " :" + msg + "\r\n"))
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "Connection.Write",
			"file":     "twitchlogic.go",
			"body":     "say",
			"error":    err,
		}).Errorln("Ошибка отправки сообщения.")
	}
	_, err = bt.FileChannelLog[channel].WriteString("[" + timeStamp() + "] Канал:" + channel +
		" Ник:" + bt.BotName + "\tСообщение:" + msg + "\n")
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "FileChannelLog[channel].WriteString",
			"file":     "twitchlogic.go",
			"body":     "say",
			"error":    err,
		}).Errorln("Ошибка записи лога.")
	}
	log.Infoln("Канал:" + channel + "\tНик:" + bt.BotName + "\tСообщение:" + msg)
	//fmt.Println("[" + timeStamp() + "] Канал:" + channel + "\tНик:" + bt.BotName + "\tСообщение:" + msg + "\n")
}
