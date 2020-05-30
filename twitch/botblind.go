package twitch

import (
	"time"
)

var channelBlinde string = "blindwalkerboy"

func (self *TwitchBot) sendBlindRepeatMessage() {
	for {
		if self.handleApiRequest("", channelBlinde, "", "streamStatus") == "online" {
			self.say("Привет, дружище! Приглашаю тебя в лучшее сообщество по Requiem"+
				" и ламповое убежище для настоящих мужчин - https://discord.gg/4yqdafW", channelBlinde)
		}
		time.Sleep(15 * time.Minute)
	}
}

func (self *TwitchBot) handleBlindCMD(userName, message, cmd string) string {
	switch cmd {
	case "!чезаигра":
		return "Skyrim с модификацией Requiem от Xandr'а. Ебашим на максимальной сложности без смертей, чтобы жизнь мёдом не казалась PepeSmoke"
	case "!скиллуха":
		return "Попробуй найди peepoClown"
	case "!вызватьсанитаров":
		return "Corpsman Corpsman Corpsman Corpsman Corpsman"
	case "!труба":
		return "Наш канал на Трубе, где мы занимаемся максимально грязными делишками - https://www.youtube.com/blindwalker"
	case "!вырубай":
		return "А хуй тебе pepoGun"
	case "!билд":
		return "Суровая Арбалетчица в Тяжёлых доспехах Bratishka"
	case "!тыктоблять":
		return "30-летний дед Гришаня из Краснодара, находящийся на грани нервного срыва"
	case "!дискорд":
		return "Секретное логово для опытных мужчин, а также сообщество по Requiem для Работяг - https://discord.gg/4yqdafW"
	default:
		return "none"
	}
}
