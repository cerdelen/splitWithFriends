package messages

import (
	"log"

	// "github.com/cerdelen/splitWithFriends/globals"
	// "github.com/cerdelen/splitWithFriends/split"
	"github.com/cerdelen/splitWithFriends/user"
	"github.com/cerdelen/splitWithFriends/user/userstates"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.Chat.ID
	// text := update.Message.Text

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
        case userstates.RequestFromSingleContact:
        case userstates.Start:
        default:
            log.Panicf("unexpected userstates.UserState to receive a Message\nState: %s", user.Users[userID].State)
	}
}
