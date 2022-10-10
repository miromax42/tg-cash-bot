package util

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Telegram ConfigTelegram `mapstructure:",squash"`
	Exchange ConfigExchange `mapstructure:",squash"`
	DB       ConfigDB       `mapstructure:",squash"`
}

type ConfigTelegram struct {
	Token string `mapstructure:"TLG_TOKEN"`
}

type ConfigExchange struct {
	Token        string `mapstructure:"EXCHANGE_TOKEN"`
	BaseCurrency string `mapstructure:"EXCHANGE_BASE_CURRENCY"`
}

type ConfigDB struct {
	URL string `mapstructure:"DB_URL"`
}

func NewConfig() (cfg *Config, err error) {
	viper.SetDefault("TLG_TOKEN", "")

	viper.SetDefault("EXCHANGE_TOKEN", "")
	viper.SetDefault("EXCHANGE_BASE_CURRENCY", "RUB")

	viper.SetDefault("DB_URL", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")

	viper.AutomaticEnv()

	cfg = &Config{}
	if err = viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("config unmarchal: %w", err)
	}

	return
}
