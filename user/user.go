package user

import (
	"errors"
	"fmt"
	"log"

	"github.com/cerdelen/splitWithFriends/globals"
	"github.com/cerdelen/splitWithFriends/user/userstates"
)


type User struct {
    ID int64
    Name string
	State userstates.UserState
}

var Users = make(map[int64]*User)

func AddIfNewUser(userID int64, name string) {
    if _, ok := Users[userID]; !ok {
        Users[userID] = &User{ID: userID, Name: name, State: userstates.None}
    }
}

func (u *User) Transition(newState userstates.UserState) {
    switch u.State {
    }
}

func RegisterToBotMessages(userID int64) error {
	log.Printf("registerUser %d", userID)
	if _, exists := globals.RegisteredUsers[userID]; exists {
		return errors.New("userId is already registered")
	}

    globals.RegisteredUsers[userID] = struct{}{}

	return nil
}

func DeregisterToBotMessages(userID int64) error {
	log.Printf("deregisterUser %d", userID)
	delete(globals.RegisteredUsers, userID)
	return nil
}

func GetUserName(userID int64) (string, error) {
	if u, exists := Users[userID]; exists {
		return u.Name, nil
	}
	return "", errors.New("userId is not registered")
}

func IsRegistered(userID int64) bool {
	_, exists := globals.RegisteredUsers[userID]
	return exists
}

func PrintAllUserStates() {
	fmt.Println("Current state of userStates:")
	for _, u := range Users {
		fmt.Printf("User: %s (%d), State: %s\n", u.Name, u.ID, u.State)
	}
}
