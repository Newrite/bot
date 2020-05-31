package goodgame

import (
	"bot/reso"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func (bgg *BotGoodGame) handleChat() error {
	response := bgg.readServer()
	fmt.Println(response)
	if !strings.Contains(response, "\"type\":\"message\"") {
		return nil
	}
	var username, channel, message string = bgg.handleLine(response)
	fmt.Println(username, channel, message)
	bgg.checkReact(channel, message)
	message = strings.ToLower(message)
	bgg.checkCMD(channel, username, message)
	time.Sleep(10 * time.Millisecond)
	return nil
}

func (bgg *BotGoodGame) checkReact(channel, message string) {
	for key := range react {
		if strings.Contains(message, key) {
			bgg.say(react[key], channel)
			break
		}
	}
}

func (bgg *BotGoodGame) checkCMD(channel, userName, message string) {
	for key, value := range cmd {
		if strings.HasPrefix(message, key) && value != "_" {
			bgg.say(userName+", "+cmd[key], channel)
			break
		}
		if strings.HasPrefix(message, key) && value == "_" {
			go bgg.say(bgg.handleInteractiveCMD(key, channel, userName, message), channel)
			break
		}
	}
}

func (bgg *BotGoodGame) handleInteractiveCMD(cmd, channel, userName, message string) string {
	switch cmd {
	case "!roll":
		return userName + ", " + strconv.Itoa(rand.Intn(21))
	case "!eva":
		return reso.EvaAnswers[rand.Intn(16)]
	case "!билд":
		return reso.BuildAnswers[rand.Intn(16)]
	case "!uptime":
		return bgg.TwitchPtr.HandleRequests("uptime")
	default:
		return ""
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
