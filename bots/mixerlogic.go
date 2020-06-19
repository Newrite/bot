package bots

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

func (bm *BotMixer) initBot() {
	botFile, err := ioutil.ReadFile("MixerBotData.json")
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "ReadFile",
			"file":     "mixerlogic.go",
			"body":     "initBot",
			"error":    err,
		}).Errorln("Ошибка чтения данных бота (MixerBotData.json)," +
			" должно находиться в корневой папке с исполняемым файлом")
	}
	err = json.Unmarshal(botFile, bm)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "ReadFile",
			"file":     "mixerlogic.go",
			"body":     "initBot",
			"error":    err,
		}).Errorln("Ошибка конвертирования структуры из файла в структуру бота.")
	}
}

func (bm *BotMixer) initMixerConfig() {
	mixConfig, err := ioutil.ReadFile("MixerConfig.json")
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "ioutil.ReadAll",
			"file":     "mixerlogic.go",
			"body":     "initMixerConfig",
			"error":    err,
		}).Errorln("Ошибка чтения данных апи для бота (MixerConfig.json), должно находиться" +
			" в корневой папке с исполняемым файлом.")
	}
	err = json.Unmarshal(mixConfig, &bm.ApiConf)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "json.Unmarshal",
			"file":     "mixerlogic.go",
			"body":     "initMixerConfig",
			"error":    err,
		}).Errorln("Ошибка конвертирования структуры из файла в апи структуру бота.")
	}
}

func (bm *BotMixer) Start() {
	bm.initBot()
	bm.initMixerConfig()
	for {
		authKey := requestChat()
		bm.connect()
		defer mixer.Connection.Close()
		err := bm.joinChannel(authKey)
		if err != nil {
			fmt.Println(err)
			time.Sleep(5 * time.Second)
			continue
		}
		bm.uptime = time.Now().Unix()
		for {
			err = bm.listenChannel()
			if err != nil {
				fmt.Println(err)
				break
			}
		}
	}
}

func (bm *BotMixer) connect() {
	var err error
	head := http.Header{}
	head.Add("Authorization", SingleMixer().ApiConf.Access_token)
	mixer.Connection, _, err = websocket.DefaultDialer.Dial("wss://chat.mixer.com:443", head)
	if err != nil {
		fmt.Println(err)
		time.Sleep(5 * time.Second)
		bm.connect()
	}
}

func (bm *BotMixer) joinChannel(authKey string) error {
	err := mixer.Connection.WriteMessage(websocket.TextMessage, []byte(`{"type":"method","method":"auth","arguments":[`+newriteChannelID+`,`+immersiveEvaUserID+`,"`+authKey+`"],"id":0}`))
	if err != nil {
		return err
	}
	return nil
}

func main() {
	SingleMixer().Start()
}

func templateRequest(method, url, headAuth string) []byte {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if req == nil || err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "NewRequest",
			"file":     "mixerlogic.go",
			"body":     "templateRequest",
			"error":    err,
			"params":   method + " " + url,
		}).Errorln("Request == nil.")
		return nil
	}
	req.Header.Add("Authorization", headAuth)
	resp, err := client.Do(req)
	if err != nil || resp == nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "client.Do",
			"file":     "twitchapi.go",
			"body":     "templateRequest",
			"error":    err,
			"Request":  req,
		}).Errorln("Ошибка обработки реквеста.")
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "ioutil.ReadAll",
			"file":     "mixerlogic.go",
			"body":     "templateRequest",
			"error":    err,
			"Request":  req,
		}).Errorln("Ошибка парсинга тела реквеста в срез байтов.")
		return nil
	}
	return body
}

func requestChat() string {
	var chat = &chatConnect{}
	Url := "https://mixer.com/api/v1/chats/171903618"
	body := templateRequest("GET", Url, SingleMixer().ApiConf.Access_token)
	err := json.Unmarshal(body, chat)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "json.Unmarshal",
			"file":     "mixerlogic.go",
			"body":     "requestChat",
			"error":    err,
		}).Errorln("Ошибка конвертирования структуры.")
	}
	return chat.Authkey
}

func (bm *BotMixer) say(msg string) {
	if msg == "" {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "msg == \"\"",
			"file":     "mixerlogic.go",
			"body":     "say",
			"error":    nil,
		}).Infoln("Пустое сообщение.")
		return
	}
	err := bm.Connection.WriteMessage(websocket.TextMessage, []byte(`{"type":"method","method": "msg","arguments":["`+msg+`"],"id":2}`))
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "Connection.Write",
			"file":     "mixerlogic.go",
			"body":     "say",
			"error":    err,
		}).Errorln("Ошибка во время отправки сообщения на mixer.")
	}
}
