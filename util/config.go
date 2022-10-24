package util

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Telegram ConfigTelegram `mapstructure:",squash"`
	Exchange ConfigExchange `mapstructure:",squash"`
	DB       ConfigDB       `mapstructure:",squash"`
	Tracing  ConfigTracing  `mapstructure:",squash"`
}

type ConfigTelegram struct {
	Token string `mapstructure:"TLG_TOKEN"`
}

type ConfigExchange struct {
	Token        string `mapstructure:"EXCHANGE_TOKEN"`
	BaseCurrency string `mapstructure:"EXCHANGE_BASE_CURRENCY"`
}

type ConfigTracing struct {
	URL string `mapstructure:"TRACING_URL"`
}

type ConfigDB struct {
	URL        string `mapstructure:"DB_URL"`
	TestUserID int64  `mapstructure:"DB_TEST_USER_ID"`
}

func NewConfig() (cfg *Config, err error) {
	viper.SetDefault("TLG_TOKEN", "")

	viper.SetDefault("EXCHANGE_TOKEN", "")
	viper.SetDefault("EXCHANGE_BASE_CURRENCY", "RUB")

	viper.SetDefault("DB_URL", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	viper.SetDefault("DB_TEST_USER_ID", 0)

	viper.SetDefault("TRACING_URL", "http://localhost:14268/api/traces")

	viper.AutomaticEnv()

	cfg = &Config{}
	if err = viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("config unmarchal: %w", err)
	}

	return
}
