package goodgame

import (
	"bot/reso"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func (self *GoodGameBot) handleChat() error {
	response := self.readServer()
	fmt.Println(response)
	if !strings.Contains(response, "\"type\":\"message\"") {
		return nil
	}
	var username, channel, message string = self.handleLine(response)
	fmt.Println(username, channel, message)
	self.checkReact(channel, message)
	message = strings.ToLower(message)
	self.checkCMD(channel, username, message)
	time.Sleep(10 * time.Millisecond)
	return nil
}

func (self *GoodGameBot) checkReact(channel, message string) {
	for key := range react {
		if strings.Contains(message, key) {
			self.say(react[key], channel)
			break
		}
	}
}

func (self *GoodGameBot) checkCMD(channel, userName, message string) {
	for key, value := range cmd {
		if strings.HasPrefix(message, key) && value != "_" {
			self.say(userName+", "+cmd[key], channel)
			break
		}
		if strings.HasPrefix(message, key) && value == "_" {
			go self.say(self.handleInteractiveCMD(key, channel, userName, message), channel)
			break
		}
	}
}

func (self *GoodGameBot) handleInteractiveCMD(cmd, channel, userName, message string) string {
	switch cmd {
	case "!roll":
		return userName + ", " + strconv.Itoa(rand.Intn(21))
	case "!eva":
		return reso.EvaAnswers[rand.Intn(16)]
	case "!билд":
		return reso.BuildAnswers[rand.Intn(16)]
	default:
		return ""
	}
}

func (self *GoodGameBot) handleLine(line string) (user, channel, message string) {
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
