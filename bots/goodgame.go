package bots

import (
	"golang.org/x/net/websocket"
)

type BotGoodGame struct {
	BotName        string   `json:"bot_name"`
	Token          string   `json:"token"`
	BotId          string   `json:"bot_id"`
	Server         string   `json:"server"`
	Origin         string   `json:"origin"`
	OwnerBot       string   `json:"owner_bot"`
	Channels       []string `json:"Channels"`
	uptime         int64
	TwitchPtr      *BotTwitch
	DiscordPtr     *BotDiscord
	Connection     *websocket.Conn
	serverResponse []byte
	n              int
}
