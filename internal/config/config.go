package config

import (
	"os"
)

type Configuration struct {
	Port string
}

func NewConfig() Configuration {
	newConf := os.Getenv("PORT")
	if newConf == "" {
		return Configuration{
			Port: ":8080",
		}
	}

	return Configuration{
		Port: newConf,
	}
}
