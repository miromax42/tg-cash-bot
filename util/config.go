package util

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Telegram ConfigTelegram `mapstructure:",squash"`
}
type ConfigTelegram struct {
	TelegramToken string `mapstructure:"TLG_TOKEN"`
}

func NewConfig() (cfg *Config, err error) {
	viper.SetDefault("TLG_TOKEN", "")

	viper.AutomaticEnv()

	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("config unmarchal: %w", err)
	}

	return
}
