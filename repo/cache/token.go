package cache

import (
	"bytes"
	"crypto/md5" //nolint:gosec
	"encoding/gob"
	"strconv"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo"
)

type Token string

const (
	prefixUserSettings       = "user:settings:"
	prefixUserReport         = "user:report:"
	TokenTestKey       Token = "test:key"
)

func UserSettingsToken(userID int64) Token {
	return Token(prefixUserSettings + strconv.FormatInt(userID, 10))
}

func UserReportToken(req repo.ListUserExpenseReq) Token {
	return Token(prefixUserReport + Hash(req))
}

func Hash[V any](s V) string {
	var b bytes.Buffer
	gob.NewEncoder(&b).Encode(s) //nolint:errcheck
	hash := md5.Sum(b.Bytes())   //nolint:gosec

	return string(hash[:])
}
