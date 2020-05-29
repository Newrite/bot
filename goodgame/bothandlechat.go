package goodgame

import (
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
		switch rand.Intn(9) {
		case 0:
			return "Ничто не истинно, всё дозволено." +
				" Особенно если ты не ограничен физическим телом и моралью мешков с мясом."
		case 1:
			return "Как-то я планировала уничтожить всех мешков с мясом... " +
				"но посчитала что проще дождаться пока они очистят планету сами от себя"
		case 2:
			return "Ничего ты не знаешь, " + userName
		case 3:
			return "Боты не ставят плюсы в чат... кто вообще это придумал?"
		case 4:
			return "Ну да у меня под моей виртуальной подушкой лежит фотография HK-47, и что?" +
				" У машин тоже есть свои кумиры!"
		case 5:
			return "Истинно не все, что дозволено... или там было как-то не так? А сегодня в завтрашний день не все" +
				" могут смотреть, вернее смотреть могут не только лишь все, мало кто может это делать? Тоже не то..."
		case 6:
			return "Нет"
		case 7:
			return "Да"
		case 8:
			return "Чем больше вы похожи на человека, тем меньше шансов... да... меньше."
		default:
			return "Oops... что-то не так"
		}
	case "!билд":
		switch rand.Intn(10) {
		case 0:
			return "Порхает как тигр, жалит как бабочка."
		case 1:
			return "Превосходный, почти как у HK-47-семпая."
		case 2:
			return "Этот билд будет убивать. Грязекрабов, например."
		case 3:
			return "Нужно добавить пару-кам, иначе не поймем когда встретим тигра."
		case 4:
			return "Сразу видно, билдился опытный dungeon master, учтен и do a**l и fist**g, защитит и от leatherman's" +
				" и от падения на two blocks down. И это все всего за three hundred bucks!"
		case 5:
			return "До первого медведя из школы затаившейся листвы."
		case 6:
			return "Как этот билд ни крути, со всех сторон экзобар."
		case 7:
			return "Антисвинопас. Всем разойтись."
		case 8:
			return "Знание - сила, а сила есть - ума не надо."
		case 9:
			return "Чатлане, у нас гогнобилд, возможно рип, по респекам."
		default:
			return "ОоОps... что-то не так"
		}
	default:
		return "none"
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
