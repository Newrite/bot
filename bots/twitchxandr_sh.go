package bots

import (
	"bot/resource"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func (bt *BotTwitch) resurrected(message, channel string) string {
	messageSlice := strings.Fields(message)
	if len(messageSlice) < 3 {
		return ""
	}
	if rand.Intn(99)+1 >= 70 {
		bt.say("/untimeout "+messageSlice[2], channel)
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

func xandrSendRepeatMessage() {
	for {
		time.Sleep(15 * time.Minute)
		if SingleTwitch().handleApiRequest("", "xandr_sh", "", "streamStatus") == "online" {
			SingleTwitch().say(`команды бота: !help / VK: https://vk.com/xandr_tv / YouTube: https://www.youtube.com/channel/UC0oObsGZKntyAP_OoMnFIPA / GoodGame: https://goodgame.ru/channel/Xandr_Sh/`, "xandr_sh")
			SingleGoodGame().say(`команды бота: !help / VK: https://vk.com/xandr_tv / YouTube: https://www.youtube.com/channel/UC0oObsGZKntyAP_OoMnFIPA / GoodGame: https://goodgame.ru/channel/Xandr_Sh/`, SingleGoodGame().Channels[0])
		}
	}
}

func xandrHandleCMD(userName, message, cmd, platform, ggUserID string) string {
	switch cmd {
	case "love":
		msgSlice := strings.Fields(message)
		if len(msgSlice) < 2 {
			return `Соберитесь с мыслями и решите уже наконец, с кем вы готовы проверить свою любовь!`
		}
		switch platform {
		case TW:
			return strconv.Itoa(rand.Intn(99)+1) + `% <3 между ` + userName + ` и` + strings.TrimPrefix(message, msgSlice[0])
		case GG:
			return strconv.Itoa(rand.Intn(99)+1) + `% :love: между ` + userName + ` и` + strings.TrimPrefix(message, msgSlice[0])
		}
	case "seppuku":
		switch platform {
		case TW:
			return `/timeout @` + userName + ` 60 харакири`
		case GG:
			_, err := SingleGoodGame().Connection.Write([]byte(`{
    "type": "ban",
    "data": {
        "channel_id": "` + SingleGoodGame().Channels[0] + `",
        "ban_channel": "` + SingleGoodGame().Channels[0] + `",
        "user_id": "` + ggUserID + `",
        "duration": 60,
        "reason": "харакири",
        "comment": "!харакири",
        "show_ban": true
    }
}`))
			if err != nil {
				log.WithFields(log.Fields{
					"package":  "bots",
					"function": "Connection.Write",
					"file":     "twitchxandr_sh.go",
					"body":     "xandrHandleCMD",
					"error":    err,
				}).Errorln("Ошибка во время отправки бана.")
			}
			return ""
		}
	case "roulette":
		switch {
		case rand.Intn(99)+1 < 50:
			switch platform {
			case TW:
				SingleTwitch().say(`/me подносит револьвер к виску `+userName, "xandr_sh")
				time.Sleep(2 * time.Second)
				SingleTwitch().say(`/timeout @`+userName+` 120 рулетка`, "xandr_sh")
				return `револьвер выстреливает! ` + userName + ` погибает у чатлан на руках BibleThump 7`
			case GG:
				SingleGoodGame().say(`подносит револьвер к виску `+userName, SingleGoodGame().Channels[0])
				time.Sleep(2 * time.Second)
				_, err := SingleGoodGame().Connection.Write([]byte(`{
    "type": "ban",
    "data": {
        "channel_id": "` + SingleGoodGame().Channels[0] + `",
        "ban_channel": "` + SingleGoodGame().Channels[0] + `",
        "user_id": "` + ggUserID + `",
        "duration": 120,
        "reason": "рулетка",
        "comment": "!рулетка",
        "show_ban": true
    }
}`))
				if err != nil {
					log.WithFields(log.Fields{
						"package":  "bots",
						"function": "Connection.Write",
						"file":     "twitchxandr_sh.go",
						"body":     "xandrHandleCMD",
						"error":    err,
					}).Errorln("Ошибка во время отправки бана.")
				}
				return `револьвер выстреливает! ` + userName + ` погибает у чатлан на руках :skull:`
			}
		case rand.Intn(99)+1 >= 50:
			switch platform {
			case TW:
				SingleTwitch().say(`/me подносит револьвер к виску `+userName, "xandr_sh")
				time.Sleep(2 * time.Second)
				return `револьвер издает щелчок. ` + userName + ` выживает! PogChamp`
			case GG:
				SingleGoodGame().say(`подносит револьвер к виску `+userName, SingleGoodGame().Channels[0])
				time.Sleep(2 * time.Second)
				return `револьвер издает щелчок. ` + userName + ` выживает! :kerrisad:`
			}
		}
	case "8ball":
		return userName + `, ` + resource.Eva8ball[rand.Intn(resource.Count8Ball)]
	default:
		return "Ашибка xandrHandleCMD"
	}
	return "Ашибка xandrHandleCMD2"
}
