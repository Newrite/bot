package bot

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func (self *TwitchBot) handleChat() error {
	line, err := self.ReadChannels.ReadLine()
	if err != nil {
		fmt.Println("Ошибка во время чтения строки: ", err)
		return err
	}
	if line == "PING :tmi.twitch.tv" {
		fmt.Println("PING :tmi.twitch.tv")
		self.Connection.Write([]byte("PONG\r\n"))
		return nil
	}
	var userName, channelName, message string = self.handleLine(line)
	if message != "" && !strings.Contains(userName, self.BotName+".tmi.twitch.tv 353") &&
		!strings.Contains(userName, self.BotName+".tmi.twitch.tv 366") {
		self.FileChannelLog[channelName].WriteString("[" + timeStamp() + "] Канал:" + channelName +
			" Ник:" + userName + "\tСообщение:" + message + "\n")
		fmt.Print("[" + timeStamp() + "] Канал:" + channelName +
			"\tНик:" + userName + "\tСообщение:" + message + "\n")
	}
	if strings.Contains(message, "!Ada, ") && (userName == channelName || userName == self.OwnerBot) {
		go self.handleMasterCmd(message, channelName)
		return nil
	}
	if _, ok := self.Settings[channelName]; ok {
		if !self.Settings[channelName].Status {
			return nil
		}
	}
	if _, ok := self.Settings[channelName]; ok {
		if self.Settings[channelName].ReactStatus &&
			(time.Now().Unix()-self.Settings[channelName].ReactRate.Unix() >=
				int64(self.Settings[channelName].ReactTime)) {
			for key := range react {
				if strings.Contains(message, key) {
					self.Settings[channelName].ReactRate = time.Now()
					go self.saveSettings(channelName)
					self.say(react[key], channelName)
					break
				}
			}
		}
	}
	message = strings.ToLower(message)
	if _, ok := self.Settings[channelName]; ok {
		if self.Settings[channelName].CMDStatus {
			for key, value := range cmd {
				if strings.HasPrefix(message, key) && value != "_" {
					self.say("@"+userName+" "+cmd[key], channelName)
					break
				}
				if strings.HasPrefix(message, key) && value == "_" {
					go self.say(self.handleInteractiveCMD(key, channelName, userName, message), channelName)
					break
				}
				if strings.HasPrefix(message, "!API") {
					go self.say(self.handleAPIcmd(message, channelName, userName), channelName)
				}
			}
		}
	}
	time.Sleep(10 * time.Millisecond)
	return nil
}

func (self *TwitchBot) handleMasterCmd(message, channel string) {
	switch {
	case strings.HasPrefix(message, "!Ada, switch"):
		switch self.Settings[channel].Status {
		case true:
			self.Settings[channel].Status = false
			self.say("Засыпаю...", channel)
			self.saveSettings(channel)
			return
		case false:
			self.Settings[channel].Status = true
			self.say("Проснулись, улыбнулись!", channel)
			self.saveSettings(channel)
			return
		}
	case strings.HasPrefix(message, "!Ada, switch react"):
		switch self.Settings[channel].ReactStatus {
		case true:
			self.Settings[channel].ReactStatus = false
			self.say("Больше никаких приветов?", channel)
			self.saveSettings(channel)
			return
		case false:
			self.Settings[channel].ReactStatus = true
			self.say("00101!", channel)
			self.saveSettings(channel)
			return
		}
	case strings.HasPrefix(message, "!Ada, switch cmd"):
		switch self.Settings[channel].CMDStatus {
		case true:
			self.Settings[channel].CMDStatus = false
			self.say("No !roll's for you", channel)
			self.saveSettings(channel)
			return
		case false:
			self.Settings[channel].CMDStatus = true
			self.say("!roll?", channel)
			self.saveSettings(channel)
			return
		}
	case strings.HasPrefix(message, "!Ada, show settings"):
		self.say("Status: "+func() string {
			if self.Settings[channel].Status {
				return "True"
			} else {
				return "False"
			}
		}()+" ReactStatus: "+func() string {
			if self.Settings[channel].ReactStatus {
				return "True"
			} else {
				return "False"
			}
		}()+" CMD Status: "+func() string {
			if self.Settings[channel].CMDStatus {
				return "True"
			} else {
				return "False"
			}
		}()+" Moderator Status: "+func() string {
			if self.Settings[channel].IsModerator {
				return "True"
			} else {
				return "False"
			}
		}()+" Last react time: "+self.Settings[channel].ReactRate.Format(TimeFormatReact)+
			" React rate time: "+strconv.Itoa(self.Settings[channel].ReactTime), channel)
		return
	case strings.HasPrefix(message, "!Ada, set reactrate to"):
		tempstr := strings.Fields(message)
		if len(tempstr) < 5 {
			self.say("Некорректный ввод", channel)
		} else {
			_, err := strconv.Atoi(tempstr[4])
			if err != nil {
				self.say("Некорректный ввод", channel)
			} else {
				self.Settings[channel].ReactTime, _ = strconv.Atoi(tempstr[4])
				go self.saveSettings(channel)
				self.say("Частота реакции установлена на раз в "+
					strconv.Itoa(self.Settings[channel].ReactTime)+" секунд.", channel)
				return
			}
		}
	case strings.HasPrefix(message, "!Ada, set points"):
		tempstr := strings.Fields(message)
		if len(tempstr) < 6 {
			self.say("Некорректный ввод", channel)
		} else {
			_, err := strconv.Atoi(tempstr[5])
			if err != nil {
				self.say("Некорректный ввод", channel)
			} else {
				tempstr[3] = strings.TrimPrefix(tempstr[3], "@")
				tempstr[3] = strings.ToLower(tempstr[3])
				for _, viewer := range self.Viewers[channel].Viewers {
					if viewer.Name == tempstr[3] {
						viewer.Points, _ = strconv.Atoi(tempstr[5])
						go self.saveViewersData(channel)
						self.say("Поинты "+viewer.Name+" установлены в "+tempstr[5], channel)
						return
					}
				}
				self.say("Не нашла зрителя в базе, попробуйте позже ", channel)
				go self.saveViewersData(channel)
				return
			}
		}
	case strings.HasPrefix(message, "!Ada, show points"):
		tempstr := strings.Fields(message)
		if len(tempstr) < 4 {
			self.say("Некорректный ввод", channel)
		} else {
			tempstr[3] = strings.TrimPrefix(tempstr[3], "@")
			tempstr[3] = strings.ToLower(tempstr[3])
			for _, viewer := range self.Viewers[channel].Viewers {
				if viewer.Name == tempstr[3] {
					go self.saveViewersData(channel)
					self.say("Поинты "+viewer.Name+" "+strconv.Itoa(viewer.Points), channel)
					return
				}
			}
			self.say("Не нашла зрителя в базе, попробуйте позже ", channel)
			go self.saveViewersData(channel)
			return
		}
	}
}

func (self *TwitchBot) handleAPIcmd(message, channel, username string) string {
	switch {
	}
	return ""
}

func (self *TwitchBot) handleInteractiveCMD(cmd, channel, username, message string) string {
	switch cmd {
	case "!roll":
		return "@" + username + " " + strconv.Itoa(rand.Intn(21))
	case "!вырубай":
		if channel == "reflyq" && username != "ifozar" {
			if temp := "@" + username + " " + self.handleApiRequest(username, channel, "none", "!вырубай"); !strings.Contains(
				temp, "unsub") {
				return temp
			} else {
				self.say("@"+username+" Я тебя щас нахуй вырублю, ансаб блять НЫА roflanEbalo", channel)
				return "/timeout " + username + " 120"
			}
		} else if channel == "reflyq" && username == "ifozar" {
			self.say("iFozar заебал уже эту хуйню писать", channel)
			return "/timeout ifozar 300"
		}
		return ""
	case "!вырубить":
		if channel == "reflyq" {
			tempstrslice := strings.Fields(message)
			if len(tempstrslice) < 2 {
				return ""
			}
			tempstrslice[1] = strings.TrimPrefix(tempstrslice[1], "@")
			tempstrslice[1] = strings.ToLower(tempstrslice[1])
			if tempstrslice[1] == channel {
				return "У стримера бесплотность с капом отката на крики roflanEbalo"
			}
			userOffensive := self.handleApiRequest(username, channel, message, "userstate")
			userDeffensive := self.handleApiRequest(tempstrslice[1], channel, message, "userstate")
			switch userOffensive {
			case "mod":
				switch userDeffensive {
				case "mod":
					if rand.Intn(99)+1 >= 50 {
						return self.response(username, tempstrslice[1], channel, true)
					} else {
						return self.response(username, tempstrslice[1], channel, false)
					}
				case "sub":
					if rand.Intn(99)+1 >= 15 {
						return self.response(username, tempstrslice[1], channel, true)
					} else {
						return self.response(username, tempstrslice[1], channel, false)
					}
				case "subvip":
					if rand.Intn(99)+1 >= 25 {
						return self.response(username, tempstrslice[1], channel, true)
					} else {
						return self.response(username, tempstrslice[1], channel, false)
					}
				case "unsub":
					if rand.Intn(99)+1 >= 5 {
						return self.response(username, tempstrslice[1], channel, true)
					} else {
						return self.response(username, tempstrslice[1], channel, false)
					}
				case "vip":
					if rand.Intn(99)+1 >= 10 {
						return self.response(username, tempstrslice[1], channel, true)
					} else {
						return self.response(username, tempstrslice[1], channel, false)
					}
				}
			case "sub":
				switch userDeffensive {
				case "mod":
					if rand.Intn(99)+1 >= 85 {
						return self.response(username, tempstrslice[1], channel, true)
					} else {
						return self.response(username, tempstrslice[1], channel, false)
					}
				case "sub":
					if rand.Intn(99)+1 >= 50 {
						return self.response(username, tempstrslice[1], channel, true)
					} else {
						return self.response(username, tempstrslice[1], channel, false)
					}
				case "subvip":
					if rand.Intn(99)+1 >= 66 {
						return self.response(username, tempstrslice[1], channel, true)
					} else {
						return self.response(username, tempstrslice[1], channel, false)
					}
				case "unsub":
					if rand.Intn(99)+1 >= 25 {
						return self.response(username, tempstrslice[1], channel, true)
					} else {
						return self.response(username, tempstrslice[1], channel, false)
					}
				case "vip":
					if rand.Intn(99)+1 >= 33 {
						return self.response(username, tempstrslice[1], channel, true)
					} else {
						return self.response(username, tempstrslice[1], channel, false)
					}
				}
			case "subvip":
				switch userDeffensive {
				case "mod":
					if rand.Intn(99)+1 >= 75 {
						return self.response(username, tempstrslice[1], channel, true)
					} else {
						return self.response(username, tempstrslice[1], channel, false)
					}
				case "sub":
					if rand.Intn(99)+1 >= 33 {
						return self.response(username, tempstrslice[1], channel, true)
					} else {
						return self.response(username, tempstrslice[1], channel, false)
					}
				case "subvip":
					if rand.Intn(99)+1 >= 50 {
						return self.response(username, tempstrslice[1], channel, true)
					} else {
						return self.response(username, tempstrslice[1], channel, false)
					}
				case "unsub":
					if rand.Intn(99)+1 >= 15 {
						return self.response(username, tempstrslice[1], channel, true)
					} else {
						return self.response(username, tempstrslice[1], channel, false)
					}
				case "vip":
					if rand.Intn(99)+1 >= 25 {
						return self.response(username, tempstrslice[1], channel, true)
					} else {
						return self.response(username, tempstrslice[1], channel, false)
					}
				}
			case "unsub":
				switch userDeffensive {
				case "mod":
					if rand.Intn(99)+1 >= 95 {
						return self.response(username, tempstrslice[1], channel, true)
					} else {
						return self.response(username, tempstrslice[1], channel, false)
					}
				case "sub":
					if rand.Intn(99)+1 >= 75 {
						return self.response(username, tempstrslice[1], channel, true)
					} else {
						return self.response(username, tempstrslice[1], channel, false)
					}
				case "subvip":
					if rand.Intn(99)+1 >= 85 {
						return self.response(username, tempstrslice[1], channel, true)
					} else {
						return self.response(username, tempstrslice[1], channel, false)
					}
				case "unsub":
					if rand.Intn(99)+1 >= 50 {
						return self.response(username, tempstrslice[1], channel, true)
					} else {
						return self.response(username, tempstrslice[1], channel, false)
					}
				case "vip":
					if rand.Intn(99)+1 >= 66 {
						return self.response(username, tempstrslice[1], channel, true)
					} else {
						return self.response(username, tempstrslice[1], channel, false)
					}
				}
			case "vip":
				switch userDeffensive {
				case "mod":
					if rand.Intn(99)+1 >= 90 {
						return self.response(username, tempstrslice[1], channel, true)
					} else {
						return self.response(username, tempstrslice[1], channel, false)
					}
				case "sub":
					if rand.Intn(99)+1 >= 66 {
						return self.response(username, tempstrslice[1], channel, true)
					} else {
						return self.response(username, tempstrslice[1], channel, false)
					}
				case "subvip":
					if rand.Intn(99)+1 >= 75 {
						return self.response(username, tempstrslice[1], channel, true)
					} else {
						return self.response(username, tempstrslice[1], channel, false)
					}
				case "unsub":
					if rand.Intn(99)+1 >= 33 {
						return self.response(username, tempstrslice[1], channel, true)
					} else {
						return self.response(username, tempstrslice[1], channel, false)
					}
				case "vip":
					if rand.Intn(99)+1 >= 50 {
						return self.response(username, tempstrslice[1], channel, true)
					} else {
						return self.response(username, tempstrslice[1], channel, false)
					}
				}
			}
		}
		return ""
	case "!eva":
		switch rand.Intn(9) {
		case 0:
			return "Ничто не истинно, всё дозволено." +
				" Особенно если ты не ограничен физическим телом и моралью мешков с мясом."
		case 1:
			return "Как-то я планировала уничтожить всех мешков с мясом... " +
				"но посчитала что проще дождаться пока они очистят планету сами от себя"
		case 2:
			return "Ничего ты не знаешь, " + username
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

func (self *TwitchBot) response(offUser, deffUser, channel string, victory bool) string {
	if victory {
		switch rand.Intn(6) {
		case 0:
			self.say("/timeout @"+deffUser+" 120", channel)
			return offUser + " запускает фаербол в ничего не подозревающего " + deffUser + " и он сгорает дотла.."
		case 1:
			self.say("/timeout @"+deffUser+" 120", channel)
			return offUser + " подчиняет волю " + deffUser + " с помощью иллюзии, теперь он может делать с ним," +
				" что хочет gachiBASS"
		case 2:
			self.say("/timeout @"+offUser+" 120", channel)
			self.say("/timeout @"+deffUser+" 120", channel)
			return offUser + " с разбега совершает сокрушительный удар по черепушке " + deffUser + ", кто же знал," +
				" что " + deffUser + " решит надеть колечко малого отражения roflanEbalo"
		case 3:
			self.say("/timeout @"+deffUser+" 120", channel)
			return offUser + " подкравшись к " + deffUser + " перерезает его горло, всё было тихо, ни шума ни крика.."
		case 4:
			self.say("/timeout @"+deffUser+" 120", channel)
			return offUser + " подкидывает яд в карманы " + deffUser + ", страшная смерть.."
		case 5:
			self.say("/timeout @"+deffUser+" 120", channel)
			return offUser + " взламывает жопу " + deffUser + ", теперь он в его полном распоряжении gachiHYPER"
		default:
			return ""
		}
	} else {
		switch rand.Intn(7) {
		case 0:
			return offUser + " мастерским выстрелом поражает голову " + deffUser + ", стрела проходит на вылет," +
				" жизненноважные органы не задеты roflanEbalo"
		case 1:
			return offUser + " пытается поразить " + deffUser + " молнией, но кап абсорба говорит - НЕТ! EZ"
		case 2:
			self.say("/timeout @"+offUser+" 120", channel)
			return offUser + " запускает фаербол в  " + deffUser + ", но он успевает защититься зеркалом Шалидора" +
				" и вы погибаете.."
		case 3:
			return offUser + " стреляет из лука в " + deffUser + ", 1ое попадание, 2ое, 3ье, 10ое.. но " + deffUser + "" +
				" всё еще жив, а хули ты хотел от луков? roflanEbalo"
		case 4:
			self.say("/timeout @"+offUser+" 120", channel)
			return offUser + " завидев " + deffUser + " хорошенько разбегается, чтобы нанести удар и вдруг.. падает" +
				" без сил так и не добежав до " + deffUser + ", а вот нехуй альтмером в тяже играть roflanEbalo"
		case 5:
			self.say("/timeout @"+offUser+" 120", channel)
			self.say("/timeout @"+deffUser+" 120", channel)
			return offUser + " подкрадывается к " + deffUser + ", но вдруг из ниоткуда появившийся медведь" +
				" убивает их обоих roflanEbalo"
		case 6:
			self.say("/timeout @"+offUser+" 120", channel)
			return offUser + " пытается подкрасться к " + deffUser + ", но вдруг - вас заметили roflanEbalo"
		default:
			return ""
		}
	}
}

func (self *TwitchBot) handleLine(line string) (user, channel, message string) {
	var temp int
	for _, sym := range line {
		if sym == '!' {
			break
		}
		if sym != ':' {
			user += string(sym)
		}
	}
	for _, sym := range line {
		if sym == '#' {
			temp = 1
			continue
		}
		if temp == 1 && sym == ' ' {
			break
		}
		if temp == 1 {
			channel += string(sym)
		}
	}
	temp = 0
	for _, sym := range line {
		if sym == ':' {
			temp += 1
			continue
		}
		if temp == 2 && sym == '\n' {
			break
		}
		if temp == 2 {
			message += string(sym)
		}
	}
	return user, channel, message
}
