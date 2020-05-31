package goodgame

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"time"
)

func (bgg *BotGoodGame) readServer() string {
	var err error
	if bgg.n, err = bgg.Connection.Read(bgg.serverResponse); err != nil {
		fmt.Println(err)
	}
	return fmt.Sprintf(string(bgg.serverResponse[:bgg.n]))
}

func (bgg *BotGoodGame) initBot() {
	botFile, err := ioutil.ReadFile("GGBotData.json")
	if err != nil {
		fmt.Println("Ошибка чтения данных бота (GGBotData.Json),"+
			" должно находиться в корневой папке с исполняемым файлом: ", err)
	}
	err = json.Unmarshal(botFile, bgg)
	if err != nil {
		fmt.Println("Ошибка конвертирования структуры из файла в структуру бота: ", err)
	}
}

func (bgg *BotGoodGame) connect() {
	var err error
	bgg.Connection, err = websocket.Dial(bgg.Server, "", bgg.Origin)
	if err != nil {
		fmt.Println("Ошибка попытки соединения: ", err)
		time.Sleep(10 * time.Second)
		bgg.connect()
	}
	fmt.Println(bgg.readServer())
	_, err = bgg.Connection.Write([]byte(`{"type":"auth","data":{"user_id":"` + bgg.BotId + `","token":"` + bgg.Token + `"}}`))
	if err != nil {
		fmt.Println("Ошибка во время отправки логина: ", err)
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
			fmt.Println("Ошибка во время входа в чат-комнату: ", err)
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
			return err
		}
	}
}

func (bgg *BotGoodGame) say(msg, channel string) {
	if msg == "" {
		fmt.Println("Сообщение пустое")
	}
	bgg.Connection.Write([]byte(`{"type":"send_message","data":{"channel_id":"` + channel + `","text":"` + msg + `","hideIcon":false,"mobile":false}}`))
	fmt.Println(bgg.readServer())
}
