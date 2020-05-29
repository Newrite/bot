package twitch

import (
	"math/rand"
	"strings"
)

var channelRflyq string = "reflyq"

func (self *TwitchBot) handleReflyqCMD(userName, message, cmd string) string {
	switch cmd {
	case "!вырубай":
		if userName == "ifozar" {
			self.say("iFozar заебал уже эту хуйню писать", channelRflyq)
			return "/timeout ifozar 300"
		}
		if temp := "@" + userName + " " + self.handleApiRequest(userName, channelRflyq, "none", "!вырубай"); !strings.Contains(temp, "unsub") {
			return temp
		} else {
			self.say("@"+userName+" Я тебя щас нахуй вырублю, ансаб блять НЫА roflanEbalo", channelRflyq)
			return "/timeout " + userName + " 120"
		}
	case "!вырубить":
		tempStrSlice := strings.Fields(message)
		if len(tempStrSlice) < 2 {
			return ""
		}
		tempStrSlice[1] = strings.TrimPrefix(tempStrSlice[1], "@")
		tempStrSlice[1] = strings.ToLower(tempStrSlice[1])
		if tempStrSlice[1] == channelRflyq {
			return "У стримера бесплотность с капом отката на крики roflanEbalo"
		}
		userOffensive := self.handleApiRequest(userName, channelRflyq, message, "userstate")
		userDeffensive := self.handleApiRequest(tempStrSlice[1], channelRflyq, message, "userstate")
		switch userOffensive {
		case "mod":
			switch userDeffensive {
			case "mod":
				if rand.Intn(99)+1 >= 50 {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			case "sub":
				if rand.Intn(99)+1 >= 15 {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			case "subvip":
				if rand.Intn(99)+1 >= 25 {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			case "unsub":
				if rand.Intn(99)+1 >= 5 {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			case "vip":
				if rand.Intn(99)+1 >= 10 {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			}
		case "sub":
			switch userDeffensive {
			case "mod":
				if rand.Intn(99)+1 >= 85 {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			case "sub":
				if rand.Intn(99)+1 >= 50 {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			case "subvip":
				if rand.Intn(99)+1 >= 66 {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			case "unsub":
				if rand.Intn(99)+1 >= 25 {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			case "vip":
				if rand.Intn(99)+1 >= 33 {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			}
		case "subvip":
			switch userDeffensive {
			case "mod":
				if rand.Intn(99)+1 >= 75 {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			case "sub":
				if rand.Intn(99)+1 >= 33 {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			case "subvip":
				if rand.Intn(99)+1 >= 50 {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			case "unsub":
				if rand.Intn(99)+1 >= 15 {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			case "vip":
				if rand.Intn(99)+1 >= 25 {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			}
		case "unsub":
			switch userDeffensive {
			case "mod":
				if rand.Intn(99)+1 >= 95 {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			case "sub":
				if rand.Intn(99)+1 >= 75 {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			case "subvip":
				if rand.Intn(99)+1 >= 85 {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			case "unsub":
				if rand.Intn(99)+1 >= 50 {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			case "vip":
				if rand.Intn(99)+1 >= 66 {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			}
		case "vip":
			switch userDeffensive {
			case "mod":
				if rand.Intn(99)+1 >= 90 {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			case "sub":
				if rand.Intn(99)+1 >= 66 {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			case "subvip":
				if rand.Intn(99)+1 >= 75 {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			case "unsub":
				if rand.Intn(99)+1 >= 33 {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			case "vip":
				if rand.Intn(99)+1 >= 50 {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return self.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			}
		}
	default:
		return "none"
	}
	return "none"
}

func (self *TwitchBot) reflyqAnswer(offUser, deffUser, channel string, victory bool) string {
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
