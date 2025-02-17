package callbacks

import (
	"log"

	"github.com/cerdelen/splitWithFriends/globals"
	"github.com/cerdelen/splitWithFriends/keyboards"
	"github.com/cerdelen/splitWithFriends/user"
	"github.com/cerdelen/splitWithFriends/user/userstates"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


func HandleCallBackQueries(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Printf("Received callback query: %+v", update.CallbackQuery.Message.Chat.ID)
	// log.Println("Callback got registered", err)
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
	if _, err := bot.Request(callback); err != nil {
		log.Println("Error acknowledging callback:", err)
	}

	chatID := update.CallbackQuery.Message.Chat.ID
	callbackData := update.CallbackQuery.Data
	userID := update.CallbackQuery.Message.Chat.ID


    switch user.Users[userID].State {
        case userstates.Awaiting_amount_to_split:
        case userstates.Awaiting_for_split_by_amount:
			var responseText string
			var split_by int
			switch callbackData {
			case "split_by_2":
				split_by = 2
			case "split_by_3":
				split_by = 3
			case "split_by_4":
				split_by = 4
			case "split_by_5":
				split_by = 5
			case "split_by_6":
				split_by = 6
			case "split_by_more_than_6":
				split_by = -1
			default:
				// log.Fatal("user was in State to give spplit by amt but has unknown response\nUser ID: %d, callBackData: %s\n", userID, callbackData)
			}
			if split_by > 0 {
				responseText = "What is the amount to be split?"
                user.Users[userID].State = userstates.Awaiting_amount_to_split
				// userStates[userID] = UserState{State: "awaiting_amount_to_split"}
				globals.SplitByValue[userID] = split_by
			} else {
				responseText = "Please enter by how many you want to Split instead?"
				// userStates[userID] = UserState{State: "awaiting_amount_to_split"}
			}
			msg := tgbotapi.NewMessage(chatID, responseText)
			bot.Send(msg)
        case userstates.Awaiting_for_split_contacts:
        case userstates.Awaiting_new_contact_name:
        case userstates.None:
        case userstates.Start:
            switch callbackData {
                case "simple_split":
                    user.Users[userID].State = userstates.Awaiting_for_split_by_amount
                    // userStates[userID] = UserState{State: "waiting_for_split_by_amount"}
                    msg := tgbotapi.NewMessage(chatID, "By how many people do you want to split the Bill?")
                    msg.ReplyMarkup = keyboards.Split_by_amt_keyboard
                    bot.Send(msg)
            }
        case userstates.Configuration:
            switch callbackData {
                case "register_self":
                    // userName := update.CallbackQuery.From.UserName
                    user.RegisterToBotMessages(userID)
                case "deregister_self":
                    user.DeregisterToBotMessages(userID)

            }
        case userstates.RequestFromSingleContact:

        default:
            panic("unexpected userstates.UserState")
	}

	editMsg := tgbotapi.NewEditMessageText(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		"You selected: "+callbackData,
	)
	bot.Send(editMsg)

	// if state, ok := userStates[userID]; ok {
	// } else {
	// 	var responseText string
	// 	switch callbackData {
	// 	case "new_contact":
	// 		responseText = "New Contact"
	// 		userStates[userID] = UserState{State: "awaiting_new_contact_name"}
	// 		// msg := tgbotapi.NewMessage(chatID, )
	// 		responseText = "Please type the UserName you want to add as a Contact!"
	//
 //        case "finished_selecting_contacts":
 //            userStates[userID] = UserState{State: "awaiting_amount_to_split"}
 //            split := currentSplit[userID]
 //            split.divisor = len(split.splitWith)
 //            currentSplit[userID] = split
	// 		msg := tgbotapi.NewMessage(chatID, "What is the amount to be split?")
	// 		bot.Send(msg)
	// 	case "split_with_contacts":
	// 		userStates[userID] = UserState{State: "waiting_for_split_contacts"}
	// 		msg := tgbotapi.NewMessage(chatID, "Which Contacts do you want to split with?")
 //            keyboard, err := keyboards.BuildSplitContactKeyboard(userID)
 //            if err != nil {
 //                if split, exists := currentSplit[userID]; exists {
 //                    if len(split.splitWith) == 0 {
 //                        // TODO No Contacts but also none contacts put so far, ERROR!
 //                    } else {
 //                        // TODO No more contacts but it should be done alrady
 //                    }
 //                }
 //            } else {
 //                msg.ReplyMarkup = keyboard
 //                bot.Send(msg)
 //            }
	// 	case "new_Split":
	// 		currentSplit[userID] = Split{
	// 			author:     userID,
	// 			authorName: update.CallbackQuery.From.UserName,
	// 			divisor:    -1,
	// 			amt:        -1.0,
	// 			splitWith:  []int64{},
	// 		}
	// 		msg := tgbotapi.NewMessage(chatID, "Click a button:")
	// 		msg.ReplyMarkup = keyboards.New_split_keyboard
	// 		bot.Send(msg)
	// 	default:
	// 		responseText = "Unknown option."
	// 	}
	// 	msg := tgbotapi.NewMessage(chatID, responseText)
	// 	bot.Send(msg)
	// }
}

