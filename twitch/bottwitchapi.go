package twitch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func (bt *BotTwitch) templateRequest(method, url, headAuth string) []byte {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	req.Header.Set("Accept", "application/vnd.twitchtv.v5+json")
	req.Header.Add("Authorization", headAuth)
	req.Header.Add("Client-ID", bt.ApiConf.Client_id)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	return body
}

func (bt *BotTwitch) initApiConfig() {
	botConfig, err := ioutil.ReadFile("BotApiConfig.json")
	if err != nil {
		fmt.Print("Ошибка чтения данных апи для бота (BotApiConfig.json),"+
			" должно находиться в корневой папке с исполняемым файлом: ", err)
	}
	err = json.Unmarshal(botConfig, &bt.ApiConf)
	if err != nil {
		fmt.Print("Ошибка конвертирования структуры из файла в апи структуру бота: ", err)
	}
	bt.requestInitStreamersID()
}

func (bt *BotTwitch) requestInitStreamersID() {
	var users usersData
	bt.ApiConf.Url = "https://api.twitch.tv/helix/users?login=" + bt.Channels[0]
	bt.ApiConf.ChannelsID = make(map[string]string)
	for _, channel := range bt.Channels {
		bt.ApiConf.ChannelsID[channel] = ""
		if channel != bt.Channels[0] {
			bt.ApiConf.Url += "&login=" + channel
		}
	}
	body := bt.templateRequest("GET", bt.ApiConf.Url, bt.ApiConf.Bearer)
	json.Unmarshal(body, &users)
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
	case "requestallviewers":
		bt.requestChatterData(channel, "", "initviewers")
		return ""
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
	bt.ApiConf.Url = "https://tmi.twitch.tv/group/user/" + channel + "/chatters"
	body := bt.templateRequest("GET", bt.ApiConf.Url, bt.ApiConf.O_Auth)
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
	case "initviewers":
		tempstr := make([]string, 0, 0)
		for _, name := range chatters.Chatters.Viewers {
			tempstr = append(tempstr, name)
		}
		for _, name := range chatters.Chatters.Moderators {
			tempstr = append(tempstr, name)
		}
		for _, name := range chatters.Chatters.Global_mods {
			tempstr = append(tempstr, name)
		}
		for _, name := range chatters.Chatters.Admins {
			tempstr = append(tempstr, name)
		}
		for _, name := range chatters.Chatters.Vips {
			tempstr = append(tempstr, name)
		}
		for _, name := range chatters.Chatters.Broadcaster {
			tempstr = append(tempstr, name)
		}
		for _, name := range chatters.Chatters.Staff {
			tempstr = append(tempstr, name)
		}
		for _, name := range tempstr {
			if bt.checkToAddViewer(name, channel) {
				tempstruct := &viewer{
					Name:   name,
					Points: 0,
				}
				bt.Viewers[channel].Viewers = append(bt.Viewers[channel].Viewers, tempstruct)
			}
		}
	}
	return "Nothing"
}

func (bt *BotTwitch) checkToAddViewer(name, channel string) bool {
	for _, viewer := range bt.Viewers[channel].Viewers {
		if name == viewer.Name {
			return false
		}
	}
	return true
}

func (bt *BotTwitch) requestUsersData(channel, username, cmd string) string {
	var users usersData
	bt.ApiConf.Url = "https://api.twitch.tv/helix/users?login=" + channel
	body := bt.templateRequest("GET", bt.ApiConf.Url, bt.ApiConf.Bearer)
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

func (bt *BotTwitch) requestBroadcasterSubscriptionsData(channel, username, cmd string) string {
	var broadcasterSubscriptions broadcasterSubscriptionsData
	bt.ApiConf.Url = "https://api.twitch.tv/helix/subscriptions?broadcaster_id=" + bt.ApiConf.ChannelsID[channel]
	body := bt.templateRequest("GET", bt.ApiConf.Url, bt.ApiConf.Bearer)
	json.Unmarshal(body, &broadcasterSubscriptions)
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
	bt.ApiConf.Url = "https://api.twitch.tv/helix/streams?user_login=" + channel
	body := bt.templateRequest("GET", bt.ApiConf.Url, bt.ApiConf.Bearer)
	json.Unmarshal(body, &stream)
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
