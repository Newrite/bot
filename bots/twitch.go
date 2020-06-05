package bots

import (
	"golang.org/x/net/websocket"
	"net"
	"net/textproto"
	"os"
	"time"
)

type BotTwitch struct {
	BotName             string   `json:"bot_name"`
	OAuth               string   `json:"o_auth"`
	Server              string   `json:"server"`
	Port                string   `json:"port"`
	OwnerBot            string   `json:"owner_bot"`
	Channels            []string `json:"Channels"`
	Connection          net.Conn
	WebSocketConnection *websocket.Conn
	ReadChannels        *textproto.Reader
	ApiConf             *apiConfig
	GoodGameBotPtr      *BotGoodGame
	DiscordPtr          *BotDiscord
	MutedUsers          string
	serverResponse      []byte
	n                   int
	FileChannelLog      map[string]*os.File
	Settings            map[string]*botSettings
	Viewers             map[string]*viewersData
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
