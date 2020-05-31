package goodgame

import (
	"bot/twitch"
	"golang.org/x/net/websocket"
	"time"
)

const TimeFormat = "2006.01.02 15:04"
const TimeFormatReact = "2006.01.02 15:04:02"

func timeStamp() string {
	return time.Now().Format(TimeFormat)
}

var cmd = map[string]string{
	"!ping": "pong!",
	"!бот": "AdaIsEva, чат-бот для GG и twitch, написана на GoLang v1.14 без использования сторонних библиотек. " +
		"Живет на VPS с убунтой размещенном в москоу сити. Рекомендации, пожелания и" +
		" прочая можно присылать на adaiseva.newrite@gmail.com",
	"!help": "Доступные комманды: !ping, !бот, !roll, !help, !Eva",
	//"!master help": "Владелец бота либо канала может переключить активность бота коммандой !Ada, switch." +
	//	" Реакции на всякое разное командой !Ada, switch react. " +
	//	"Переключить отзыв на различные команды !Ada, switch cmd." +
	//	" !Ada, set reactrate to <значение> выставляет настройку частоты реакции на различные сообщения в чате",
	"!roll":             "_",
	"!вырубай":          "_",
	"!eva":              "_",
	"!билд":             "_",
	"!вырубить":         "_",
	"!чезаигра":         "_",
	"!скиллуха":         "_",
	"!вызватьсанитаров": "_",
	"!труба":            "_",
	"!тыктоблять":       "_",
	"!дискорд":          "_",
	"!uptime":           "_",
}

var react = map[string]string{
	"Привет":  ":skull: 10000010 10010000 11000010 00100000 01000011 11101000 01100101 00001100",
	"привет":  ":skull: 10000010 10010000 11000010 00100000 01000011 11101000 01100101 00001100",
	"Hello":   ":skull: 10000010 10010000 11000010 00100000 01000011 11101000 01100101 00001100",
	"hello":   ":skull: 10000010 10010000 11000010 00100000 01000011 11101000 01100101 00001100",
	"+ в чат": "+",
}

type BotGoodGame struct {
	BotName        string   `json:"bot_name"`
	Token          string   `json:"token"`
	BotId          string   `json:"bot_id"`
	Server         string   `json:"server"`
	Origin         string   `json:"origin"`
	OwnerBot       string   `json:"owner_bot"`
	Channels       []string `json:"channels"`
	TwitchPtr      *twitch.BotTwitch
	Connection     *websocket.Conn
	serverResponse []byte
	n              int
}
