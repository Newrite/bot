package bots

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func (bt *BotTwitch) templateRequest(method, url, headAuth string) []byte {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if req == nil || err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "NewRequest",
			"file":     "twitchapi.go",
			"body":     "templateRequest",
			"error":    err,
			"params":   method + " " + url,
		}).Errorln("Request == nil.")
		return nil
	}
	req.Header.Set("Accept", "application/vnd.twitchtv.v5+json")
	req.Header.Add("Authorization", headAuth)
	req.Header.Add("Client-ID", bt.ApiConf.Client_id)
	resp, err := client.Do(req)
	if err != nil || resp == nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "client.Do",
			"file":     "twitchapi.go",
			"body":     "templateRequest",
			"error":    err,
			"Request":  req,
		}).Errorln("Ошибка обработки реквеста.")
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "ioutil.ReadAll",
			"file":     "twitchapi.go",
			"body":     "templateRequest",
			"error":    err,
			"Request":  req,
		}).Errorln("Ошибка парсинга тела реквеста в срез байтов.")
		return nil
	}
	return body
}

func (bt *BotTwitch) initApiConfig() {
	botConfig, err := ioutil.ReadFile("BotApiConfig.json")
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "ioutil.ReadAll",
			"file":     "twitchapi.go",
			"body":     "initApiConfig",
			"error":    err,
		}).Errorln("Ошибка чтения данных апи для бота (BotApiConfig.json), должно находиться" +
			" в корневой папке с исполняемым файлом.")
	}
	err = json.Unmarshal(botConfig, &bt.ApiConf)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "json.Unmarshal",
			"file":     "twitchapi.go",
			"body":     "initApiConfig",
			"error":    err,
		}).Errorln("Ошибка конвертирования структуры из файла в апи структуру бота.")
	}
	bt.requestInitStreamersID()
}

func (bt *BotTwitch) requestInitStreamersID() {
	var users usersData
	bt.ApiConf.Url = "https://api.Twitch.tv/helix/Users?login=" + bt.Channels[0]
	bt.ApiConf.ChannelsID = make(map[string]string)
	for _, channel := range bt.Channels {
		bt.ApiConf.ChannelsID[channel] = ""
		if channel != bt.Channels[0] {
			bt.ApiConf.Url += "&login=" + channel
		}
	}
	body := bt.templateRequest("GET", bt.ApiConf.Url, bt.ApiConf.Bearer)
	err := json.Unmarshal(body, &users)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "json.Unmarshal",
			"file":     "twitchapi.go",
			"body":     "requestInitStreamersID",
			"error":    err,
		}).Errorln("Ошибка конвертирования структуры.")
	}
	for key, _ := range bt.ApiConf.ChannelsID {
		for _, channel := range users.User {
			if channel.Login == key {
				bt.ApiConf.ChannelsID[key] = channel.Id
			}
		}
	}
}

func (bt *BotTwitch) handleApiRequest(username, channel, message, cmd string) string {
	switch cmd {
	case "!вырубай":
		if bt.requestChatterData(channel, username, "mod") == "mod" {
			return "Моё уважение модераторскому корпусу, но нет roflanZdarova"
		}
		if bt.requestBroadcasterSubscriptionsData(channel, username, "subus") == "Саб" {
			if bt.requestChatterData(channel, username, "vip") == "vip" {
				return "Можно пожалуйста постримить? PepeHands"
			} else {
				return "Зачем ты это делаешь? roflanZachto"
			}
		} else {
			if bt.requestChatterData(channel, username, "vip") == "vip" {
				return "Ты ходишь по тонкому льду, випчик.. Ладно живи roflanEbalo"
			} else {
				return "unsub"
			}
		}
	case "userstate":
		if bt.requestChatterData(channel, username, "mod") == "mod" {
			return "mod"
		}
		if bt.requestBroadcasterSubscriptionsData(channel, username, "subus") == "Саб" {
			if bt.requestChatterData(channel, username, "vip") == "vip" {
				return "subvip"
			} else {
				return "sub"
			}
		} else {
			if bt.requestChatterData(channel, username, "vip") == "vip" {
				return "vip"
			} else {
				return "unsub"
			}
		}
	case "!evaismod":
		if bt.requestChatterData(channel, bt.BotName, "mod") == "mod" {
			return "true"
		} else {
			return "false"
		}
	case "streamStatus":
		return bt.requestStreamData(channel, username, cmd)
	case "uptime":
		return bt.requestStreamData(channel, username, cmd)
	default:
		return "error"
	}
}

func (bt *BotTwitch) requestChatterData(channel, username, cmd string) string {
	var chatters chattersData
	bt.ApiConf.Url = "https://tmi.Twitch.tv/group/user/" + channel + "/chatters"
	body := bt.templateRequest("GET", bt.ApiConf.Url, bt.ApiConf.O_Auth)
	err := json.Unmarshal(body, &chatters)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "json.Unmarshal",
			"file":     "twitchapi.go",
			"body":     "requestChatterData",
			"error":    err,
		}).Errorln("Ошибка конвертирования структуры.")
	}
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

func (bt *BotTwitch) requestUsersData(channel, username, cmd string) string {
	var users usersData
	bt.ApiConf.Url = "https://api.Twitch.tv/helix/Users?login=" + channel
	body := bt.templateRequest("GET", bt.ApiConf.Url, bt.ApiConf.Bearer)
	err := json.Unmarshal(body, &users)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "json.Unmarshal",
			"file":     "twitchapi.go",
			"body":     "requestUsersData",
			"error":    err,
		}).Errorln("Ошибка конвертирования структуры.")
	}
	if len(users.User) < 1 {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "len(users.User) < 1 ",
			"file":     "twitchapi.go",
			"body":     "requestUsersData",
			"error":    err,
		}).Errorln("Пустой срез юзверей.")
		return ""
	}
	switch cmd {
	case "DisNam":
		return users.User[0].Display_name
	default:
		return "Nothing"
	}
}

func (bt *BotTwitch) requestBroadcasterSubscriptionsData(channel, username, cmd string) string {
	var broadcasterSubscriptions broadcasterSubscriptionsData
	bt.ApiConf.Url = "https://api.Twitch.tv/helix/subscriptions?broadcaster_id=" + bt.ApiConf.ChannelsID[channel]
	body := bt.templateRequest("GET", bt.ApiConf.Url, bt.ApiConf.Bearer)
	err := json.Unmarshal(body, &broadcasterSubscriptions)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "json.Unmarshal",
			"file":     "twitchapi.go",
			"body":     "requestBroadcasterSubscriptionsData",
			"error":    err,
		}).Errorln("Ошибка конвертирования структуры.")
	}
	switch cmd {
	case "subus":
		username = bt.requestUsersData(username, channel, "DisNam")
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

func (bt *BotTwitch) requestStreamData(channel, username, cmd string) string {
	var stream streamData
	bt.ApiConf.Url = "https://api.Twitch.tv/helix/streams?user_login=" + channel
	body := bt.templateRequest("GET", bt.ApiConf.Url, bt.ApiConf.Bearer)
	err := json.Unmarshal(body, &stream)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "json.Unmarshal",
			"file":     "twitchapi.go",
			"body":     "requestStreamData",
			"error":    err,
		}).Errorln("Ошибка конвертирования структуры.")
	}
	if len(stream.Data) < 1 {
		return "offline"
	}
	switch cmd {
	case "streamStatus":
		return "online"
	case "uptime":
		twitchParser, _ := time.Parse(time.RFC3339, stream.Data[0].Started_at)
		tempstr := strings.Replace(time.Since(twitchParser).Truncate(time.Second).String(), "h", "ч", -1)
		tempstr = strings.Replace(tempstr, "m", "м", -1)
		tempstr = strings.Replace(tempstr, "s", "с", -1)
		return " " + tempstr
	default:
		return ""
	}
}
