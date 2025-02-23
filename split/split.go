package split

import (
	// "fmt"
	// "log"
	"errors"
	// "log"
	"regexp"

	// "strconv"

	// "github.com/cerdelen/splitWithFriends/globals"
	"github.com/cerdelen/splitWithFriends/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	amt             float64
}

var CurrentSplits = make(map[int64]*Split)
var CurrentDirectRequests = make(map[int64]*DirectRequest)

func DirectRequestsInit(author *user.User, recipient *user.User) *DirectRequest {
    return &DirectRequest{author: author, recipient: recipient}
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

func isValidPositiveNumber(s string) bool {
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

