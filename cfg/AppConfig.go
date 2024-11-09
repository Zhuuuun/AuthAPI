package cfg

import (
	"log"

	"github.com/spf13/viper"
)

type AppEnvironment string

const (
	PROD AppEnvironment = "PROD"
	QA   AppEnvironment = "QA"
)

type AppConfig struct {
	JwtSecret          string
	GoogleClientID     string
	GoogleClientSecret string
	Port               string
	AppEnvironment     AppEnvironment
}

func LoadCfg() (*AppConfig, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Config file not found:", err)
	}

	config := &AppConfig{
		JwtSecret:          viper.GetString("JWT_SECRET"),
		GoogleClientID:     viper.GetString("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: viper.GetString("GOOGLE_CLIENT_SECRET"),
		Port:               viper.GetString("PORT"),
		AppEnvironment:     loadAppEnvironmentCfg(viper.GetString("APP_ENV")),
	}
	return config, nil
}

func loadAppEnvironmentCfg(appEnv string) AppEnvironment {
	switch appEnv {
	case string(PROD):
		return PROD
	case string(QA):
		return QA
	default:
		log.Fatal("No environment was provided. Set environment to QA")
		return QA
	}
}
