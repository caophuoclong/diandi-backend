package lib

import (
	"github.com/spf13/viper"
	"log"
)

type Env struct {
	ServerPort string `mapstructure:"SERVER_PORT"`
}

func NewEnv() Env {
	env := Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal("Cannot read configuration")
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	return env
}
