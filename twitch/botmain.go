package twitch

import (
	"golang.org/x/net/websocket"
	//"github.com/sirupsen/logrus"
	"net"
	"net/textproto"
	"os"
	"time"
)

const TimeFormat = "2006.01.02 15:04"
const TimeFormatReact = "2006.01.02 15:04:02"

func timeStamp() string {
	return time.Now().Format(TimeFormat)
}

var cmd = map[string]string{
	"!ping": "pong!",
	"!бот": "AdaIsEva, написана на GoLang v1.14 без использования сторонних библиотек. " +
		"Живет на VPS с убунтой размещенном в москоу сити. Рекомендации, пожелания и" +
		" прочая можно присылать на adaiseva.newrite@gmail.com",
	"!help": "Доступные комманды: !ping, !бот, !roll, !help, !master help, !Eva, !uptime",
	"!master help": "Владелец бота либо канала может переключить активность бота коммандой !Ada, switch." +
		" Реакции на всякое разное командой !Ada, switch react. " +
		"Переключить отзыв на различные команды !Ada, switch cmd." +
		" !Ada, set reactrate to <значение> выставляет настройку частоты реакции на различные сообщения в чате",
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
	"револьвер выстреливает!": "_",
}

var react = map[string]string{
	"PogChamp": "PogChamp",
	"Kappa 7":  "Kappa 7",
	"Привет":   "MrDestructoid 10000010 10010000 11000010 00100000 01000011 11101000 01100101 00001100",
	"привет":   "MrDestructoid 10000010 10010000 11000010 00100000 01000011 11101000 01100101 00001100",
	"Hello":    "MrDestructoid 10000010 10010000 11000010 00100000 01000011 11101000 01100101 00001100",
	"hello":    "MrDestructoid 10000010 10010000 11000010 00100000 01000011 11101000 01100101 00001100",
	"SMOrc":    "SMOrc",
	"+ в чат":  "+",
}

type BotTwitch struct {
	BotName             string   `json:"bot_name"`
	OAuth               string   `json:"o_auth"`
	Server              string   `json:"server"`
	Port                string   `json:"port"`
	OwnerBot            string   `json:"owner_bot"`
	Channels            []string `json:"channels"`
	Connection          net.Conn
	WebSocketConnection *websocket.Conn
	ReadChannels        *textproto.Reader
	ApiConf             *apiConfig
	//GoodGameBotPtr      *goodgame.BotGoodGame
	MutedUsers     string
	serverResponse []byte
	n              int
	FileChannelLog map[string]*os.File
	Settings       map[string]*botSettings
	Viewers        map[string]*viewersData
}

type botSettings struct {
	Status      bool
	ReactStatus bool
	CMDStatus   bool
	ReactRate   time.Time
	ReactTime   int
	IsModerator bool
}

type apiConfig struct {
	Client_id  string `json:"client_id"`
	O_Auth     string `json:"o_auth"`
	Bearer     string `json:"bearer"`
	Secret_id  string `json:"secret_id"`
	ReflyToken string `json:"refly_token"`
	Url        string
	ChannelsID map[string]string
}

type viewersData struct {
	Viewers []*viewer
}

type viewer struct {
	Name   string
	Points int
}
