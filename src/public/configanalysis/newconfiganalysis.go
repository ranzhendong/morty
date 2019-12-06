package configanalysis

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
)

func NewLoadConfig() (err error) {
	var (
		pwd, newConfFilePath, tokenFile string
		config, token                   *viper.Viper
	)

	config = viper.New()
	token = viper.New()

	//config init
	config.SetConfigName("config")
	config.SetConfigType("yaml")
	config.AddConfigPath(".")
	config.AddConfigPath("./config/")
	token.SetConfigName("token")
	token.AddConfigPath("./token")
	token.AddConfigPath("./config/")

	//watch the config change
	config.WatchConfig()
	config.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
	})
	token.WatchConfig()
	token.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Token file changed:", e.Name)
	})

	if pwd, err = os.Getwd(); err != nil {
		os.Exit(1)
		return
	}
	log.Println("Morty Is Running, Execute Path", pwd)

	//if cli has two parameters -f  and -t
	newConfFilePath, tokenFile = conf(changePath(pwd) + "/")
	if newConfFilePath != "" {
		config.AddConfigPath(newConfFilePath)
	}
	if tokenFile != "" {
		token.AddConfigPath(tokenFile)
	}
	//Find and read the config and token file
	if err = config.ReadInConfig(); err != nil {
		log.Printf("Fatal Error Config File: %s \n", err)
		return
	}
	if err = token.ReadInConfig(); err != nil {
		log.Printf("Fatal Error Token File: %s \n", err)
		return
	}
	//if err = viper.Unmarshal(&c); err != nil {
	//	log.Fatalf("Unable To Decode Into struct, %v", err)
	//	return
	//}
	//log.Println(c)
	return
}
