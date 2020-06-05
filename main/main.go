package main

import (
	"bot/bots"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	go bots.SingleTwitch().Start()
	go bots.SingleGoodGame().Start()
	bots.SingleDiscord().Start()
}
