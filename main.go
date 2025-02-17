package main

import (
	// "errors"
	// "fmt"
	"log"
	"os"

	// "regexp"
	// "strconv"

	"github.com/cerdelen/splitWithFriends/globals"
	"github.com/cerdelen/splitWithFriends/keyboards"
	"github.com/cerdelen/splitWithFriends/updates/callbacks"
	"github.com/cerdelen/splitWithFriends/updates/messages"
	"github.com/cerdelen/splitWithFriends/user"
	"github.com/cerdelen/splitWithFriends/user/userstates"

	// "github.com/cerdelen/splitWithFriends/user/userstates"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)


const MAXRETRIES = 3

type UserState struct {
	State string
}


// var userStates = make(map[int64]UserState)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func checkRetryLeft(userID int64) int {
	if retriesLeft, exists := globals.RetryCounter[userID]; exists {
		retriesLeft -= 1
		globals.RetryCounter[userID] = retriesLeft
		return retriesLeft
	}
	globals.RetryCounter[userID] = MAXRETRIES
	return MAXRETRIES
}

func removeRetries(userID int64) {
	delete(globals.RetryCounter, userID)
}

// func checkUserNameExists(userName string) bool {
// 	for _, v := range globals.RegisteredUsers {
// 		if v == userName {
// 			return true
// 		}
// 	}
// 	return false
// }

func returnHelpMessage (bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    userID := update.Message.Chat.ID
    msg := tgbotapi.NewMessage(userID, "Start the Bot by sending \"/start\" and follow the instructions!")
    bot.Send(msg)
}


func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN not set")
	} else {
		log.Println(botToken)
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	botInfo, err := bot.GetMe()
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Bot Username: %s", botInfo.UserName)
	log.Printf("Can Join Groups: %v", botInfo.CanJoinGroups)
	log.Printf("Can Read Messages: %v", botInfo.CanReadAllGroupMessages)
	log.Printf("Supports Inline Queries: %v", botInfo.SupportsInlineQueries)

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.AllowedUpdates = []string{"message", "edited_channel_post", "callback_query"}
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		log.Printf("Received update: %+v", update)
        if update.CallbackQuery != nil {
            callbacks.HandleCallBackQueries(bot, update)
        } else if update.Message != nil {
            userID := update.Message.Chat.ID
            user.AddIfNewUser(userID, update.Message.Chat.UserName)
            switch user.Users[userID].State {
                case userstates.None:
                    switch update.Message.Text {
                    case "/start":
                        user.Users[userID].State = userstates.Start
                        // userStates[userID] = UserState{State: "start"}
                        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Click a button:")
                        msg.ReplyMarkup = keyboards.StartKeyboard
                        bot.Send(msg)
                    default:
                        returnHelpMessage(bot, update)
                    }
                default:
                    messages.HandleMessage(bot, update)
            }
		}
	}
}
