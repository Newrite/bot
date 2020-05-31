package discord

import (
	"bot/reso"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

const prefix = "~"

type DiscordBot struct {
	Session *discordgo.Session
	token   string
}

func (db *DiscordBot) Start() {
	var err error
	db.token = "NzE2NjQ1Mzk1MDE0NjE1MDkw.XtOyFw.dG3BnWN9ydODZWGJOrIFlK4lKUk"
	db.Session, err = discordgo.New("Bot " + db.token)
	if err != nil {
		fmt.Println(err)
	}
	// Register the messageCreate func as a callback for MessageCreate events.
	db.Session.AddHandler(db.messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = db.Session.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	db.Session.Close()
}

func (db *DiscordBot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("Дискорд Сервер:<#"+ m.GuildID+"> Канал:<#"+m.ChannelID+"> Юзер:",m.Author,"Сообщение:",m.Content)
	if m.Author.ID == s.State.User.ID {
		return
	}
	m.Content = strings.ToLower(m.Content)
	switch m.Content {
	case prefix+"roll":
		s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">, "+strconv.Itoa(rand.Intn(21)))
	case prefix+"eva":
		s.ChannelMessageSend(m.ChannelID, reso.EvaAnswers[rand.Intn(16)])
	case prefix+"билд":
		s.ChannelMessageSend(m.ChannelID, reso.BuildAnswers[rand.Intn(16)])
	case prefix+"help":
		s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">, ~билд,  ~roll, ~eva, ~help, ~ada")
	case prefix+"ada":
		s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">, AdaIsEva, Discord чат бот, написана на" +
			" GoLang v1.14 с использование оболочки DiscordGo by bwmarrin. " +
		"Живет на VPS с убунтой размещенном в москоу сити. Рекомендации, пожелания и" +
		" прочая можно присылать на adaiseva.newrite@gmail.com")
	}
}
