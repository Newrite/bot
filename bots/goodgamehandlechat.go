package bots

import (
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

func (bgg *BotGoodGame) handleChat() error {
	response := bgg.readServer()
	if !strings.Contains(response, "\"type\":\"message\"") {
		return nil
	}
	var userName, channel, message string = bgg.handleLine(response)
	log.Infof("Ник: %s Канал: %s Сообщение: %s\n", userName, channel, message)
	bgg.checkReact(channel, message)
	message = strings.ToLower(message)
	if strings.HasPrefix(message, GgPrefix) {
		msgSl := strings.Fields(message)
		bgg.say(checkCMD(userName, channel, msgSl[0], "GG", message), channel)
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

func (bgg *BotGoodGame) handleLine(line string) (user, channel, message string) {
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
	return lineSlice[15], lineSlice[9], message
}
