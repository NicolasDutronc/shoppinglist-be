package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("config.yml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println(viper.GetString("app.key.size"))
	os.Setenv("APP_KEY_SIZE", "32")
	fmt.Println(viper.GetString("app.key.size"))

	var c map[string]interface{}
	if err := viper.Unmarshal(&c); err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%v", c)
}
