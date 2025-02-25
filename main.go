package main

import (
	"log"
	"os"

	"github.com/cerdelen/splitWithFriends/keyboards"
	"github.com/cerdelen/splitWithFriends/updates/callbacks"
	"github.com/cerdelen/splitWithFriends/updates/messages"
	"github.com/cerdelen/splitWithFriends/user"
	"github.com/cerdelen/splitWithFriends/user/userstates"

	botWrap "github.com/cerdelen/splitWithFriends/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

// var botWrap.Bot *tgbotapi.BotAPI

type UserState struct {
	State string
}

func returnHelpMessage (update tgbotapi.Update) {
    userID := update.Message.Chat.ID
    msg := tgbotapi.NewMessage(userID, "Start the botWrap.Bot by sending \"/start\" and follow the instructions!")
    botWrap.Bot.Send(msg)
}

func main() {
	err := godotenv.Load()

	if err != nil { log.Fatal("Error loading .env file") }

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN not set")
	}

	botWrap.Bot, err = tgbotapi.NewBotAPI(botToken)

	if err != nil {
		log.Panic(err)
	}

	botWrap.Bot.Debug = true

	// botInfo, err := botWrap.Bot.GetMe()
	// if err != nil {
	// 	log.Panic(err)
	// }
	// log.Printf("botWrap.Bot Username: %s", botInfo.UserName)
	// log.Printf("Can Join Groups: %v", botInfo.CanJoinGroups)
	// log.Printf("Can Read Messages: %v", botInfo.CanReadAllGroupMessages)
	// log.Printf("Supports Inline Queries: %v", botInfo.SupportsInlineQueries)

	log.Printf("Authorized on account %s", botWrap.Bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.AllowedUpdates = []string{"message", "edited_channel_post", "callback_query"}
	u.Timeout = 60

	updates := botWrap.Bot.GetUpdatesChan(u)

    user.AddIfNewUser(123, "123Name")
    user.AddIfNewUser(321, "321Name")
    user.AddIfNewUser(1234, "1234Name")
    user.AddIfNewUser(4321, "4321Name")
    user.AddIfNewUser(12345, "12345Name")
    user.AddIfNewUser(54321, "54321Name")
    user.AddIfNewUser(123456, "123456Name")
    user.AddIfNewUser(654321, "654321Name")
    user.AddIfNewUser(1234567, "1234567Name")
    user.AddIfNewUser(7654321, "7654321Name")
    user.AddIfNewUser(12345678, "12345678Name")
    user.AddIfNewUser(87654321, "87654321Name")
    user.RegisterToBotMessages(321)
    user.RegisterToBotMessages(4321)
    user.RegisterToBotMessages(54321)
    user.RegisterToBotMessages(654321)
    user.RegisterToBotMessages(7654321)
    user.RegisterToBotMessages(87654321)

    user.AddIfNewUser(1233, "BBC")
    user.RegisterToBotMessages(1233)
    user.AddIfNewUser(12334, "CBC")
    user.RegisterToBotMessages(12334)

    user.AddIfNewUser(54362, "abc")
    user.RegisterToBotMessages(54362)
    user.AddIfNewUser(123, "ABC")
    user.RegisterToBotMessages(123)


    var meID int64 = 476541234
    user.AddIfNewUser(meID, "Cerdelen")
    user.RegisterToBotMessages(meID)


	for update := range updates {
        log.Println("")
		log.Printf("Received update: %+v", update)
        if update.CallbackQuery != nil {
            callbacks.HandleCallBackQueries(botWrap.Bot, update)
        } else if update.Message != nil {
            userID := update.Message.Chat.ID
            user.AddIfNewUser(userID, update.Message.Chat.UserName)
            switch user.Users[userID].State {
                case userstates.None:
                    switch update.Message.Text {
                        case "/start":
                            user.Users[userID].State = userstates.Start
                            msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Click a button:")
                            msg.ReplyMarkup = keyboards.StartKeyboard
                            botWrap.Bot.Send(msg)
                        default:
                            returnHelpMessage(update)
                    }
                default:
                    messages.HandleMessage(botWrap.Bot, update)
            }
		}
	}
}

