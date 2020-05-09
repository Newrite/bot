package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/textproto"
	"time"
	"strings"
)

type TwitchBot struct {
	BotName    string   `json:"bot_name"`
	OAuth      string   `json:"o_auth"`
	Server     string   `json:"server"`
	Port       string   `json:"port"`
	OwnerBot   string   `json:"owner_bot"`
	Channels   []string `json:"channels"`
	Connection net.Conn
	ReadChannels *textproto.Reader
}

func (bot *TwitchBot) connect() {
	var err error
	bot.Connection, err = net.Dial("tcp", bot.Server+":"+bot.Port)
	if err != nil {
		fmt.Print("Ошибка попытки соединения: ", err)
		time.Sleep(300000000)
		bot.connect()
	}
	_, err = bot.Connection.Write([]byte("PASS " + bot.OAuth + "\r\n"))
	if err != nil {
		fmt.Print("Ошибка во время отправки oauth-key: ", err)
		time.Sleep(300000000)
		bot.connect()
	}
	_, err = bot.Connection.Write([]byte("NICK " + bot.BotName + "\r\n"))
	if err != nil {
		fmt.Print("Ошибка во время отправки логина: ", err)
		time.Sleep(300000000)
		bot.connect()
	}
}

func (bot *TwitchBot) joinChannels() {
	for _, channel := range bot.Channels {
		_, err := bot.Connection.Write([]byte("JOIN #" + channel + "\r\n"))
		if err != nil {
			fmt.Print("Ошибка во время входа в чат-комнату: ", err)
		}
	}
}

func (bot *TwitchBot) listenChannels() {
	for {
		line, err := bot.ReadChannels.ReadLine()
		if err != nil {
			fmt.Println("Error: ", err)
		}
		fmt.Printf(line + "\n")
		if line == "PING :tmi.twitch.tv" {
			bot.Connection.Write([]byte("PONG\r\n"))
		}
		var userName, channelName, message string = bot.handleLine(line)
		//fmt.Println("DEBUG 1: ", "User: ", userName, " Channel: ", channelName, " MSG: ", message)
		if strings.Contains(message, "!бот") || strings.Contains(message, "!бот") {
			bot.say("@"+userName+" AdaIsEva, написана на GoLang v1.14 без использования сторонних библиотек.", channelName)
		}
		time.Sleep(1000000)
	}
}

func (bot *TwitchBot) say(msg, channel string) {
	if msg == "" {
		fmt.Println("Сообщение пустое")
		return
	}
	_, err := bot.Connection.Write([]byte("PRIVMSG #"+channel+" :"+msg+"\r\n"))
	if err != nil {
		fmt.Println("Ошибка отрпавки сообщения: ", err)
	}
}

func (bot *TwitchBot) handleLine(line  string) (user, channel, message string) {
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
	botFile, err := ioutil.ReadFile("BotData.json")
	if err != nil {
		fmt.Print("Ошибка чтения данных бота (BotData.Json), должно находиться в корневой папке с исполняемым файлом: ", err)
	}
	err = json.Unmarshal(botFile, &bot)
	if err != nil {
		fmt.Print("Ошибка конвертирования структуры из файла в структуру бота: ", err)
	}
	fmt.Println(bot)
	bot.connect()
	bot.joinChannels()
	bot.ReadChannels = textproto.NewReader(bufio.NewReader(bot.Connection))
	go bot.listenChannels()
	for {
		time.Sleep(100000000)
	}
	defer bot.Connection.Close()
}
