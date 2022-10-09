package telegram

type ContextKey string

const (
	SettingsKey ContextKey = "userSettings"
)

func (k ContextKey) String() string {
	return string(SettingsKey)
}
