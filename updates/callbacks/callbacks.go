package callbacks

import (
	"log"
	"strconv"

	"github.com/cerdelen/splitWithFriends/globals"
	"github.com/cerdelen/splitWithFriends/keyboards"
	"github.com/cerdelen/splitWithFriends/user"
	"github.com/cerdelen/splitWithFriends/user/userstates"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


func parseSplitByAmount(update tgbotapi.Update) (string, tgbotapi.InlineKeyboardMarkup, error) {
        // chatID := update.CallbackQuery.Message.Chat.ID
        callbackData := update.CallbackQuery.Data
        userID := update.CallbackQuery.Message.Chat.ID
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
        return responseText, tgbotapi.NewInlineKeyboardMarkup(), nil
}

func parseStartCallBack(update tgbotapi.Update) (string, tgbotapi.InlineKeyboardMarkup, error) {
	// chatID := update.CallbackQuery.Message.Chat.ID
	callbackData := update.CallbackQuery.Data
	userID := update.CallbackQuery.Message.Chat.ID
    var responseText string
    var keyboard tgbotapi.InlineKeyboardMarkup
    switch callbackData {
        case "configuration":
            user.Users[userID].State = userstates.Configuration
            // userStates[userID] = UserState{State: "waiting_for_split_by_amount"}
            // msg := tgbotapi.NewMessage(chatID, "Configuration")
            // bot.Send(msg)
            responseText = "Configuration"
            keyboard = keyboards.ConfigurationKeyboard
        case "new_split":
        case "direct_request":
        // case "simple_split":
        //     user.Users[userID].State = userstates.Awaiting_for_split_by_amount
        //     // userStates[userID] = UserState{State: "waiting_for_split_by_amount"}
        //     msg := tgbotapi.NewMessage(chatID, "By how many people do you want to split the Bill?")
        //     msg.ReplyMarkup = keyboards.Split_by_amt_keyboard
        //     bot.Send(msg)
    }
        return responseText, keyboard, nil
}

func parseConfigurationCallBack(update tgbotapi.Update) (string, tgbotapi.InlineKeyboardMarkup, error) {
	// chatID := update.CallbackQuery.Message.Chat.ID
	callbackData := update.CallbackQuery.Data
	userID := update.CallbackQuery.Message.Chat.ID
    var responseText string
    var keyboard tgbotapi.InlineKeyboardMarkup
    var err error
    switch callbackData {
        case "register_self":
            if err = user.RegisterToBotMessages(userID); err == nil {
                responseText = "Successfully registered"
                user.Users[userID].State = userstates.None
            }
        case "deregister_self":
            if err = user.DeregisterToBotMessages(userID); err == nil {
                responseText = "Successfully deregistered"
                user.Users[userID].State = userstates.None
            }
        case "add_contact":
            user.Users[userID].State = userstates.AddingContact
            // var msg tgbotapi.MessageConfig
            if keyboard, err = keyboards.BuildAddingContactKeyboard(userID); err != nil {
                log.Println(err.Error())
                responseText = "Error Parsing Adding Contact Keyboard"
            } else {
                responseText = "What User do you want to add as a Contact?"
                // msg.ReplyMarkup = keyboard
                // log.Printf("%+v", keyboard)
            }
        case "remove_contact":
            user.Users[userID].State = userstates.RemovingContact
            // var msg tgbotapi.MessageConfig
            if keyboard, err = keyboards.BuildRemovingContactKeyboard(userID); err != nil {
                log.Println(err.Error())
                responseText = "Error Parsing Removing Contact Keyboard"
            } else {
                responseText = "What User do you want to remove as a Contact?"
                // msg.ReplyMarkup = keyboard
                // log.Printf("%+v", keyboard)
            }
    }
    return responseText, keyboard, nil
}

func parseRemoveContact(update tgbotapi.Update) (string, tgbotapi.InlineKeyboardMarkup, error) {
	callbackData := update.CallbackQuery.Data
	userID := update.CallbackQuery.Message.Chat.ID

    var responseText string
    var keyboard tgbotapi.InlineKeyboardMarkup
    var err error
    switch callbackData {
        case "load_more_contacts":
            // user.Users[userID].State = userstates.AddingContact
        case "finished_selecting_contacts":
            user.Users[userID].State = userstates.None
        default:
            var otherId int
            if otherId, err = strconv.Atoi(callbackData); err == nil {
                user.Users[userID].RemoveContact(int64(otherId))
                if keyboard, err = keyboards.BuildRemovingContactKeyboard(userID); err != nil {
                    log.Println(err.Error())
                    responseText = "Error Parsing Contact Keyboard"
                } else {
                    responseText = "What User do you want to add as a Contact?"
                }
            }
    }
    return responseText, keyboard, nil
}

func parseAddingContact(update tgbotapi.Update) (string, tgbotapi.InlineKeyboardMarkup, error) {
	callbackData := update.CallbackQuery.Data
	userID := update.CallbackQuery.Message.Chat.ID

    var responseText string
    var keyboard tgbotapi.InlineKeyboardMarkup
    var err error
    switch callbackData {
        case "load_more_contacts":
            // user.Users[userID].State = userstates.AddingContact
        case "finished_selecting_contacts":
            user.Users[userID].State = userstates.None
        default:
            var otherId int
            if otherId, err = strconv.Atoi(callbackData); err == nil {
                if err = user.Users[userID].AddContact(int64(otherId)); err == nil {
                    if keyboard, err = keyboards.BuildAddingContactKeyboard(userID); err != nil {
                        log.Println(err.Error())
                        responseText = "Error Parsing Contact Keyboard"
                    } else {
                        responseText = "What User do you want to add as a Contact?"
                    }
                } else {
                    log.Printf("Adding %d as a Contact to %d\n%s", userID, otherId, err.Error())
                }
            }
    }
    return responseText, keyboard, nil
}

func HandleCallBackQueries(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Printf("Received callback query: %+v", update.CallbackQuery.Message.Chat.ID)
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
	if _, err := bot.Request(callback); err != nil {
		log.Println("Error acknowledging callback:", err)
	}

	userID := update.CallbackQuery.Message.Chat.ID

    var msg string
    var err error
    var keyboard tgbotapi.InlineKeyboardMarkup

    switch user.Users[userID].State {
        case userstates.Awaiting_amount_to_split:
        case userstates.Awaiting_for_split_by_amount:
            msg, keyboard, err = parseSplitByAmount(update)
        case userstates.Awaiting_for_split_contacts:
        case userstates.Awaiting_new_contact_name:
        case userstates.None:
        case userstates.AddingContact:
            msg, keyboard, err = parseAddingContact(update)
        case userstates.Start:
            msg, keyboard, err = parseStartCallBack(update)
        case userstates.Configuration:
            msg, keyboard, err = parseConfigurationCallBack(update)
        case userstates.RequestFromSingleContact:
        case userstates.RemovingContact:
            msg, keyboard, err = parseRemoveContact(update)

        default:
            panic("unexpected userstates.UserState")
	}

    if err == nil {
        editMsg := tgbotapi.NewEditMessageText(
        	update.CallbackQuery.Message.Chat.ID,
        	update.CallbackQuery.Message.MessageID,
            msg,
        )
        bot.Send(editMsg)
        editKeyboard := tgbotapi.NewEditMessageReplyMarkup(
        	update.CallbackQuery.Message.Chat.ID,
        	update.CallbackQuery.Message.MessageID,
            keyboard,
        )
        bot.Send(editKeyboard)
    } else {
    }

}

