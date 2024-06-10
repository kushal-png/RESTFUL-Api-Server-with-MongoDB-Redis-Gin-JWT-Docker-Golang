package initializers

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBUri    string `mapstructure:"MONGO_LOCAL_URL"`
	RedisUri string `mapstructure:"REDIS_URL"`
	Port     string `mapstructure:"PORT"`

	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge      int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`

	ClientOrigin string `mapstructure:"CLIENT_ORIGIN"`

	SmtpHost  string `mapstructure:"SMTP_HOST"`
	SmtpUser  string `mapstructure:"SMTP_USER"`
	SmtpPass  string `mapstructure:"SMTP_PASS"`
	SmtpPort  int    `mapstructure:"SMTP_PORT"`
	EmailFrom string `mapstructure:"EMAIL_FROM"`
}

func LoadConfig(path string) (Config, error) {
	var config Config
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Error in loading configuation")
		return config, err
	}

	err = viper.Unmarshal(&config)
	return config, err
}
