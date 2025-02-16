package keyboards

import (
	"errors"
	"log"
	"strconv"

	"github.com/cerdelen/splitWithFriends/globals"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var StartKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Configuration", "configuration"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("New Split", "new_split"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Request directly from a Contact", "direct_request"),
	),
)

var New_split_keyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Split with Contacts", "split_with_contacts"),
		tgbotapi.NewInlineKeyboardButtonData("Simple Split", "simple_split"),
	),
)

var Split_by_amt_keyboard = tgbotapi.NewInlineKeyboardMarkup(
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

func BuildSplitContactKeyboard(userID int64) (tgbotapi.InlineKeyboardMarkup, error) {
	var keyboardRows [][]tgbotapi.InlineKeyboardButton

	for contactID, contactName := range globals.RegisteredUsers {
        if len(keyboardRows) == 5 {
            row := tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonData("Load more", "load_more_contacts"),
            )
            keyboardRows = append(keyboardRows, row)
            break
        }
        if contactID != userID {
            row := tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonData(contactName, strconv.FormatInt(contactID, 10)),
            )

            keyboardRows = append(keyboardRows, row)
        }
	}

    if len(keyboardRows) == 0 {
        return tgbotapi.InlineKeyboardMarkup{}, errors.New("No Contacts")
    }

    row := tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("That was all!", "finished_selecting_contacts"),
    )

    keyboardRows = append(keyboardRows, row)

    log.Printf("This is the Keyboard %+v", keyboardRows)

	return tgbotapi.NewInlineKeyboardMarkup(keyboardRows...), nil
}
