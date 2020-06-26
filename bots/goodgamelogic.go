package bots

import (
	"bot/controllers"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"time"
)

func (bgg *BotGoodGame) readServer() string {
	var err error
	if bgg.n, err = bgg.Connection.Read(bgg.serverResponse); err != nil {
		controllers.SingleLog().WithFields(log.Fields{
			"package":  "bots",
			"function": "Connection.Read",
			"file":     "goodgamelogic.go",
			"body":     "readServer",
			"error":    err,
		}).Errorln("Ошибка чтения ответа от сервера.")
	}
	return fmt.Sprintf(string(bgg.serverResponse[:bgg.n]))
}

func (bgg *BotGoodGame) initBot() {
	botFile, err := ioutil.ReadFile("GGBotData.json")
	if err != nil {
		controllers.SingleLog().WithFields(log.Fields{
			"package":  "bots",
			"function": "ReadFile",
			"file":     "goodgamelogic.go",
			"body":     "initBots",
			"error":    err,
		}).Errorln("Ошибка чтения данных бота (GGBotData.Json)," +
			" должно находиться в корневой папке с исполняемым файлом")
	}
	err = json.Unmarshal(botFile, bgg)
	if err != nil {
		controllers.SingleLog().WithFields(log.Fields{
			"package":  "bots",
			"function": "ReadFile",
			"file":     "goodgamelogic.go",
			"body":     "initBots",
			"error":    err,
		}).Errorln("Ошибка конвертирования структуры из файла в структуру бота.")
	}
}

func (bgg *BotGoodGame) connect() {
	var err error
	bgg.Connection, err = websocket.Dial(bgg.Server, "", bgg.Origin)
	if err != nil {
		controllers.SingleLog().WithFields(log.Fields{
			"package":  "bots",
			"function": "Dial",
			"file":     "goodgamelogic.go",
			"body":     "connect",
			"error":    err,
		}).Errorln("Ошибка установки соединения.")
		time.Sleep(10 * time.Second)
		bgg.connect()
	}
	fmt.Println(bgg.readServer())
	_, err = bgg.Connection.Write([]byte(`{"type":"auth","data":{"user_id":"` + bgg.BotId + `","token":"` + bgg.Token + `"}}`))
	if err != nil {
		controllers.SingleLog().WithFields(log.Fields{
			"package":  "bots",
			"function": "Connection.Write",
			"file":     "goodgamelogic.go",
			"body":     "connect",
			"error":    err,
		}).Errorln("Ошибка во время отправки логина.")
		time.Sleep(10 * time.Second)
	}
	fmt.Println(bgg.readServer())
}

func (bgg *BotGoodGame) joinChannels() error {
	var err error
	for _, channel := range bgg.Channels {
		_, err = bgg.Connection.Write([]byte(`{"type":"join","data":{"channel_id":"` + channel + `","hidden":false}}`))
		fmt.Println(bgg.readServer())
		if err != nil {
			controllers.SingleLog().WithFields(log.Fields{
				"package":  "bots",
				"function": "Connection.Write",
				"file":     "goodgamelogic.go",
				"body":     "joinChannels",
				"error":    err,
			}).Errorln("Ошибка во время входа в чат-комнату.")
			return err
		}
	}
	return nil
}

func (bgg *BotGoodGame) Start() {
	var err error
	bgg.initBot()
	bgg.serverResponse = make([]byte, 1024)
	for {
		bgg.connect()
		err = bgg.joinChannels()
		if err != nil {
			err = nil
			continue
		}
		bgg.uptime = time.Now().Unix()
		err = bgg.listenChannels()
		if err != nil {
			err = nil
			time.Sleep(10 * time.Second)
			continue
		} else {
			break
		}
	}
	defer bgg.Connection.Close()
}

func (bgg *BotGoodGame) listenChannels() error {
	var err error
	for {
		if err = bgg.handleChat(); err != nil {
			controllers.SingleLog().WithFields(log.Fields{
				"package":  "bots",
				"function": "handleChat",
				"file":     "goodgamelogic.go",
				"body":     "listenChannels",
				"error":    err,
			}).Errorln("Получена ошибка во время прослушки чата.")
			return err
		}
	}
}

func (bgg *BotGoodGame) say(msg, channel string) {
	if msg == "" {
		controllers.SingleLog().WithFields(log.Fields{
			"package":  "bots",
			"function": "msg == \"\"",
			"file":     "goodgamelogic.go",
			"body":     "say",
			"error":    nil,
		}).Infoln("Пустое сообщение.")
		return
	}
	fmt.Println(msg)
	fmt.Println(channel)
	_, err := bgg.Connection.Write([]byte(`{"type":"send_message","data":{"channel_id":"` + channel + `","text":"` + msg + `","hideIcon":false,"mobile":false}}`))
	if err != nil {
		controllers.SingleLog().WithFields(log.Fields{
			"package":  "bots",
			"function": "Connection.Write",
			"file":     "goodgamelogic.go",
			"body":     "say",
			"error":    err,
		}).Errorln("Ошибка во время отправки сообщения на гг.")
	}
	fmt.Println(bgg.readServer())
}
