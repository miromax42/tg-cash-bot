package cache

import "strconv"

type Token string

const (
	prefixUserSettings       = "user:settings:"
	TokenTestKey       Token = "test:key"
)

func UserSettingsToken(userID int64) Token {
	return Token(prefixUserSettings + strconv.FormatInt(userID, 10))
}
