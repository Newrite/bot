package twitch

import (
	"math/rand"
	"strings"
)

func (self *TwitchBot) resurrected(message, channel string) string {
	messageSlice := strings.Fields(message)
	if len(messageSlice) < 3 {
		return ""
	}
	if rand.Intn(99)+1 >= 70 {
		self.say("/untimeout "+messageSlice[2], channel)
		switch rand.Intn(4) {
		case 0:
			return messageSlice[2] + " Восстань, избранный мертвец!"
		case 1:
			return messageSlice[2] + " Встань, вечный воитель, твой час еще не пробил."
		case 2:
			return messageSlice[2] + " На этот раз черный камень душ спас тебя, некромант, но так будет не всегда."
		case 3:
			return messageSlice[2] + " Колесо времени повернулось на тебе, снова живи."
		default:
			return "Oops, что-то пошло не так"
		}
	} else {
		return ""
	}
}
