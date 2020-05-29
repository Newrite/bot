package goodgame

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"time"
)

func (self *GoodGameBot) readServer() string {
	var err error
	if self.n, err = self.Connection.Read(self.serverResponse); err != nil {
		fmt.Println(err)
	}
	return fmt.Sprintf(string(self.serverResponse[:self.n]))
}

func (self *GoodGameBot) initBot() {
	botFile, err := ioutil.ReadFile("GGBotData.json")
	if err != nil {
		fmt.Println("Ошибка чтения данных бота (GGBotData.Json),"+
			" должно находиться в корневой папке с исполняемым файлом: ", err)
	}
	err = json.Unmarshal(botFile, self)
	if err != nil {
		fmt.Println("Ошибка конвертирования структуры из файла в структуру бота: ", err)
	}
}

func (self *GoodGameBot) connect() {
	var err error
	self.Connection, err = websocket.Dial(self.Server, "", self.Origin)
	if err != nil {
		fmt.Println("Ошибка попытки соединения: ", err)
		time.Sleep(10 * time.Second)
		self.connect()
	}
	fmt.Println(self.readServer())
	_, err = self.Connection.Write([]byte(`{"type":"auth","data":{"user_id":"` + self.BotId + `","token":"` + self.Token + `"}}`))
	if err != nil {
		fmt.Println("Ошибка во время отправки логина: ", err)
		time.Sleep(10 * time.Second)
	}
	fmt.Println(self.readServer())
}

func (self *GoodGameBot) joinChannels() error {
	var err error
	for _, channel := range self.Channels {
		_, err = self.Connection.Write([]byte(`{"type":"join","data":{"channel_id":"` + channel + `","hidden":false}}`))
		fmt.Println(self.readServer())
		if err != nil {
			fmt.Println("Ошибка во время входа в чат-комнату: ", err)
			return err
		}
	}
	return nil
}

func (self *GoodGameBot) Start() {
	var err error
	self.initBot()
	self.serverResponse = make([]byte, 1024)
	for {
		self.connect()
		err = self.joinChannels()
		if err != nil {
			err = nil
			continue
		}
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

func (self *GoodGameBot) listenChannels() error {
	var err error
	for {
		if err = self.handleChat(); err != nil {
			return err
		}
	}
}

func (self *GoodGameBot) say(msg, channel string) {
	if msg == "" {
		fmt.Println("Сообщение пустое")
	}
	self.Connection.Write([]byte(`{"type":"send_message","data":{"channel_id":"` + channel + `","text":"` + msg + `","hideIcon":false,"mobile":false}}`))
	fmt.Println(self.readServer())
}
