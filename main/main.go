package main

import (
	"bot/twitch"
	"bot/goodgame"
	"math/rand"
	"time"
)

func main() {
	var twitchBot twitch.TwitchBot
	var goodGameBot goodgame.GoodGameBot
	rand.Seed(time.Now().Unix())
	go goodGameBot.Start()
	twitchBot.Start()
}
