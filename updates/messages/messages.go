package messages

import (
	"log"
	"strconv"

	// "github.com/cerdelen/splitWithFriends/globals"
	"github.com/cerdelen/splitWithFriends/split"
	"github.com/cerdelen/splitWithFriends/user"
	"github.com/cerdelen/splitWithFriends/user/userstates"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	botWrap "github.com/cerdelen/splitWithFriends/bot"
)

func parseDirectRequestAmt(update tgbotapi.Update) (string, error) {
	callbackData := update.CallbackQuery.Data
	userID := update.CallbackQuery.Message.Chat.ID

	var responseText string = ""
	var err error = nil
	switch callbackData {
	default:
        if split.IsValidAmt(callbackData) {
            var amt float64
            amt, err = strconv.ParseFloat(callbackData, 32)
            responseText, err = split.CurrentDirectRequests[userID].ResolveRequest(amt)
            delete(split.CurrentDirectRequests, userID)
        } else {
        }
	}
	return responseText, err
}

func HandleMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.Chat.ID
	// text := update.Message.Text
	var msg string
	var err error

    switch user.Users[userID].State {
        case userstates.Awaiting_amount_to_split:
            // split.CurrentSplits[userID].HandleSplit(bot, update)
			// split.HandleSplit(bot, update, globals.
        case userstates.Awaiting_for_split_by_amount:
        case userstates.Awaiting_for_split_contacts:
        case userstates.Awaiting_new_contact_name:
            // username := text
            // var responseText string
            // if checkUserNameExists(username) {
            // } else {
            //     responseText = "Username doesnt exist!"
            // }
            // msg := tgbotapi.NewMessage(userID, responseText)
            // bot.Send(msg)
            //
            // delete(userStates, userID)
        case userstates.Configuration:
        case userstates.None:
        case userstates.NewDirectRequest:
        case userstates.Start:
        case userstates.AwaitingAmountDirectRequest:
            msg, err = parseDirectRequestAmt(update)
        default:
            log.Panicf("unexpected userstates.UserState to receive a Message\nState: %s", user.Users[userID].State)
	}
	if err == nil {
        botWrap.EditMessage(
			update.CallbackQuery.Message.Chat.ID,
			update.CallbackQuery.Message.MessageID,
            msg,
        )
	} else {
		log.Printf("Error in HandleCallBackQueries: %s", err.Error())
	}
}
