package main

import (
	"bot/bot"
	"math/rand"
	"time"
)

func main() {
	var twitchBot bot.TwitchBot
	rand.Seed(time.Now().Unix())
	twitchBot.Start()
}
