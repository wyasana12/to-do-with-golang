package config

import (
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

type Config struct {
	PORT           string
	DB_HOST        string
	DB_PORT        string
	DB_USER        string
	DB_PASSWORD    string
	DB_DATABASE    string
	SMTP_HOST      string
	SMTP_PORT      string
	SMTP_EMAIL     string
	SMTP_PASSWORD  string
	APP_URL        string
	JWT_SECRET_KEY string
}

var ENV Config
var Validate *validator.Validate

func LoadConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal(err)
	}

	if err := viper.Unmarshal(&ENV); err != nil {
		log.Fatal(err)
	}

	log.Println("Load Server Successfull")
}

func init() {
	Validate = validator.New()
	log.Println("Validator Running Success.")
}
