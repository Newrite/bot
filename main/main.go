package main

import (
	"bot/goodgame"
	"bot/twitch"
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
