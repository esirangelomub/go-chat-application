package configs

import (
	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
)

type Conf struct {
	DBDriver               string `mapstructure:"DB_DRIVER"`
	DBHost                 string `mapstructure:"DB_HOST"`
	DBPort                 string `mapstructure:"DB_PORT"`
	DBUser                 string `mapstructure:"DB_USER"`
	DBPassword             string `mapstructure:"DB_PASSWORD"`
	DBName                 string `mapstructure:"DB_NAME"`
	WebServerPort          string `mapstructure:"WEB_SERVER_PORT"`
	JwtSecret              string `mapstructure:"JWT_SECRET"`
	JwtExpiresIn           int    `mapstructure:"JWT_EXPIRES_IN"`
	RabbitMQHost           string `mapstructure:"RABBITMQ_HOST"`
	RabbitMQPort           string `mapstructure:"RABBITMQ_PORT"`
	RabbitMQUser           string `mapstructure:"RABBITMQ_USER"`
	RabbitMQPassword       string `mapstructure:"RABBITMQ_PASSWORD"`
	RabbitMQQueueBot       string `mapstructure:"RABBITMQ_QUEUE_BOT"`
	RabbitMQQueueWebSocket string `mapstructure:"RABBITMQ_QUEUE_WEBSOCKET"`
	RabbitMQExchange       string `mapstructure:"RABBITMQ_EXCHANGE"`
	SecretKey              string
	TokenAuth              *jwtauth.JWTAuth
	BotURL                 string `mapstructure:"BOT_URL"`
}

func LoadConfig(path string) (*Conf, error) {
	var cfg *Conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	cfg.TokenAuth = jwtauth.New("HS256", []byte(cfg.JwtSecret), nil)

	return cfg, nil
}
