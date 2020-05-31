package main

import (
	"bot/discord"
	"bot/goodgame"
	"bot/twitch"
	"math/rand"
	"time"
)

type Bot struct {
	TBot  *twitch.BotTwitch
	GGBot *goodgame.BotGoodGame
	DBot  *discord.DiscordBot
}

func main() {
	bot := &Bot{
		&twitch.BotTwitch{},
		&goodgame.BotGoodGame{},
		&discord.DiscordBot{},
	}
	rand.Seed(time.Now().Unix())
	bot.GGBot.TwitchPtr = bot.TBot
	go bot.GGBot.Start()
	go bot.TBot.Start()
	bot.DBot.Start()
}
