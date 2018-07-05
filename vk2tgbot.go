package main

import (
	"log"
	"fmt"
	"os"
	"strings"
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
	msgTxt := update.Message.CommandArguments()
	msgParts := strings.Split(msgTxt, "_")
	if len(msgParts) < 3 {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "/newvote command take a vote name and alternatives as arguments, see /help"))
	} else {
		replay_txt := `You specified a vote ` + msgParts[1] + ` 
		with variants:` + "\n" + strings.Join(msgParts[2:], "\n")
		replay := tgbotapi.NewMessage(update.Message.Chat.ID, replay_txt)
		replay.ReplyToMessageID = update.Message.MessageID
		bot.Send(replay)
	}
}

func repli(bot tgbotapi.BotAPI, update tgbotapi.Update) {
		msg_common := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg_common.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg_common)
	}

func main() {
	var c conf
    c.getConf()
	
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
		
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		
		if update.Message.IsCommand() {
			fmt.Println("This is commannddd!!!")
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "help":
				msg.Text = "type /start or /newvote."				
			case "sayhi":
				msg.Text = "Hi :)"				
			case "start":
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, getBotStartMsg(bot.Self.UserName)))				
			case "newvote":
				newVote(*bot, update)				
			default:
				msg.Text = "I don't know that command"
			} 
			bot.Send(msg)
		} else { 
			repli(*bot, update) 
		}		
	}
}