package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (self *TwitchBot) templateRequest(method, url, headAuth string) []byte {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	req.Header.Set("Accept", "application/vnd.twitchtv.v5+json")
	req.Header.Set("Authorization", headAuth)
	req.Header.Add("Client-ID", self.ApiConf.Client_id)
	if err != nil {
		panic(err.Error())
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	return body
}

func (self *TwitchBot) initApiConfig() {
	botConfig, err := ioutil.ReadFile("BotApiConfig.json")
	if err != nil {
		fmt.Print("Ошибка чтения данных апи для бота (BotApiConfig.json),"+
			" должно находиться в корневой папке с исполняемым файлом: ", err)
	}
	err = json.Unmarshal(botConfig, &self.ApiConf)
	if err != nil {
		fmt.Print("Ошибка конвертирования структуры из файла в апи структуру бота: ", err)
	}
	self.requestInitStreamersID()
}

func (self TwitchBot) requestInitStreamersID() {
	var users usersData
	self.ApiConf.Url = "https://api.twitch.tv/helix/users?login=" + self.Channels[0]
	self.ApiConf.ChannelsID = make(map[string]string)
	for _, channel := range self.Channels {
		self.ApiConf.ChannelsID[channel] = ""
		if channel != self.Channels[0] {
			self.ApiConf.Url += "&login=" + channel
		}
	}
	body := self.templateRequest("GET", self.ApiConf.Url, self.ApiConf.Bearer)
	json.Unmarshal(body, &users)
	for key, _ := range self.ApiConf.ChannelsID {
		for _, channel := range users.User {
			if channel.Login == key {
				self.ApiConf.ChannelsID[key] = channel.Id
			}
		}
	}
}

func (self *TwitchBot) handleApiRequest(username, channel, message, cmd string) string {
	switch cmd {
	case "!вырубай":
		if self.requestChatterData(channel, username, "mod") == "mod" {
			return "Моё уважение модераторскому корпусу, но нет roflanZdarova"
		}
		if self.requestBroadcasterSubscriptionsData(channel, username, "subus") == "Саб" {
			if self.requestChatterData(channel, username, "vip") == "vip" {
				return "Можно пожалуйста постримить? PepeHands"
			} else {
				return "Зачем ты это делаешь? roflanZachto"
			}
		} else {
			if self.requestChatterData(channel, username, "vip") == "vip" {
				return "Ты ходишь по тонкому льду, випчик.. Ладно живи roflanEbalo"
			} else {
				return "unsub"
			}
		}
	case "!evaismod":
		if self.requestChatterData(channel, self.BotName, "mod") == "mod" {
			return "true"
		} else {
			return "false"
		}
	default:
		return "error"
	}
}

func (self *TwitchBot) requestChatterData(channel, username, cmd string) string {
	var chatters chattersData
	self.ApiConf.Url = "https://tmi.twitch.tv/group/user/" + channel + "/chatters"
	body := self.templateRequest("GET", self.ApiConf.Url, self.ApiConf.O_Auth)
	json.Unmarshal(body, &chatters)
	switch cmd {
	case "vip":
		for _, name := range chatters.Chatters.Vips {
			if name == username {
				return "vip"
			}
		}
	case "mod":
		for _, name := range chatters.Chatters.Moderators {
			if name == username {
				return "mod"
			}
		}
	}
	return "Nothing"
}

func (self TwitchBot) requestUsersData(channel, username, cmd string) string {
	var users usersData
	self.ApiConf.Url = "https://api.twitch.tv/helix/users?login=" + channel
	body := self.templateRequest("GET", self.ApiConf.Url, self.ApiConf.Bearer)
	json.Unmarshal(body, &users)
	if len(users.User) < 1 {
		fmt.Println("requestUsersData аут оф аррей")
		return ""
	}
	switch cmd {
	case "DisNam":
		return users.User[0].Display_name
	default:
		return "Nothing"
	}
}

func (self TwitchBot) requestBroadcasterSubscriptionsData(channel, username, cmd string) string {
	var broadcasterSubscriptions broadcasterSubscriptionsData
	self.ApiConf.Url = "https://api.twitch.tv/helix/subscriptions?broadcaster_id=" + self.ApiConf.ChannelsID[channel]
	body := self.templateRequest("GET", self.ApiConf.Url, self.ApiConf.Bearer)
	json.Unmarshal(body, &broadcasterSubscriptions)
	switch cmd {
	case "subus":
		username = self.requestUsersData(channel, username, "DisNam")
		for _, name := range broadcasterSubscriptions.Subscriptions {
			if name.User_name == username {
				return "Саб"
			}
		}
		return "Не саб"
	default:
		return "Nothing"
	}
}
