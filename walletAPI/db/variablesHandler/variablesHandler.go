package variablesHandler

import (
	"os"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type DbConfig struct {
	Host     			string
	Port     			string
	User     			string
	Password 			string
	DbName   			string
	ApiKey   			string
	BackgroundChecksUrl string
}

func LoadEnvVariables() (DbConfig, error) {
	err := godotenv.Load("./config/global.env")
	if err != nil {
		return DbConfig{}, err
	}
	return DbConfig{
		Host:     os.Getenv("host"),
		Port:     os.Getenv("port"),
		User:     os.Getenv("user"),
		Password: os.Getenv("password"),
		DbName:   os.Getenv("dbname"),
		ApiKey:   os.Getenv("apiKey"),
		BackgroundChecksUrl: os.Getenv("backgroundChecksUrl"),
	}, nil
}
