package constants

import "os"

type UserKey string

var UserIdKey UserKey = UserKey(os.Getenv("TASKIFY_USER_ID_KEY"))
