package main

import (
	"bot/bots"
	_ "github.com/mattn/go-sqlite3"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	go bots.SingleTwitch().Start()
	go bots.SingleGoodGame().Start()
	//go bots.SingleMixer().Start()
	bots.SingleDiscord().Start()
}
