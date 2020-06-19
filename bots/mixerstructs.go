package bots

import "github.com/gorilla/websocket"

type mixerConfig struct {
	Client_ID     string `json:"client_id"`
	Access_token  string `json:"access_token"`
	Refresh_token string `json:"refresh_token"`
	Client_secret string `json:"client_secret"`
	Expires_in    int    `json:"expires_in"`
}

type channelID struct {
	ID int `json:"id"`
}

type replayAuth struct {
	Type  string `json:"type"`
	Error string `json:"error"`
	ID    int    `json:"id"`
	Data  struct {
		Authenticated bool     `json:"authenticated"`
		Roles         []string `json:"roles"`
	} `json:"data"`
}

type chatMessage struct {
	Type  string `json:"type"`
	Error string `json:"error"`
	ID    int    `json:"id"`
	Data  struct {
		Channel     int      `json:"channel"`
		ID          string   `json:"id"`
		User_name   string   `json:"user_name"`
		User_ID     int      `json:"user_id"`
		User_level  int      `json:"user_level"`
		User_avatar string   `json:"user_avatar"`
		User_roles  []string `json:"user_roles"`
		Message     struct {
			Message []struct {
				Type string `json:"type"`
				Data string `json:"data"`
				Text string `json:"text"`
			} `json:"message"`
			Meta struct{} `json:"meta"`
		} `json:"message"`
	} `json:"data"`
}

type chatConnect struct {
	Roles       []string `json:"roles"`
	Authkey     string   `json:"authkey"`
	Permissions []string `json:"permissions"`
	Endpoints   []string `json:"endpoints"`
	IsLoadShed  bool     `json:"is_load_shed"`
}

type BotMixer struct {
	Bot_name   string `json:"bot_name"`
	Owner_bot  string `json:"owner_bot"`
	ApiConf    *mixerConfig
	Connection *websocket.Conn
	Message    *chatMessage
	uptime     int64
}
