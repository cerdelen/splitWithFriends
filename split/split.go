package split

import (
	"fmt"
	"errors"
	"log"
	"regexp"

	"github.com/cerdelen/splitWithFriends/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    botWrap "github.com/cerdelen/splitWithFriends/bot"
)

type Split struct {
	author          *user.User
    splitWith       []*user.User
	divisor         int
	amt             float64
    // authorName      string
    // simpleSplit     bool
}

type DirectRequest struct {
	author          *user.User
    recipient       *user.User
}

var CurrentSplits = make(map[int64]*Split)
var CurrentDirectRequests = make(map[int64]*DirectRequest)

func (r *DirectRequest)ResolveRequest(amt float64) (string, error) {
    var err error = nil
    var text string = ""
    log.Printf("Looking for %d", r.recipient.ID)
    if user.IsRegistered(r.recipient.ID) {
        responseText := fmt.Sprintf("%s has requested %.2f %s from you!", r.author.Name, amt, r.author.Currency)
        msg := tgbotapi.NewMessage(r.recipient.ID, responseText)
        botWrap.SendMessage(msg)
        text = fmt.Sprintf("You have requested %.2f %s from %s!", amt, r.author.Currency, r.recipient.Name)
    } else {
        log.Printf("Unregistered Splitter (i guess unregistered while Split was still assembled: %s", r.recipient.Name)
        err = errors.New("Recipient is not Registered")
    }
    return text, err
}

func DirectRequestsInit(author *user.User, recipient *user.User) {
    CurrentDirectRequests[author.ID] = &DirectRequest{author: author, recipient: recipient}
}

func (s *Split)isValidSplit() bool {
	if s.divisor > 0 && s.amt > 0 {
		return true
	}
	return false
}

func Init(author *user.User) *Split {
    return &Split{author: author}
}

func IsValidAmt(s string) bool {
	re := regexp.MustCompile(`^[0-9]+(\.[0-9]{1,2})?$`)
	return re.MatchString(s)
}

func (s *Split)AddSplitter(other *user.User) error {
    if s.author.HasContact(other.ID) {
        s.splitWith = append(s.splitWith, other)
        return nil
    }
    return errors.New("User you tried to add as a Splitter is not a Contact!")
}

func (s *Split)HasSplitter(userID *user.User) bool {
    for _, splitter := range s.splitWith {
        if splitter.ID == userID.ID {
            return true
        }
    }
    return false
}

func (s *Split) HandleSplit(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// userID := update.Message.Chat.ID
	// if isValidSplit(s) {
	// 	// divisor := max(split.divisor, len(split.splitWith))
 //        if s.simpleSplit {
	// 		amt := s.amt / float64(s.divisor)
	// 		responseText := fmt.Sprintf("Amount everyone has to Pay is: %.2f", amt)
	// 		msg := tgbotapi.NewMessage(userID, responseText)
	// 		bot.Send(msg)
	// 	} else {
	// 		// userName := getUserName(split.author)
	//
	// 		for _, splitter := range s.splitWith {
	// 			// splitter_name, err := getUserName(splitter)
	// 			if user.IsRegistered(splitter) {
	// 				amt := s.amt / float64(len(s.splitWith))
	// 				responseText := fmt.Sprintf("%s splits a Bill of %.2f by %d People.\nThe Amount you have to pay is %.2f", s.authorName, s.amt, s.divisor, amt)
	// 				msg := tgbotapi.NewMessage(splitter, responseText)
	// 				bot.Send(msg)
	// 			} else {
	// 				log.Printf("Unregistered Splitter (i guess unregistered while Split was still assembled: %d", splitter)
	// 			}
	// 		}
	// 	}
	//
	// } else {
	// 	log.Println("Invalid Split Object")
	// }
}

