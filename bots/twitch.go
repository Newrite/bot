package bots

import (
	"golang.org/x/net/websocket"
	"net"
	"net/textproto"
	"os"
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
	uptime              int64
	FileChannelLog      map[string]*os.File
	Settings            map[string]*botSettings
}

type botSettings struct {
	Status        bool
	ReactStatus   bool
	CMDStatus     bool
	ReactRate     int
	LastReactTime int64
	IsModerator   bool
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
