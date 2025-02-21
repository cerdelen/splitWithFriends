package keyboards

import (
	"errors"
	"log"
	"strconv"

	// "github.com/cerdelen/splitWithFriends/globals"
	"github.com/cerdelen/splitWithFriends/split"
	"github.com/cerdelen/splitWithFriends/user"
	"github.com/cerdelen/splitWithFriends/utils"
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

var NewSplitKeyboard = tgbotapi.NewInlineKeyboardMarkup(
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

var ConfigurationKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Register to receiving Bot Messages", "register_self"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Deregister to receiving Bot Messages", "deregister_self"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Add Contact", "add_contact"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Remove Contact", "remove_contact"),
	),
)

func BuildSplitContactKeyboard(userID int64) (tgbotapi.InlineKeyboardMarkup, error) {
	var keyboardRows [][]tgbotapi.InlineKeyboardButton

	for _, registeredUser := range user.RegisteredUsers {
        if len(keyboardRows) == 5 {
            row := tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonData("Load more", "load_more_contacts"),
            )
            keyboardRows = append(keyboardRows, row)
            break
        }
        if registeredUser.ID != userID {
            userName, err := user.GetUserName(registeredUser.ID)
            if err != nil { log.Fatalf("Couldnt retrieve Username for %d\nMap of Users %+v", registeredUser.ID, user.Users)}

            row := tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonData(userName, strconv.FormatInt(registeredUser.ID, 10)),
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

const ContactKeyboardContactFields= 6

func BuildAddSplitterKeyboard (userID int64) (tgbotapi.InlineKeyboardMarkup, error) {
    contacts := user.Users[userID].Contacts
    var filteredContacts []*user.User
    for _, contact := range contacts {
        if !split.CurrentSplits[userID].HasSplitter(contact) {
            filteredContacts = append(filteredContacts, contact)
        }
    }
    if len(filteredContacts) > ContactKeyboardContactFields {
        user.Users[userID].ContactIndexing = utils.Min(user.Users[userID].ContactIndexing, len(filteredContacts) - ContactKeyboardContactFields)
    }

    return BuildingContactKeyboard(filteredContacts, user.Users[userID].ContactIndexing)
}

func BuildAddingContactKeyboard (userID int64) (tgbotapi.InlineKeyboardMarkup, error) {
    possibleContacts := user.Users[userID].GetPossibleContacts()
    if len(possibleContacts) > ContactKeyboardContactFields {
        user.Users[userID].ContactIndexing = utils.Min(user.Users[userID].ContactIndexing, len(possibleContacts) - ContactKeyboardContactFields)
    }

    return BuildingContactKeyboard(possibleContacts, user.Users[userID].ContactIndexing)
}

func BuildRemovingContactKeyboard (userID int64) (tgbotapi.InlineKeyboardMarkup, error) {
    contacts := user.Users[userID].Contacts
    if len(contacts) > ContactKeyboardContactFields {
        user.Users[userID].ContactIndexing = utils.Min(user.Users[userID].ContactIndexing, len(contacts) - ContactKeyboardContactFields)
    }

    return BuildingContactKeyboard(contacts, user.Users[userID].ContactIndexing)
}

func BuildingContactKeyboard(contacts []*user.User, skip int) (tgbotapi.InlineKeyboardMarkup, error) {
	var keyboardRows [][]tgbotapi.InlineKeyboardButton
    log.Printf("BuildingCotnactKeyboard: %d", len(contacts))

    if skip > 0 {
        row := tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("Load previous Contacts", "load_prev_contacts"),
        )
        keyboardRows = append(keyboardRows, row)
    }
    for i, contact := range contacts {
        if i < skip {
            continue
        }
        if len(keyboardRows) == ContactKeyboardContactFields && (len(contacts) - skip) > ContactKeyboardContactFields {
            row := tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonData("Load more", "load_more_contacts"),
            )
            keyboardRows = append(keyboardRows, row)
            break
        }
        row := tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData(contact.Name, strconv.FormatInt(contact.ID, 10)),
            )
        keyboardRows = append(keyboardRows, row)
    }

    if len(keyboardRows) == 0 {
        row := tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("No Contacts!", "finished_selecting_contacts"),
        )
        keyboardRows = append(keyboardRows, row)
    } else {
        row := tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("That was all!", "finished_selecting_contacts"),
        )
        keyboardRows = append(keyboardRows, row)
    }

	return tgbotapi.NewInlineKeyboardMarkup(keyboardRows...), nil
}

