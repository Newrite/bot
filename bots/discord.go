package bots

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type BotDiscord struct {
	TwitchPtr   *BotTwitch
	GoodGamePtr *BotGoodGame
	Session     *discordgo.Session
	token       string
}

func (db *BotDiscord) Start() {
	var err error
	db.token = "NzE2NjQ1Mzk1MDE0NjE1MDkw.XtOyFw.dG3BnWN9ydODZWGJOrIFlK4lKUk"
	db.Session, err = discordgo.New("Bot " + db.token)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "discordgo.New",
			"file":     "discord.go",
			"body":     "Start",
			"error":    err,
		}).Fatal("Ошибка создания сессии.")
	}
	db.Session.AddHandler(db.messageCreate)
	err = db.Session.Open()
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "bots",
			"function": "Session.Open",
			"file":     "discord.go",
			"body":     "Start",
			"error":    err,
		}).Error("Ошибка открытия соединения.")
		return
	}
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	db.Session.Close()
}

func (db *BotDiscord) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("Дискорд Сервер:<#"+m.GuildID+"> Канал:<#"+m.ChannelID+"> Юзер:", m.Author, "Сообщение:", m.Content)
	if m.Author.ID == s.State.User.ID {
		return
	}
	message := strings.ToLower(m.Content)
	if strings.HasPrefix(message, DisPrefix) {
		msgSl := strings.Fields(message)
		_, err := s.ChannelMessageSend(m.ChannelID, checkCMD("<@"+m.Author.ID+">", m.ChannelID, msgSl[0], "DIS", message))
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "bots",
				"function": "discordgo.New",
				"file":     "discord.go",
				"body":     "messageCreate",
				"error":    err,
			}).Error("Ошибка отправки сообщения")
		}
	}
}
