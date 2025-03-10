package callbacks

import (
	"errors"
	// "fmt"
	"log"
	"strconv"

	"github.com/cerdelen/splitWithFriends/globals"
	"github.com/cerdelen/splitWithFriends/keyboards"
	"github.com/cerdelen/splitWithFriends/split"
	"github.com/cerdelen/splitWithFriends/user"
	"github.com/cerdelen/splitWithFriends/user/userstates"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	botwrap "github.com/cerdelen/splitWithFriends/bot"
)

func parseSplitByAmount(update tgbotapi.Update) (string, tgbotapi.InlineKeyboardMarkup, error) {
	// chatID := update.CallbackQuery.Message.Chat.ID
	callbackData := update.CallbackQuery.Data
	userID := update.CallbackQuery.Message.Chat.ID
	var responseText string
	var split_by int
	var err error = nil
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
	return responseText, tgbotapi.NewInlineKeyboardMarkup(), err
}

func parseDirectRequest(update tgbotapi.Update) (string, tgbotapi.InlineKeyboardMarkup, error) {
	callbackData := update.CallbackQuery.Data
	userID := update.CallbackQuery.Message.Chat.ID
	var responseText string
	var keyboard tgbotapi.InlineKeyboardMarkup
	var err error = nil
	switch callbackData {
	case "load_prev_contacts":
		user.Users[userID].ContactIndexing -= 4
		keyboard, err = keyboards.BuildContactKeyboard(userID)
		responseText = "What User do you want to send a Request"
	case "load_more_contacts":
		user.Users[userID].ContactIndexing += 4
		keyboard, err = keyboards.BuildContactKeyboard(userID)
		responseText = "What User do you want to send a Request"
	case "finished_selecting_contacts":
		user.Users[userID].State = userstates.None
		responseText = "Canceled Request"
	default:
		var otherId int
		if otherId, err = strconv.Atoi(callbackData); err == nil {
			if !user.Users[userID].HasContact(int64(otherId)) {
				err = errors.New("User is not a Contact")
			} else {
				me := user.Users[userID]
				recipient := user.Users[int64(otherId)]
				split.DirectRequestsInit(me, recipient)
				// split.CurrentDirectRequests[userID] = split.DirectRequestsInit(me, recipient)
				me.State = userstates.AwaitingAmountDirectRequest
			}
		}
	}
	return responseText, keyboard, err
}

func parseStartCallBack(update tgbotapi.Update) (string, tgbotapi.InlineKeyboardMarkup, error) {
	// chatID := update.CallbackQuery.Message.Chat.ID
	callbackData := update.CallbackQuery.Data
	userID := update.CallbackQuery.Message.Chat.ID
	var responseText string
	var keyboard tgbotapi.InlineKeyboardMarkup
	var err error = nil
	switch callbackData {
	case "configuration":
		user.Users[userID].State = userstates.Configuration
		responseText = "Configuration"
		keyboard = keyboards.ConfigurationKeyboard
	case "new_split":
		user.Users[userID].State = userstates.NewSplit
		responseText = "New Split"
		keyboard = keyboards.NewSplitKeyboard
	case "direct_request":
		user.Users[userID].State = userstates.NewDirectRequest
		if keyboard, err = keyboards.BuildContactKeyboard(userID); err != nil {
			log.Println(err.Error())
			responseText = "Error Parsing Contact Keyboard for direct request"
		} else {
			responseText = "Direct Request"
		}
	}
	return responseText, keyboard, err
}

func parseNewSplit(update tgbotapi.Update) (string, tgbotapi.InlineKeyboardMarkup, error) {
	callbackData := update.CallbackQuery.Data
	userID := update.CallbackQuery.Message.Chat.ID
	var responseText string
	var keyboard tgbotapi.InlineKeyboardMarkup
	var err error = nil
	switch callbackData {
	case "split_with_contacts":
		user.Users[userID].State = userstates.AddContactsToSplit
		user.Users[userID].ContactIndexing = 0
		split.CurrentSplits[userID] = split.Init(user.Users[userID])
		if keyboard, err = keyboards.BuildAddSplitterKeyboard(userID); err != nil {
			log.Println(err.Error())
			responseText = "Error Parsing Adding Contact Keyboard"
		} else {
			responseText = "What User do you want to add as a Contact?"
		}
	default:
		log.Fatalf("In parseNewSplit calbackdata is %s", callbackData)
	}
	return responseText, keyboard, err
}

func parseConfigurationCallBack(update tgbotapi.Update) (string, tgbotapi.InlineKeyboardMarkup, error) {
	callbackData := update.CallbackQuery.Data
	userID := update.CallbackQuery.Message.Chat.ID
	var responseText string
	var keyboard tgbotapi.InlineKeyboardMarkup
	var err error = nil
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
	case "currency_change":
		user.Users[userID].State = userstates.None
		currency := user.Users[userID].ChangeCurrency()
		responseText = "Changed Currency to " + currency
	case "add_contact":
		user.Users[userID].State = userstates.AddingContact
		user.Users[userID].ContactIndexing = 0
		if keyboard, err = keyboards.BuildAddingContactKeyboard(userID); err != nil {
			log.Println(err.Error())
			responseText = "Error Parsing Adding Contact Keyboard"
		} else {
			responseText = "What User do you want to add as a Contact?"
		}
	case "remove_contact":
		user.Users[userID].State = userstates.RemovingContact
		user.Users[userID].ContactIndexing = 0
		if keyboard, err = keyboards.BuildRemovingContactKeyboard(userID); err != nil {
			log.Println(err.Error())
			responseText = "Error Parsing Removing Contact Keyboard"
		} else {
			responseText = "What User do you want to remove as a Contact?"
		}
	}
	return responseText, keyboard, err
}

func parseRemoveContact(update tgbotapi.Update) (string, tgbotapi.InlineKeyboardMarkup, error) {
	callbackData := update.CallbackQuery.Data
	userID := update.CallbackQuery.Message.Chat.ID

	var responseText string
	var keyboard tgbotapi.InlineKeyboardMarkup
	var err error
	switch callbackData {
	case "load_prev_contacts":
		user.Users[userID].ContactIndexing -= 4
		keyboard, err = keyboards.BuildRemovingContactKeyboard(userID)
		responseText = "What User do you want to remove as a Contact?"
	case "load_more_contacts":
		user.Users[userID].ContactIndexing += 4
		keyboard, err = keyboards.BuildRemovingContactKeyboard(userID)
		responseText = "What User do you want to remove as a Contact?"
	case "finished_selecting_contacts":
		user.Users[userID].State = userstates.None
		responseText = "Finished Removing Contacts"
	default:
		var otherId int
		if otherId, err = strconv.Atoi(callbackData); err == nil {
			user.Users[userID].RemoveContact(int64(otherId))
			if keyboard, err = keyboards.BuildRemovingContactKeyboard(userID); err != nil {
				log.Println(err.Error())
				responseText = "Error Parsing Contact Keyboard"
			} else {
				responseText = "What User do you want to remove as a Contact?"
			}
		}
	}
	return responseText, keyboard, nil
}

func parseAddingSplitter(update tgbotapi.Update) (string, tgbotapi.InlineKeyboardMarkup, error) {
	callbackData := update.CallbackQuery.Data
	userID := update.CallbackQuery.Message.Chat.ID

	var responseText string
	var keyboard tgbotapi.InlineKeyboardMarkup
	var err error = nil
	switch callbackData {
	case "load_prev_contacts":
		user.Users[userID].ContactIndexing -= 4
		keyboard, err = keyboards.BuildAddSplitterKeyboard(userID)
		responseText = "What User do you want to add to the Split"
	case "load_more_contacts":
		user.Users[userID].ContactIndexing += 4
		keyboard, err = keyboards.BuildAddSplitterKeyboard(userID)
		responseText = "What User do you want to add to the Split"
	case "finished_selecting_contacts":
		user.Users[userID].State = userstates.None
		responseText = "Finished Adding Contacts"
	default:
		var otherId int
		if otherId, err = strconv.Atoi(callbackData); err == nil {
			// if err = user.Users[userID].AddContactToSplit(int64(otherId)); err == nil {
			if err = split.CurrentSplits[userID].AddSplitter(user.Users[int64(otherId)]); err == nil {
				keyboard, err = keyboards.BuildAddSplitterKeyboard(userID)
				responseText = "What User do you want to add to the Split"
			} else {
				log.Printf("Adding %d as a Splitter to %d's Split\n%s", userID, otherId, err.Error())
			}
		}
	}
	return responseText, keyboard, err
}

func parseAddingContact(update tgbotapi.Update) (string, tgbotapi.InlineKeyboardMarkup, error) {
	callbackData := update.CallbackQuery.Data
	userID := update.CallbackQuery.Message.Chat.ID

	var responseText string
	var keyboard tgbotapi.InlineKeyboardMarkup
	var err error
	switch callbackData {
	case "load_prev_contacts":
		user.Users[userID].ContactIndexing -= 4
		keyboard, err = keyboards.BuildAddingContactKeyboard(userID)
		responseText = "What User do you want to add as a Contact?"
	case "load_more_contacts":
		user.Users[userID].ContactIndexing += 4
		keyboard, err = keyboards.BuildAddingContactKeyboard(userID)
		responseText = "What User do you want to add as a Contact?"
	case "finished_selecting_contacts":
		user.Users[userID].State = userstates.None
		responseText = "Finished Adding Contacts"
	default:
		var otherId int
		if otherId, err = strconv.Atoi(callbackData); err == nil {
			if err = user.Users[userID].AddContact(int64(otherId)); err == nil {
				keyboard, err = keyboards.BuildAddingContactKeyboard(userID)
				responseText = "What User do you want to add as a Contact?"
				if user.Users[userID].ContactIndexing > 0 {
					user.Users[userID].ContactIndexing -= 4
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
	case userstates.AddContactsToSplit:
		msg, keyboard, err = parseAddingSplitter(update)
	case userstates.AddingContact:
		msg, keyboard, err = parseAddingContact(update)
	case userstates.Start:
		msg, keyboard, err = parseStartCallBack(update)
	case userstates.Configuration:
		msg, keyboard, err = parseConfigurationCallBack(update)
	case userstates.RemovingContact:
		msg, keyboard, err = parseRemoveContact(update)
	case userstates.NewDirectRequest:
		msg, keyboard, err = parseDirectRequest(update)
	case userstates.NewSplit:
		msg, keyboard, err = parseNewSplit(update)
	case userstates.DirectRequest:

	default:
		panic("unexpected userstates.UserState")
	}

	log.Printf("Keyboard: %+v", keyboard)

	if err == nil {
        botwrap.EditMessageAndKeyboard(
			update.CallbackQuery.Message.Chat.ID,
			update.CallbackQuery.Message.MessageID,
            msg,
            keyboard,
        )
	} else {
		log.Printf("Error in HandleCallBackQueries: %s", err.Error())
	}
}
