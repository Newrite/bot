package bots

import (
	"bot/resource"
	"fmt"
	"strings"
	"time"
)

func (bgg *BotGoodGame) handleChat() error {
	response := bgg.readServer()
	fmt.Println(response)
	if !strings.Contains(response, "\"type\":\"message\"") {
		return nil
	}
	var userID, userName, channel, message string = bgg.handleLine(response)
	fmt.Print("[" + timeStamp() + "] [GOODGAME] Канал:" + channel + " " +
		"Ник:" + userName + "\tСообщение:" + message + "\n")
	if !strings.HasPrefix(message, GgPrefix) {
		SingleTwitch().MarkovChain += " " + resource.ReadTxt(message)
	}
	fmt.Println("ID:", strings.TrimSpace(userID))
	bgg.checkReact(channel, message)
	lowMessage := strings.ToLower(message)
	if strings.HasPrefix(lowMessage, GgPrefix) {
		msgSl := strings.Fields(lowMessage)
		go bgg.say(checkCMD(userName, channel, msgSl[0], GG, lowMessage, message, strings.TrimSpace(userID)), channel)
	}
	if strings.HasPrefix(lowMessage, GgPrefix) {
		msgSl := strings.Fields(lowMessage)
		go bgg.say(handleCMDfromDB(userName, "xandr_sh", strings.TrimPrefix(msgSl[0], GgPrefix)), channel)
	}
	time.Sleep(10 * time.Millisecond)
	return nil
}

func (bgg *BotGoodGame) checkReact(channel, message string) {
	for key := range reactGG {
		if strings.Contains(message, key) {
			bgg.say(reactGG[key], channel)
			break
		}
	}
}

func (bgg *BotGoodGame) handleLine(line string) (userID, user, channel, message string) {
	line = strings.Replace(line, "\"", " ", -1)
	lineSlice := strings.Fields(line)
	var tempId int
	for id, field := range lineSlice {
		if field == "text" {
			tempId = id
		}
		if tempId != 0 && id > tempId+1 && id != len(lineSlice)-1 {
			if id != len(lineSlice)-2 {
				message += field + " "
			} else {
				message += field
			}
		}
	}
	return strings.TrimSuffix(strings.TrimPrefix(lineSlice[12], `:`), `,`), lineSlice[15], lineSlice[9], message
}
