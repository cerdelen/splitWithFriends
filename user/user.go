package user

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"

	// "github.com/cerdelen/splitWithFriends/globals"
	// "github.com/cerdelen/splitWithFriends/user"
	// "github.com/cerdelen/splitWithFriends/split"
	// "github.com/cerdelen/splitWithFriends/user"
	"github.com/cerdelen/splitWithFriends/user/userstates"
)

const MAXRETRIES = 3
const RETRIESUNINIT = -12


type Currency int

const (
    Dollar                              Currency = iota
    Euro                                Currency = iota
)

func (s Currency) String() string {
	return [...]string{"$", "€"} [s]
}

func (s Currency) ExplicitString() string {
	return [...]string{"Dollar ($)", "Euro (€)"} [s]
}

type User struct {
    ID                      int64
    Name                    string
	State                   userstates.UserState
    Contacts                []*User
    Retries                 int
    ContactIndexing         int
    Currency                Currency
}

func (u *User)CheckRetryLeft() int {
    if u.Retries != RETRIESUNINIT {
        u.Retries--
    } else {
        u.Retries = MAXRETRIES
    }
    return u.Retries
}

func (u *User)CountPossibleAddableContacts() int {
    count := 0
    for _, registeredUser := range RegisteredUsers {
        if registeredUser.ID != u.ID {
            if !u.HasContact(registeredUser.ID) {
                count++
            }
        }
    }
    return count
}

func (u *User)GetPossibleContacts() []*User {
    var out []*User
    for _, registeredUser := range RegisteredUsers {
        if registeredUser.ID != u.ID {
            if !u.HasContact(registeredUser.ID) {
                out = append(out, registeredUser)
            }
        }
    }
    return out
}

func (u *User)RemoveRetries(userID int64) {
    u.Retries = RETRIESUNINIT
}

var RegisteredUsers []*User
var ToSplitUsers = make(map[int64]*User)
var Users = make(map[int64]*User)

func (u *User)AddContact(other int64) error {
    if IsRegistered(other) {
        u.Contacts = append(u.Contacts, Users[other])
        sortUserSliceByName(u.Contacts)
        return nil
    }
    return errors.New("User is not Registered to become a Contact!")
}

func (u *User)ChangeCurrency() string {
    u.Currency--
    if u.Currency < 0 {
        u.Currency = Euro
    }
    return u.Currency.ExplicitString()
}

func (u *User)RemoveContact(other int64) {
    for i, c := range u.Contacts {
        if c.ID == other {
            u.Contacts = append(u.Contacts[:i], u.Contacts[i+1:]...)
            break
        }
    }
}

func (u *User) HasContact(other int64) bool {
    for _, c := range u.Contacts {
        if c.ID == other {
            return true
        }
    }
    return false
}

func GetId(userName string) (int64, error) {
	for _, u := range Users {
		if u.Name == userName {
			return u.ID, nil
		}
	}
	return -1, errors.New("No User with that User Name")
}

func NameIsRegistered(userName string) bool {
    if id, err := GetId(userName); err == nil {
        for _, contact := range RegisteredUsers {
            if contact.ID == id {
                return true
            }
        }
    }
    return false
}

func AddIfNewUser(userID int64, name string) {
    if _, ok := Users[userID]; !ok {
        Users[userID] = &User{ID: userID, Name: name, State: userstates.None}
    }
}

func sortUserSliceByName(slice []*User) {
    sort.Slice(slice, func(i, j int) bool {
        return strings.ToLower(slice[i].Name) < strings.ToLower(slice[j].Name)
    })
}

func RegisterToBotMessages(userID int64) error {
    if RegisteredIndex(userID) != -1 {
        return errors.New("userId is already registered")
    }
    RegisteredUsers = append(RegisteredUsers, Users[userID])

    sortUserSliceByName(RegisteredUsers)
    // sort.Slice(RegisteredUsers, func(i, j int) bool {
    //     return strings.ToLower(RegisteredUsers[i].Name) < strings.ToLower(RegisteredUsers[j].Name)
    // })
    return nil
}

func DeregisterToBotMessages(userID int64) error {
    if ind := RegisteredIndex(userID); ind != -1 {
        log.Printf("deregisterUser %d", userID)
        RegisteredUsers = append(RegisteredUsers[:ind], RegisteredUsers[ind+1:]...)
    }
	return nil
}

func GetUserName(userID int64) (string, error) {
	if u, exists := Users[userID]; exists {
		return u.Name, nil
	}
	return "", errors.New("userId is not registered")
}

func IsRegistered(userID int64) bool {
    return RegisteredIndex(userID) != -1
}

func RegisteredIndex(userID int64) int {
    for i, contact := range RegisteredUsers {
        if contact.ID == userID {
            return i
        }
    }
	return -1
}

func PrintAllUserStates() {
	fmt.Println("Current state of userStates:")
	for _, u := range Users {
		fmt.Printf("User: %s (%d), State: %s\n", u.Name, u.ID, u.State)
	}
}
