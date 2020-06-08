package bots

import (
	"bot/resource"
	"math/rand"
	"strings"
	"sync"
	"time"
)

const TimeFormat = "2006.01.02 15:04"

func timeStamp() string {
	return time.Now().Format(TimeFormat)
}

const VERSION = `1.0.1`
const cmdCOUNT = 24
const TW = "TW"
const GG = "GG"
const DIS = "DIS"
const TwPrefix = "!"
const GgPrefix = "!"
const DisPrefix = "~"

var twitch = &BotTwitch{}
var goodgame = &BotGoodGame{}
var discord = &BotDiscord{}
var once sync.Once

func SingleTwitch() *BotTwitch {
	once.Do(func() {
		twitch = &BotTwitch{}
	})
	return twitch
}

func SingleGoodGame() *BotGoodGame {
	once.Do(func() {
		goodgame = &BotGoodGame{}
	})
	return goodgame
}

func SingleDiscord() *BotDiscord {
	once.Do(func() {
		discord = &BotDiscord{}
	})
	return discord
}

func checkCMD(userName, channel, cmd, platform, message, originMessage string) string {
	var pr string
	switch platform {
	case TW:
		pr = TwPrefix
	case GG:
		pr = GgPrefix
	case DIS:
		pr = DisPrefix
	}
	for _, cL := range CMDList {
		for _, pl := range cL.Platform {
			if pl == "all" || pl == platform {
				for _, ch := range cL.Channels {
					if ch == "all" || ch == channel {
						for _, us := range cL.Users {
							if us == "all" || us == userName {
								for _, cc := range cL.Command {
									if pr+cc == cmd {
										return handleCMD(userName, channel, cL.Request, platform, message, originMessage)
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return ""
}

func handleCMD(userName, channel, cmd, platform, message, originMessage string) string {
	switch cmd {
	case "add quote":
		msgSlice := strings.Fields(originMessage)
		if len(msgSlice) < 2 {
			return "Некорректный ввод"
		} else {
			quote := strings.TrimPrefix(originMessage, msgSlice[0])
			resource.AddQuoteDB(quote)
			return userName+", схоронила в бд: "+quote
		}
	case "get quote":
		return resource.DBQuote()
	case "live":
		switch platform {
		case TW:
			d := time.Since(time.Unix(SingleTwitch().uptime, 0))
			return userName + ", " + d.String()
		case GG:
			d := time.Since(time.Unix(SingleGoodGame().uptime, 0))
			return userName + ", " + d.String()
		case DIS:
			d := time.Since(time.Unix(SingleDiscord().uptime, 0))
			return userName + ", " + d.String()
		}
	case "VERSION":
		return userName + ", " + VERSION
	case "build":
		return userName + ", " + resource.BuildAnswers[rand.Intn(resource.CountBuilds)]
	case "eva":
		return userName + ", " + resource.EvaAnswers[rand.Intn(resource.CountAnswers)]
	case "roll":
		return userName + ", " + resource.Rolls(message)
	case "bot":
		return userName + ", AdaIsEva, написана на GoLang v1.14 без использования сторонних библиотек (для GG и twitch). " +
			"Для дискорда использовалось discordgo by bwmarrin." +
			"Живет на VPS с убунтой размещенном в москоу сити. Рекомендации, пожелания и" +
			" прочая можно присылать на adaiseva.newrite@gmail.com"
	case "help":
		switch platform {
		case TW:
			return userName + ", Доступные комманды: build, eva, roll, bot, uptime, live, help, master help." +
				"Взаимодействие с БД: Добавить квоту - addquote <message> или aq <message>, хранить можно даже ссылки." +
				"Получить рандомну квоту из Бд - q или quote. Используйте префикс - " + TwPrefix
		case GG:
			return userName + ", Доступные комманды: build, eva, roll, bot, uptime (берет с твича), live, help." +
				"Взаимодействие с БД: Добавить квоту - addquote <message> или aq <message>, хранить можно даже ссылки." +
				"Получить рандомну квоту из Бд - q или quote. Используйте префикс - " + GgPrefix
		case DIS:
			return userName + ", Доступные комманды: build, eva, roll, bot, live, help." +
				"Взаимодействие с БД: Добавить квоту - addquote <message> или aq <message>, хранить можно даже ссылки." +
				"Получить рандомну квоту из Бд - q или quote. Используйте префикс - " + DisPrefix
		}
	case "master help":
		switch platform {
		case TW:
			return userName + " Владелец бота либо канала может переключить активность бота коммандой !Ada, switch." +
				" Реакции на всякое разное командой !Ada, switch react. " +
				"Переключить отзыв на различные команды !Ada, switch cmd." +
				" !Ada, set reactrate to <значение> выставляет настройку частоты реакции на различные сообщения в чате"
		}
	case "uptime":
		switch platform {
		case TW:
			return userName + ", " + SingleTwitch().handleApiRequest(userName, channel, message, "uptime")
		case GG:
			return userName + ", " + SingleTwitch().handleRequests("uptime")
		}
	case "вырубить":
		return SingleTwitch().handleReflyqCMD(userName, message, cmd)
	case "вырубайReflyq":
		return SingleTwitch().handleReflyqCMD(userName, message, "вырубай")
	case "helpReflyq":
		return userName + ", Доступные комманды: build, eva, roll, bot, uptime, live, help, master help." +
			" Уникальные на канале: вырубить, вырубай." +
			"Взаимодействие с БД: Добавить квоту - addquote <message> или aq <message>, хранить можно даже ссылки." +
			"Получить рандомну квоту из Бд - q или quote. Используйте префикс - " + TwPrefix
	case "вырубайBlind":
		return SingleTwitch().handleBlindCMD(userName, message, "вырубай")
	case "чезаигра":
		return SingleTwitch().handleBlindCMD(userName, message, cmd)
	case "скиллуха":
		return SingleTwitch().handleBlindCMD(userName, message, cmd)
	case "вызватьсанитаров":
		return SingleTwitch().handleBlindCMD(userName, message, cmd)
	case "труба":
		return SingleTwitch().handleBlindCMD(userName, message, cmd)
	case "вырубай":
		return SingleTwitch().handleBlindCMD(userName, message, cmd)
	case "buildBlind":
		return SingleTwitch().handleBlindCMD(userName, message, "билд")
	case "тыктоблять":
		return SingleTwitch().handleBlindCMD(userName, message, cmd)
	case "дискорд":
		return SingleTwitch().handleBlindCMD(userName, message, cmd)
	case "helpBlind":
		return userName + ", Доступные комманды: build, eva, roll, bot, uptime, live, help, master help." +
			" Уникальные на канале: чезаигра, скиллуха, вырубай, билд, вызватьсанитаров, труба, тыктоблять, дискорд." +
			"Взаимодействие с БД: Добавить квоту - addquote <message> или aq <message>, хранить можно даже ссылки." +
			"Получить рандомну квоту из Бд - q или quote. Используйте префикс - " + TwPrefix
	}
	return "Ашибка (handleCMD)"
}

var reactTW = map[string]string{
	"PogChamp": "PogChamp",
	"Kappa 7":  "Kappa 7",
	"Привет":   "MrDestructoid 10000010 10010000 11000010 00100000 01000011 11101000 01100101 00001100",
	"привет":   "MrDestructoid 10000010 10010000 11000010 00100000 01000011 11101000 01100101 00001100",
	"Hello":    "MrDestructoid 10000010 10010000 11000010 00100000 01000011 11101000 01100101 00001100",
	"hello":    "MrDestructoid 10000010 10010000 11000010 00100000 01000011 11101000 01100101 00001100",
	"SMOrc":    "SMOrc",
	"+ в чат":  "+",
}

var reactGG = map[string]string{
	"Привет":  ":skull: 10000010 10010000 11000010 00100000 01000011 11101000 01100101 00001100",
	"привет":  ":skull: 10000010 10010000 11000010 00100000 01000011 11101000 01100101 00001100",
	"Hello":   ":skull: 10000010 10010000 11000010 00100000 01000011 11101000 01100101 00001100",
	"hello":   ":skull: 10000010 10010000 11000010 00100000 01000011 11101000 01100101 00001100",
	"+ в чат": "+",
}

var CMDList = [cmdCOUNT]resource.Commands{
	//Reflyq
	{Command: []string{"вырубить"}, Platform: []string{TW}, Channels: []string{"reflyq"},
		Users: []string{"all"}, Request: "вырубить"},
	{Command: []string{"вырубай"}, Platform: []string{TW}, Channels: []string{"reflyq"},
		Users: []string{"all"}, Request: "вырубайReflyq"},
	{Command: []string{"help", "h", "помощь", "п"}, Platform: []string{TW}, Channels: []string{"reflyq"},
		Users: []string{"all"}, Request: "helpReflyq"},

	//Blind
	{Command: []string{"чезаигра"}, Platform: []string{TW}, Channels: []string{"blindwalkerboy"},
		Users: []string{"all"}, Request: "чезаигра"},
	{Command: []string{"скиллуха"}, Platform: []string{TW}, Channels: []string{"blindwalkerboy"},
		Users: []string{"all"}, Request: "скиллуха"},
	{Command: []string{"вызватьсанитаров"}, Platform: []string{TW}, Channels: []string{"blindwalkerboy"},
		Users: []string{"all"}, Request: "вызватьсанитаров"},
	{Command: []string{"труба", "youtube"}, Platform: []string{TW}, Channels: []string{"blindwalkerboy"},
		Users: []string{"all"}, Request: "труба"},
	{Command: []string{"тыктоблять"}, Platform: []string{TW}, Channels: []string{"blindwalkerboy"},
		Users: []string{"all"}, Request: "тыктоблять"},
	{Command: []string{"дискорд", "discord"}, Platform: []string{TW}, Channels: []string{"blindwalkerboy"},
		Users: []string{"all"}, Request: "дискорд"},
	{Command: []string{"вырубай"}, Platform: []string{TW}, Channels: []string{"blindwalkerboy"},
		Users: []string{"all"}, Request: "вырубайBlind"},
	{Command: []string{"build", "билд", "билдец"}, Platform: []string{TW}, Channels: []string{"blindwalkerboy"},
		Users: []string{"all"}, Request: "buildBlind"},
	{Command: []string{"help", "h", "помощь", "п"}, Platform: []string{TW}, Channels: []string{"blindwalkerboy"},
		Users: []string{"all"}, Request: "helpBlind"},

	//All
	{Command: []string{"build", "билд", "билдец"}, Platform: []string{"all"}, Channels: []string{"all"},
		Users: []string{"all"}, Request: "build"},
	{Command: []string{"eva", "ева"}, Platform: []string{"all"}, Channels: []string{"all"},
		Users: []string{"all"}, Request: "eva"},
	{Command: []string{"roll", "ролл"}, Platform: []string{"all"}, Channels: []string{"all"},
		Users: []string{"all"}, Request: "roll"},
	{Command: []string{"ping", "pong"}, Platform: []string{"all"}, Channels: []string{"all"},
		Users: []string{"all"}, Request: "pong"},
	{Command: []string{"bot", "бот"}, Platform: []string{"all"}, Channels: []string{"all"},
		Users: []string{"all"}, Request: "bot"},
	{Command: []string{"help", "h", "помощь", "п"}, Platform: []string{"all"}, Channels: []string{"all"},
		Users: []string{"all"}, Request: "help"},
	{Command: []string{"master help", "mh"}, Platform: []string{"all"}, Channels: []string{"all"},
		Users: []string{"all"}, Request: "master help"},
	{Command: []string{"uptime", "аптайм"}, Platform: []string{GG, TW}, Channels: []string{"all"},
		Users: []string{"all"}, Request: "uptime"},
	{Command: []string{"live", "жива"}, Platform: []string{"all"}, Channels: []string{"all"},
		Users: []string{"all"}, Request: "live"},
	{Command: []string{"v", "VERSION", "версия"}, Platform: []string{"all"}, Channels: []string{"all"},
		Users: []string{"all"}, Request: "VERSION"},
	{Command: []string{"aq", "addquote"}, Platform: []string{"all"}, Channels: []string{"all"},
		Users: []string{"all"}, Request: "add quote"},
	{Command: []string{"q", "quote"}, Platform: []string{"all"}, Channels: []string{"all"},
		Users: []string{"all"}, Request: "get quote"},
}
