package util

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Telegram ConfigTelegram `mapstructure:",squash"`
	Exchange ConfigExchange `mapstructure:",squash"`
	DB       ConfigDB       `mapstructure:",squash"`
	Tracing  ConfigTracing  `mapstructure:",squash"`
	Cache    ConfigCache    `mapstructure:",squash"`
	HTTP     ConfigHTTP     `mapstructure:",squash"`
	GRPC     ConfigGRPC     `mapstructure:",squash"`
	Kafka    ConfigKafka    `mapstructure:",squash"`
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

type ConfigCache struct {
	Redis             ConfigRedis   `mapstructure:",squash"`
	LocalTTL          time.Duration `mapstructure:"CACHE_LOCAL_TTL"`
	LocalObjectsCount int           `mapstructure:"CACHE_LOCAL_OBJECTS_COUNT"`
	RedisTTL          time.Duration `mapstructure:"CACHE_REDIS_TTL"`
}

type ConfigRedis struct {
	SocketAddr string `mapstructure:"REDIS_SOCKET_ADDRESS"`
	Password   string `mapstructure:"REDIS_PASSWORD"`
	DB         int    `mapstructure:"REDIS_DB"`
}

//nolint:nakedret
type ConfigHTTP struct {
	Address     string `mapstructure:"HTTP_SERVER_ADDRESS"`
	MetricsPort int    `mapstructure:"HTTP_METRICS_PORT"`
}

type ConfigGRPC struct {
	Address string `mapstructure:"GRPC_SERVER_ADDRESS"`
}

type ConfigKafka struct {
	Address string `mapstructure:"KAFKA_ADDRESS"`
}

func NewConfig() (cfg *Config, err error) {
	viper.SetDefault("TLG_TOKEN", "")

	viper.SetDefault("EXCHANGE_TOKEN", "")
	viper.SetDefault("EXCHANGE_BASE_CURRENCY", "RUB")

	viper.SetDefault("DB_URL", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	viper.SetDefault("DB_TEST_USER_ID", 0)

	viper.SetDefault("CACHE_LOCAL_OBJECTS_COUNT", 10000)
	viper.SetDefault("CACHE_LOCAL_TTL", time.Minute)
	viper.SetDefault("CACHE_REDIS_TTL", time.Hour)
	viper.SetDefault("REDIS_SOCKET_ADDRESS", "localhost:6379")
	viper.SetDefault("REDIS_DB", 0)

	viper.SetDefault("TRACING_URL", "http://localhost:14268/api/traces")

	viper.SetDefault("HTTP_SERVER_ADDRESS", "0.0.0.0:8080")
	viper.SetDefault("HTTP_METRICS_PORT", 2112)
	viper.SetDefault("GRPC_SERVER_ADDRESS", "0.0.0.0:50051")

	viper.SetDefault("KAFKA_ADDRESS", "localhost:9092")

	viper.AutomaticEnv()

	cfg = &Config{}
	if err = viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("config unmarchal: %w", err)
	}

	return cfg, nil
}
