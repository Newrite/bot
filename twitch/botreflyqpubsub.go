package twitch

import (
	"fmt"
	"golang.org/x/net/websocket"
	"strings"
	"time"
)

const URL string = "wss://pubsub-edge.twitch.tv/"
const ORIGIN string = "https://pubsub-edge.twitch.tv/"

func (bt *BotTwitch) startPubSub() {
	var err error
	bt.WebSocketConnection, err = websocket.Dial(URL, "", ORIGIN)
	if err != nil {
		fmt.Println(err)
	}
	go bt.pingSub()
	bt.startListenChannelPoints()
	bt.serverResponse = make([]byte, 4096)
	bt.listenPubSub()
}

func (bt *BotTwitch) pingSub() {
	var err error
	for {
		_, err = bt.WebSocketConnection.Write([]byte(`{"type":"PING"}`))
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(5 * time.Minute)
	}
}

func (bt *BotTwitch) startListenChannelPoints() {
	var err error
	_, err = bt.WebSocketConnection.Write([]byte(`{"type":"LISTEN","nonce":"qweprotiyunb","data":{"topics":["channel-points-channel-v1.` + bt.ApiConf.ChannelsID[channelRflyq] + `"],"auth_token":"` + bt.ApiConf.ReflyToken + `"}}`))
	if err != nil {
		fmt.Println(err)
	}
}

func (bt *BotTwitch) listenPubSub() {
	var err error
	for {
		bt.n, err = bt.WebSocketConnection.Read(bt.serverResponse)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Recived: %s \n", bt.serverResponse[:bt.n])
		if strings.Contains(string(bt.serverResponse[:bt.n]), "\"topic\":\"channel-points-channel-v1") {
			var username, message, cmd string = bt.handlePubSub(string(bt.serverResponse[:bt.n]))
			fmt.Println(username, message, cmd)
		}
		time.Sleep(1 * time.Second)
	}
}

func (bt *BotTwitch) handlePubSub(body string) (username, message, cmd string) {
	body = strings.Replace(body, "\"", " ", -1)
	body = strings.Replace(body, ":", " ", -1)
	body = strings.Replace(body, "\\", " ", -1)
	body = strings.Replace(body, "}", " ", -1)
	body = strings.Replace(body, "{", " ", -1)
	body = strings.Replace(body, ",", " ", -1)
	var check int
	bodySlice := strings.Fields(body)
	for id, sl := range bodySlice {
		if check == 1 && sl != "prompt" {
			cmd += " " + sl
		}
		if sl == "title" {
			cmd = bodySlice[id+1]
			check = 1
		}
		if sl == "prompt" {
			check = 2
		}
		if check == 3 && id < len(bodySlice)-2 && bodySlice[id-1] != "prompt" {
			message += " " + sl
		}
		if sl == "user_input" {
			message = bodySlice[id+1]
			check = 3
		}
		//fmt.Println("ID:", id, "Text: ", sl)
	}
	username = bodySlice[20]
	return username, message, cmd
}
