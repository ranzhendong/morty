package configanalysis

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
)

func NewLoadConfig() (err error, token *viper.Viper) {
	var (
		pwd, newConfFilePath, tokenFile string
		//config, token                   *viper.Viper
	)

	//读取文件初始化
	token = viper.New()

	//config and token init
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config/")
	token.SetConfigName("token")
	token.SetConfigType("json")
	token.AddConfigPath(".")
	token.AddConfigPath("./config/")

	//watch the config change
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
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
		viper.AddConfigPath(newConfFilePath)
	}
	if tokenFile != "" {
		token.AddConfigPath(tokenFile)
	}

	//Find and read the config and token file
	if err = viper.ReadInConfig(); err != nil {
		log.Printf("Fatal Error Config File: %s \n", err)
		return
	}
	if err = token.ReadInConfig(); err != nil {
		log.Printf("Fatal Error Token File: %s \n", err)
		return
	}

	return
}
