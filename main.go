package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var new_request_keyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Add new Contact", "new_contact"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Allow Bot Messages", "register_self"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Remove Bot Messages", "deregister_self"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Start new Split", "new_Split"),
	),
)

var new_split_keyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Split with Contacts", "split_with_contacts"),
		tgbotapi.NewInlineKeyboardButtonData("Simple Split", "simple_split"),
	),
)

var split_by_amt_keyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("2", "split_by_2"),
		tgbotapi.NewInlineKeyboardButtonData("3", "split_by_3"),
		tgbotapi.NewInlineKeyboardButtonData("4", "split_by_4"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("5", "split_by_5"),
		tgbotapi.NewInlineKeyboardButtonData("6", "split_by_6"),
		tgbotapi.NewInlineKeyboardButtonData("more", "split_by_more_than_6"),
	),
)

func isValidPositiveNumber(s string) bool {
	re := regexp.MustCompile(`^[0-9]+(\.[0-9]{1,2})?$`)
	return re.MatchString(s)
}

const MAXRETRIES = 3

type UserState struct {
	State string
}

type Split struct {
	author     int64
	authorName string
	divisor    int
	amt        float64
	splitWith  []int64
}

var userStates = make(map[int64]UserState)
var registeredUsers = make(map[int64]string)
var splitByValue = make(map[int64]int)
var retryCounter = make(map[int64]int)
var currentSplit = make(map[int64]Split)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func checkRetryLeft(userID int64) int {
	if retriesLeft, exists := retryCounter[userID]; exists {
		retriesLeft -= 1
		retryCounter[userID] = retriesLeft
		return retriesLeft
	}
	retryCounter[userID] = MAXRETRIES
	return MAXRETRIES
}

func removeRetries(userID int64) {
	delete(retryCounter, userID)
}

func isValidSplit(split Split) bool {
	if split.divisor > 0 && split.amt > 0 {
		return true
	}
	return false
}

func getUserName(userID int64) (string, error) {
	if name, exists := registeredUsers[userID]; exists {
		return name, nil
	}
	return "", errors.New("userId is not registered")

}

func ssplit(bot *tgbotapi.BotAPI, update tgbotapi.Update, split Split) {
	userID := update.Message.Chat.ID
	if isValidSplit(split) {
		divisor := max(split.divisor, len(split.splitWith))
		if len(split.splitWith) == 0 {
			amt := split.amt / float64(divisor)
			responseText := fmt.Sprintf("Amount everyone has to Pay is: %.2f", amt)
			msg := tgbotapi.NewMessage(userID, responseText)
			bot.Send(msg)
		} else {
			// userName := getUserName(split.author)

			for _, splitter := range split.splitWith {
				// splitter_name, err := getUserName(splitter)
				if isRegistered(splitter) {
					amt := split.amt / float64(divisor)
					responseText := fmt.Sprintf("%s splits a Bill of %.2f by %d People.\nThe Amount you have to pay is %.2f", split.authorName, split.amt, amt)
					msg := tgbotapi.NewMessage(splitter, responseText)
					bot.Send(msg)
				} else {
					log.Printf("Unregistered Splitter (i guess unregistered while Split was still assembled: %d", splitter)
				}
				// if err != nil {
				// } else {
				// }
			}
		}

	} else {
		log.Println("Invalid Split Object")
	}
}

func split(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.Chat.ID
	text := update.Message.Text
	var responseText string
	if isValidPositiveNumber(text) {
		num, err := strconv.ParseFloat(text, 64)
		if err != nil || num < 0 {
			retries := checkRetryLeft(userID)
			if retries == 0 {
				responseText = fmt.Sprintf("The value entered is not a valid Amount!\nNo retries Left")
				delete(userStates, userID)
			} else {
				responseText = fmt.Sprintf("The value entered is not a valid Amount!\nYou have %d retries left!", retries)
			}
			responseText = fmt.Sprintf("The value entered is not a valid Amount!\n You have %d retries left!", retries)
		} else if divisor, ok := splitByValue[userID]; ok {
			amt := num / float64(divisor)
			responseText = fmt.Sprintf("Amount everyone has to Pay is: %.2f", amt)
			delete(userStates, userID)
		}
	} else {
		retries := checkRetryLeft(userID)
		if retries == 0 {
			responseText = fmt.Sprintf("The value entered is not a valid Amount!\nNo retries Left")
			delete(userStates, userID)
		} else {
			responseText = fmt.Sprintf("The value entered is not a valid Amount!\nYou have %d retries left!", retries)
		}
	}
	msg := tgbotapi.NewMessage(userID, responseText)
	bot.Send(msg)
}

func registerUser(userID int64, userName string) error {
	log.Printf("registerUser %d, %s", userID, userName)
	if _, exists := registeredUsers[userID]; exists {
		return errors.New("userId is already registered")
	}

	for _, v := range registeredUsers {
		if v == userName {
			return errors.New("userName is already Taken")
		}
	}

	registeredUsers[userID] = userName

	return nil
}
func deregisterUser(userID int64) error {
	log.Printf("deregisterUser %d", userID)
	delete(registeredUsers, userID)
	return nil
}

func checkUserNameExists(userName string) bool {
	for _, v := range registeredUsers {
		if v == userName {
			return true
		}
	}
	return false
}

func isRegistered(userID int64) bool {
	_, exists := registeredUsers[userID]
	return exists
}

func printUserStates() {
	fmt.Println("Current state of userStates:")
	for userID, state := range userStates {
		fmt.Printf("User ID: %d, State: %s\n", userID, state.State)
	}
}

func buildSplitContactKeyboard(userID int64) tgbotapi.InlineKeyboardMarkup {
    log.Println("BUILD SPLIT CONTACT KEYBOARD CALLED")
	var keyboardRows [][]tgbotapi.InlineKeyboardButton

	for contactID, contactName := range registeredUsers {
        if contactID != userID {
            row := tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonData(contactName, strconv.FormatInt(contactID, 10)),
            )
            keyboardRows = append(keyboardRows, row)
        }
	}

    log.Printf("This is the Keyboard %+v", keyboardRows)

	return tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)
}

func callBackQueries(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Printf("Received callback query: %+v", update.CallbackQuery.Message.Chat.ID)
	// log.Println("Callback got registered", err)
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
	if _, err := bot.Request(callback); err != nil {
		log.Println("Error acknowledging callback:", err)
	}

	chatID := update.CallbackQuery.Message.Chat.ID
	callbackData := update.CallbackQuery.Data
	userID := update.CallbackQuery.Message.Chat.ID

	if state, ok := userStates[userID]; ok {
		switch state.State {
		case "waiting_for_split_by_amount":
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
				userStates[userID] = UserState{State: "awaiting_amount_to_split"}
				splitByValue[userID] = split_by
			} else {
				responseText = "Please enter by how many you want to Split instead?"
				// userStates[userID] = UserState{State: "awaiting_amount_to_split"}
			}
			msg := tgbotapi.NewMessage(chatID, responseText)
			bot.Send(msg)
		default:
			log.Fatal("user was in unknown userState\nUser ID: %d, State: %s\n", userID, state.State)
		}
	} else {
		var responseText string
		switch callbackData {
		case "new_contact":
			responseText = "New Contact"
			userStates[userID] = UserState{State: "awaiting_new_contact_name"}
			// msg := tgbotapi.NewMessage(chatID, )
			responseText = "Please type the UserName you want to add as a Contact!"

		case "register_self":
			userName := update.CallbackQuery.From.UserName
			registerUser(userID, userName)
		case "deregister_self":
			deregisterUser(userID)
		case "split_with_contacts":
			msg := tgbotapi.NewMessage(chatID, "Which Contacts do you want to split with?")
			msg.ReplyMarkup = buildSplitContactKeyboard(userID)
			bot.Send(msg)
		case "simple_split":
			userStates[userID] = UserState{State: "waiting_for_split_by_amount"}
			msg := tgbotapi.NewMessage(chatID, "By how many people do you want to split the Bill?")
			msg.ReplyMarkup = split_by_amt_keyboard
			bot.Send(msg)
		case "new_Split":
			currentSplit[userID] = Split{
				author:     userID,
				authorName: update.CallbackQuery.From.UserName,
				divisor:    -1,
				amt:        -1.0,
				splitWith:  []int64{},
			}
			msg := tgbotapi.NewMessage(chatID, "Click a button:")
			msg.ReplyMarkup = new_split_keyboard
			bot.Send(msg)
		default:
			responseText = "Unknown option."
		}
		msg := tgbotapi.NewMessage(chatID, responseText)
		bot.Send(msg)
	}

	editMsg := tgbotapi.NewEditMessageText(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		"You selected: "+callbackData,
	)
	bot.Send(editMsg)
}

func messages(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.Chat.ID
	text := update.Message.Text

	if state, ok := userStates[userID]; ok {
		switch state.State {
		case "awaiting_new_contact_name":
			username := text
			var responseText string
			if checkUserNameExists(username) {
			} else {
				responseText = "Username doesnt exist!"
			}
			msg := tgbotapi.NewMessage(userID, responseText)
			bot.Send(msg)

			delete(userStates, userID)
		case "awaiting_username":
			// log.Printf("Received username: %s", update.Message.Text)
			username := text
			msg := tgbotapi.NewMessage(userID, "Registration complete! Your username is: "+username)
			bot.Send(msg)

			delete(userStates, userID)
		case "awaiting_amount_to_split":
			split(bot, update)
		default:
			// log.Fatal("user was in an unknown userState: %s")
			log.Fatal("user was in unknown userState\nUser ID: %d, State: %s\n", userID, state.State)
		}

	} else {
		log.Printf("Received message: %s", update.Message.Text)
		switch update.Message.Text {
		// case "/help":
		//     msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Possibl")
		//     bot.Send(msg)
		case "/new":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Click a button:")
			msg.ReplyMarkup = new_request_keyboard
			bot.Send(msg)
		case "/register":
			userStates[userID] = UserState{State: "awaiting_username"}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please provide a username:")
			bot.Send(msg)
		default:
			log.Printf("unknown Command\ncommand: %s", update.Message.Text)
		}
	}
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
		log.Printf(botToken)
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

	// Log bot info
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
		printUserStates()
		if update.Message != nil { // Handle regular messages
			messages(bot, update)
		}
		if update.CallbackQuery != nil { // Handle inline button clicks
			callBackQueries(bot, update)
		}
	}
}
