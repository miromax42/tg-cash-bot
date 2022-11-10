package telegram

import (
	"context"

	tele "gopkg.in/telebot.v3"
)

type ContextKey string

const (
	SettingsKey       ContextKey = "userSettings"
	RequestContextKey ContextKey = "requestContext"
)

func (k ContextKey) String() string {
	return string(k)
}

func RequestContext(c tele.Context) context.Context {
	return c.Get(RequestContextKey.String()).(context.Context) //nolint:forcetypeassert
}
