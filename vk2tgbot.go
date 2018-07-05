package main

import (
	"log"
	"fmt"
	"os"
	"io/ioutil"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"gopkg.in/yaml.v2"
)

type conf struct {
    Token string `yaml:"token"`
	Use_proxy bool `yaml:"use_proxy"`
}


func check(e error) {
    if e != nil {
        log.Panic(e)
    }
}

func (c *conf) getConf() *conf {
	config, err := ioutil.ReadFile("config.yml")
	check(err)
	err = yaml.Unmarshal(config, c)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }

    return c
}

func getBotStartMsg(botName string) string {

	botMsg := `Welcome to @` + botName + `!
This is start message which needs to be updated
This is a super druper bot which needs to be developed`
	return botMsg
}

func newVote(bot tgbotapi.BotAPI, update tgbotapi.Update) {
	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, `Specify the name of the Vote`))
	replay :=  tgbotapi.NewMessage(update.Message.Chat.ID, `You specified "` + update.Message.Text + `"`)
	replay.ReplyToMessageID = update.Message.MessageID
	bot.Send(replay)
}

func main() {
	var c conf
    c.getConf()
	fmt.Printf(c.Token)
	if c.Token == "" { 
		log.Fatalf("There is no token!")
		os.Exit(1)
	}
	bot, err := tgbotapi.NewBotAPI(c.Token)
	check(err)

//	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	
	updates, err := bot.GetUpdatesChan(u)
	check(err)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.Text == "/start" || update.Message.Text == "/start@"+bot.Self.UserName {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, getBotStartMsg(bot.Self.UserName)))
			continue
		} else if update.Message.Text == "/newvote" || update.Message.Text == "/newvote@"+bot.Self.UserName {
			newVote(*bot, update)			
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}