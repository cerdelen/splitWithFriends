package globals

import "github.com/cerdelen/splitWithFriends/user"


var RegisteredUsers = make(map[int64]string)
var SplitByValue = make(map[int64]int)
var RetryCounter = make(map[int64]int)
var Users = make(map[int64]*user.User)

