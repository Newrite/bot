package bots

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const inMuteTime int = 30
const useMuteTime int = 180

var channelRflyq string = "reflyq"

func (bt *BotTwitch) handleReflyqCMD(userName, message, cmd string) string {
	switch cmd {
	case "вырубай":
		if userName == "ifozar" {
			bt.say("iFozar заебал уже эту хуйню писать", channelRflyq)
			return "/timeout ifozar 300"
		}
		if temp := "@" + userName + " " + bt.handleApiRequest(userName, channelRflyq, "none", "!вырубай"); !strings.Contains(temp, "unsub") {
			return temp
		} else {
			bt.say("@"+userName+" Я тебя щас нахуй вырублю, ансаб блять НЫА roflanEbalo", channelRflyq)
			return "/timeout " + userName + " 120"
		}
	case "вырубить":
		tempStrSlice := strings.Fields(message)
		if len(tempStrSlice) < 2 {
			return ""
		}
		tempStrSlice[1] = strings.TrimPrefix(tempStrSlice[1], "@")
		tempStrSlice[1] = strings.ToLower(tempStrSlice[1])
		userOffensive := bt.handleApiRequest(userName, channelRflyq, message, "userstate")
		userDeffensive := bt.handleApiRequest(tempStrSlice[1], channelRflyq, message, "userstate")
		switch bt.handleExeption(userName, tempStrSlice[1], userOffensive, userDeffensive) {
		case "killer":
			bt.say("/timeout @"+userName+" 10", channelRflyq)
			return "Камень бьет ножницы, а я бью твое ебало спамер, НЫА roflanEbalo"
		case "streamerDeff":
			return "У стримера бесплотность с капом отката на крики roflanEbalo"
		case "killed":
			return tempStrSlice[1] + " уже вырублен"
		case "modOff":
			return userName + ", ты что забыл свой банхаммер дома? monkaHmm"
		case "modDeff":
			return userName + ", Agakakskagesh Agakakskagesh Agakakskagesh"
		case "reflyqkiller":
			go bt.userInMuteFunc(tempStrSlice[1], 120)
			bt.say("/timeout @"+tempStrSlice[1]+" "+strconv.Itoa(inMuteTime), channelRflyq)
			return "Reflyq произносит YOL TooR Shul и испепеляет " + tempStrSlice[1] + " monkaX"
		case "shiza":
			return "Осуждаю roflanEbalo"
		}
		switch userOffensive {
		case "sub":
			switch userDeffensive {
			case "sub":
				if rand.Intn(99)+1 >= 50 {
					return bt.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return bt.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			case "subvip":
				if rand.Intn(99)+1 >= 75 {
					return bt.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return bt.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			case "unsub":
				if rand.Intn(99)+1 >= 15 {
					return bt.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return bt.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			case "vip":
				if rand.Intn(99)+1 >= 25 {
					return bt.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return bt.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			}
		case "subvip":
			switch userDeffensive {
			case "sub":
				if rand.Intn(99)+1 >= 25 {
					return bt.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return bt.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			case "subvip":
				if rand.Intn(99)+1 >= 50 {
					return bt.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return bt.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			case "unsub":
				if rand.Intn(99)+1 >= 5 {
					return bt.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return bt.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			case "vip":
				if rand.Intn(99)+1 >= 15 {
					return bt.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return bt.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			}
		case "unsub":
			switch userDeffensive {
			case "sub":
				if rand.Intn(99)+1 >= 85 {
					return bt.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return bt.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			case "subvip":
				if rand.Intn(99)+1 >= 95 {
					return bt.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return bt.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			case "unsub":
				if rand.Intn(99)+1 >= 50 {
					return bt.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return bt.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			case "vip":
				if rand.Intn(99)+1 >= 75 {
					return bt.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return bt.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			}
		case "vip":
			switch userDeffensive {
			case "sub":
				if rand.Intn(99)+1 >= 75 {
					return bt.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return bt.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			case "subvip":
				if rand.Intn(99)+1 >= 85 {
					return bt.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return bt.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			case "unsub":
				if rand.Intn(99)+1 >= 25 {
					return bt.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return bt.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			case "vip":
				if rand.Intn(99)+1 >= 50 {
					return bt.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, true)
				} else {
					return bt.reflyqAnswer(userName, tempStrSlice[1], channelRflyq, false)
				}
			}
		}
	default:
		return "none"
	}
	return "Ашибка (handleReflyqCMD)"
}

func (bt *BotTwitch) reflyqAnswer(offUser, deffUser, channel string, victory bool) string {
	go bt.userUseMuteFunc(offUser, time.Duration(useMuteTime))
	if victory {
		switch rand.Intn(6) {
		case 0:
			go bt.userInMuteFunc(deffUser, time.Duration(inMuteTime))
			bt.say("/timeout @"+deffUser+" "+strconv.Itoa(inMuteTime), channel)
			return offUser + " запускает фаербол в ничего не подозревающего " + deffUser + " и он сгорает дотла.."
		case 1:
			go bt.userInMuteFunc(deffUser, time.Duration(inMuteTime))
			bt.say("/timeout @"+deffUser+" "+strconv.Itoa(inMuteTime), channel)
			return offUser + " подчиняет волю " + deffUser + " с помощью иллюзии, теперь он может делать с ним," +
				" что хочет gachiBASS"
		case 2:
			go bt.userInMuteFunc(offUser, time.Duration(inMuteTime))
			go bt.userInMuteFunc(deffUser, time.Duration(inMuteTime))
			bt.say("/timeout @"+offUser+" "+strconv.Itoa(inMuteTime), channel)
			bt.say("/timeout @"+deffUser+" "+strconv.Itoa(inMuteTime), channel)
			return offUser + " с разбега совершает сокрушительный удар по черепушке " + deffUser + ", кто же знал," +
				" что " + deffUser + " решит надеть колечко малого отражения roflanEbalo"
		case 3:
			go bt.userInMuteFunc(deffUser, time.Duration(inMuteTime))
			bt.say("/timeout @"+deffUser+" "+strconv.Itoa(inMuteTime), channel)
			return offUser + " подкравшись к " + deffUser + " перерезает его горло, всё было тихо, ни шума ни крика.."
		case 4:
			go bt.userInMuteFunc(deffUser, time.Duration(inMuteTime))
			bt.say("/timeout @"+deffUser+" "+strconv.Itoa(inMuteTime), channel)
			return offUser + " подкидывает яд в карманы " + deffUser + ", страшная смерть.."
		case 5:
			go bt.userInMuteFunc(deffUser, time.Duration(inMuteTime))
			bt.say("/timeout @"+deffUser+" "+strconv.Itoa(inMuteTime), channel)
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
			go bt.userInMuteFunc(offUser, time.Duration(inMuteTime))
			bt.say("/timeout @"+offUser+" "+strconv.Itoa(inMuteTime), channel)
			return offUser + " запускает фаербол в  " + deffUser + ", но он успевает защититься зеркалом Шалидора" +
				" и вы погибаете.."
		case 3:
			return offUser + " стреляет из лука в " + deffUser + ", 1ое попадание, 2ое, 3ье, 10ое.. но " + deffUser + "" +
				" всё еще жив, а хули ты хотел от луков? roflanEbalo"
		case 4:
			go bt.userInMuteFunc(offUser, time.Duration(inMuteTime))
			bt.say("/timeout @"+offUser+" "+strconv.Itoa(inMuteTime), channel)
			return offUser + " завидев " + deffUser + " хорошенько разбегается, чтобы нанести удар и вдруг.. падает" +
				" без сил так и не добежав до " + deffUser + ", а вот нехуй альтмером в тяже играть roflanEbalo"
		case 5:
			go bt.userInMuteFunc(offUser, time.Duration(inMuteTime))
			go bt.userInMuteFunc(deffUser, time.Duration(inMuteTime))
			bt.say("/timeout @"+offUser+" "+strconv.Itoa(inMuteTime), channel)
			bt.say("/timeout @"+deffUser+" "+strconv.Itoa(inMuteTime), channel)
			return offUser + " подкрадывается к " + deffUser + ", но вдруг из ниоткуда появившийся медведь" +
				" убивает их обоих roflanEbalo"
		case 6:
			go bt.userInMuteFunc(offUser, time.Duration(inMuteTime))
			bt.say("/timeout @"+offUser+" "+strconv.Itoa(inMuteTime), channel)
			return offUser + " пытается подкрасться к " + deffUser + ", но вдруг - вас заметили roflanEbalo"
		default:
			return ""
		}
	}
}

func (bt *BotTwitch) userInMuteFunc(user string, duration time.Duration) {
	bt.MutedUsers[user] = true
	time.Sleep(duration * time.Second)
	bt.MutedUsers[user] = false
}

func (bt *BotTwitch) userUseMuteFunc(user string, duration time.Duration) {
	bt.UsersMuted[user] = true
	time.Sleep(duration * time.Second)
	bt.UsersMuted[user] = false
}

func (bt *BotTwitch) handleExeption(userOff, userDeff, userOffStatus, userDeffStatus string) string {
	killer := false
	killed := false
	if _, ok := bt.UsersMuted[userOff]; ok {
		if bt.UsersMuted[userOff] {
			killer = true
		}
	}
	if _, ok := bt.MutedUsers[userDeff]; ok {
		if bt.MutedUsers[userOff] {
			killed = true
		}
	}
	switch {
	case killer:
		return "killer"
	case userOff == userDeff:
		return "shiza"
	case userDeff == channelRflyq:
		return "streamerDeff"
	case killed:
		return "killed"
	case userOff == channelRflyq:
		return "reflyqkiller"
	case userOffStatus == "mod":
		return "modOff"
	case userDeffStatus == "mod":
		return "modDeff"
	default:
		return ""
	}
}
