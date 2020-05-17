package TwitchAPI

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const clientId string = "gp762nuuoqcoxypju8c569th9wz7q5"
const oauth string = "Bearer 2m0hfirzasam3tkwo0ty9wwf1wjmz3"
const bearer string = "Bearer 2m0hfirzasam3tkwo0ty9wwf1wjmz3"
const clientIdRef string = "cuj07xgrz9hv2rfq2gev4lkvlu3tlg"
const oauthRef string = "Bearer 6jcz1sstmyh83pzjn5rdjyxsssyl00"
const reflyID string = "54987522"

type chattersData struct {
	Chatters struct {
		Vips []string `json:"vips"`
		Moderators []string `json:"moderators"`
	} `json:"chatters"`
}

type usersData struct {
	User []userData `json:"data"`
}

type userData struct {
	Id                string `json:"id"`
	Login             string `json:"login"`
	Display_name      string `json:"display_name"`
	Type              string `json:"type"`
	Broadcaster_type  string `json:"broadcaster_type"`
	Description       string `json:"description"`
	Profile_image_url string `json:"profile_image_url"`
	Offline_image_url string `json:"offline_image_url"`
	View_count        int    `json:"view_count"`
}

type channelData struct {
	Id                              int    `json:"id"`
	Broadcaster_language            string `json:"broadcaster_language"`
	Created_at                      string `json:"created_at"`
	Display_name                    string `json:"Display_name"`
	Followers                       int    `json:"followers"`
	Game                            string `json:"game"`
	Language                        string `json:"language"`
	Logo                            string `json:"logo"`
	Mature                          bool   `json:"mature"`
	Name                            string `json:"name"`
	Partner                         bool   `json:"partner"`
	Profile_banner                  bool   `json:"profile_banner"`
	Profile_banner_background_color bool   `json:"profile_banner_background_color"`
	Status                          string `json:"status"`
	Updated_at                      string `json:"updated_at"`
	Url                             string `json:"url"`
	Video_banner                    bool   `json"video_banner"`
	Views                           int    `json:"views"`
}

type streamsData struct {
	Stream []streamData `json:"data"`
}

type streamData struct {
	Started_at time.Time `json:"started_at"`
}

type subsData struct {
	SubData []subData `json:"data"`
}

type subData struct {
	User_name string `json:"user_name"`
}

func subsDataParse(username string, Data *usersData) string {
	client := &http.Client{}
	url := "https://api.twitch.tv/helix/subscriptions?broadcaster_id="+Data.User[0].Id
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/vnd.twitchtv.v5+json")
	req.Header.Set("Authorization", oauthRef)
	req.Header.Add("Client-ID", clientIdRef)
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
	var data subsData
	err = json.Unmarshal(body, &data)
	username = GOTwitch(username, "realname", username)
	if err != nil {
		return "Ошибка парсинга в json streamsData"
	}
	for _, name := range data.SubData {
		if name.User_name == username {
			return "Саб"
		}
	}
	return "Не саб"
}

func channelDataParse(data *usersData, cmd, username string) string {
	client := &http.Client{}
	if len(data.User) == 0 {
		return "Не удалось получить данные пользователей"
	}
	url := "https://api.twitch.tv/kraken/channels/" + data.User[0].Id
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/vnd.twitchtv.v5+json")
	req.Header.Add("Authorization", oauth)
	req.Header.Add("Client-ID", clientId)
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
	var dataChannel channelData
	json.Unmarshal(body, &dataChannel)
	fmt.Println(string(body))
	fmt.Println("Data:", dataChannel)
	switch cmd {
	case "game":
		return dataChannel.Game
	case "followers":
		return strconv.Itoa(dataChannel.Followers)
	case "status":
		return dataChannel.Status
	case "channelcreated":
		return dataChannel.Created_at
	}
	return "Ошибка"
}

func GOTwitch(channel, cmd, username string) string {
	client := &http.Client{}
	url := "https://api.twitch.tv/helix/users?login=" + channel
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/vnd.twitchtv.v5+json")
	req.Header.Add("Authorization", oauth)
	req.Header.Add("Client-ID", clientId)
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
	var data usersData
	json.Unmarshal(body, &data)
	switch cmd {
	case "game":
		return channelDataParse(&data, cmd, username)
	case "followers":
		return channelDataParse(&data, cmd, username)
	case "channelcreated":
		return channelDataParse(&data, cmd, username)
	case "status":
		return channelDataParse(&data, cmd, username)
	case "realname":
		if len(data.User) > 0 {
			return data.User[0].Display_name
		} else {
			return "Не удалось получить данные"
		}
	case "uptime":
		return streamDataParse(channel, cmd)
	case "sub":
		return subsDataParse(username, &data)
	case "mod":
		return moderOrVip(channel, cmd, username)
	case "reflysub":
		if moderOrVip(channel, "mod", username) == "mod" {
			return "Моё уважение модераторскому корпусу, но нет roflanZdarova"
		}
		if subsDataParse(username, &data) == "Саб" {
			if moderOrVip(channel, "vip", username) == "vip" {
				return "Можно пожалуйста постримить? PepeHands"
			} else {
				return "Зачем ты это делаешь? roflanZachto"
			}
		} else {
			if moderOrVip(channel, "vip", username) == "vip" {
				return "Ты ходишь по тонкому льду, випчик.. Ладно живи roflanEbalo"
			} else {
				return "unsub"
			}
		}
	}
	return "ничего"
}

func moderOrVip(channel, cmd, username string) string {
	client := &http.Client{}
	url := "https://tmi.twitch.tv/group/user/reflyq/chatters"
	req, err := http.NewRequest("GET", url, nil)
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
	var chat chattersData
	json.Unmarshal(body, &chat)
	switch cmd {
	case "vip":
		for _, name := range chat.Chatters.Vips {
			if name == username {
				return "vip"
			}
		}
	case "mod":
		for _, name := range chat.Chatters.Moderators {
			if name == username {
				return "mod"
			}
		}
	}
	return "nothing"
}

func streamDataParse(channel, cmd string) string {
	client := &http.Client{}
	url := "https://api.twitch.tv/helix/streams?user_login=" + channel
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/vnd.twitchtv.v5+json")
	req.Header.Add("Authorization", oauth)
	req.Header.Add("Client-ID", clientId)
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
	var data streamsData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "Ошибка парсинга в json streamsData"
	}
	if len(data.Stream) == 0 {
		return "Стрим офлайн"
	}
	switch cmd {
	case "uptime":
		timeSince := time.Since(data.Stream[0].Started_at)
		sinceSplit := strings.Split(timeSince.String(), ".")
		return sinceSplit[0]
	default:
		return "ничего"
	}
}
