package user

import (
	"github.com/cerdelen/splitWithFriends/user/userstates"
)


type User struct {
    ID int64
    Name string
	State userstates.UserState
}

func AddIfNewUser(userID int64, name string, users map[int64]*User) {
    if _, ok := users[userID]; !ok {
        users[userID] = &User{ID: userID, Name: name, State: userstates.None}
    }
}

func (u *User) Transition(newState userstates.UserState) {
    switch u.State {
    }
}
