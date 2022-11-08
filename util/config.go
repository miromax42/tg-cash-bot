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
	HTTP     ConfigHTTP     `mapstructure:",squash"`
	GRPC     ConfigGRPC     `mapstructure:",squash"`
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

type ConfigHTTP struct {
	Address     string `mapstructure:"HTTP_SERVER_ADDRESS"`
	MetricsPort int    `mapstructure:"HTTP_METRICS_PORT"`
}

type ConfigGRPC struct {
	Address string `mapstructure:"GRPC_SERVER_ADDRESS"`
}

func NewConfig() (cfg *Config, err error) {
	viper.SetDefault("TLG_TOKEN", "")

	viper.SetDefault("EXCHANGE_TOKEN", "")
	viper.SetDefault("EXCHANGE_BASE_CURRENCY", "RUB")

	viper.SetDefault("DB_URL", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	viper.SetDefault("DB_TEST_USER_ID", 0)

	viper.SetDefault("TRACING_URL", "http://localhost:14268/api/traces")

	viper.SetDefault("HTTP_SERVER_ADDRESS", "0.0.0.0:8080")
	viper.SetDefault("HTTP_METRICS_PORT", 2112)
	viper.SetDefault("GRPC_SERVER_ADDRESS", "0.0.0.0:50051")

	viper.AutomaticEnv()

	cfg = &Config{}
	if err = viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("config unmarchal: %w", err)
	}

	return
}
