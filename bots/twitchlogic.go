package bots

import (
	"bufio"
	"encoding/json"
	"fmt"
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
	bt.Settings = make(map[string]*botSettings)
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
	go xandrSendRepeatMessage()
	go bt.sendBlindRepeatMessage()
	go bt.initApiConfig()
	for {
		bt.connect()
		err = bt.joinChannels()
		if err != nil {
			err = nil
			continue
		}
		bt.ReadChannels = textproto.NewReader(bufio.NewReader(bt.Connection))
		bt.uptime = time.Now().Unix()
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
	_, err = bt.Connection.Write([]byte("CAP REQ :twitch.tv/tags twitch.tv/commands twitch.tv/membership\r\n"))
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
	bt.initChannelSettings()
	for _, channelFile := range bt.FileChannelLog {
		defer channelFile.Close()
	}
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
	_, err = bt.FileChannelLog[channel].WriteString("[" + timeStamp() + "] [TWITCH] Канал:" + channel +
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
	fmt.Print("[" + timeStamp() + "] [TWITCH] Канал:" + channel + " Ник:" + bt.BotName + "\tСообщение:" + msg + "\n")
}
