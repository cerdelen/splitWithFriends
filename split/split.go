package split

import (
	"fmt"
	"log"
	"regexp"
	// "strconv"

	// "github.com/cerdelen/splitWithFriends/globals"
	"github.com/cerdelen/splitWithFriends/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Split struct {
	author          int64
	authorName      string
	divisor         int
	amt             float64
	splitWith       []int64
    simpleSplit     bool
}

var CurrentSplits = make(map[int64]Split)

func isValidSplit(split Split) bool {
	if split.divisor > 0 && split.amt > 0 {
		return true
	}
	return false
}

func isValidPositiveNumber(s string) bool {
	re := regexp.MustCompile(`^[0-9]+(\.[0-9]{1,2})?$`)
	return re.MatchString(s)
}

func (s Split) HandleSplit(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.Chat.ID
	if isValidSplit(s) {
		// divisor := max(split.divisor, len(split.splitWith))
        if s.simpleSplit {
			amt := s.amt / float64(s.divisor)
			responseText := fmt.Sprintf("Amount everyone has to Pay is: %.2f", amt)
			msg := tgbotapi.NewMessage(userID, responseText)
			bot.Send(msg)
		} else {
			// userName := getUserName(split.author)

			for _, splitter := range s.splitWith {
				// splitter_name, err := getUserName(splitter)
				if user.IsRegistered(splitter) {
					amt := s.amt / float64(len(s.splitWith))
					responseText := fmt.Sprintf("%s splits a Bill of %.2f by %d People.\nThe Amount you have to pay is %.2f", s.authorName, s.amt, s.divisor, amt)
					msg := tgbotapi.NewMessage(splitter, responseText)
					bot.Send(msg)
				} else {
					log.Printf("Unregistered Splitter (i guess unregistered while Split was still assembled: %d", splitter)
				}
			}
		}

	} else {
		log.Println("Invalid Split Object")
	}
}

