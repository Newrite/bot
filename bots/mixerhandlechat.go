package bots

import (
	"bot/resource"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func (bm *BotMixer) listenChannel() error {
	_, message, err := bm.Connection.ReadMessage()
	if err != nil {
		return err
	}
	if strings.HasPrefix(string(message), `{"type":"event","event":"ChatMessage","data"`) {
		err = json.Unmarshal(message, &bm.Message)
		if err != nil {
			return err
		}
		var str string
		for _, text := range bm.Message.Data.Message.Message {
			str += text.Text
		}
		fmt.Print("[" + timeStamp() + "] [MIXER] Канал: " + "Newrite" + " " +
			"Ник: " + bm.Message.Data.User_name + "\tСообщение: " + str + "\n")
		if !strings.HasPrefix(str, MixPrefix) {
			SingleTwitch().MarkovChain += " " + resource.ReadTxt(str)
		}
		lowMessage := strings.ToLower(str)
		if strings.HasPrefix(lowMessage, MixPrefix) {
			msgSl := strings.Fields(lowMessage)
			go bm.say(checkCMD(bm.Message.Data.User_name, "Newrite", msgSl[0], MIX, lowMessage, str, ""))
		}
	} else {
		fmt.Println(string(message))
	}
	time.Sleep(100 * time.Millisecond)
	return nil
}
