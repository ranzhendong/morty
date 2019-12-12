package configanalysis

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
)

func NewLoadConfig() (err error, token *viper.Viper) {
	var (
		pwd, absoluteConf, tokenFile string
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
		log.Println("[LoadConfig] Config file changed:", e.Name)
	})
	token.WatchConfig()
	token.OnConfigChange(func(e fsnotify.Event) {
		log.Println("[LoadConfig] Token file changed:", e.Name)
	})

	if pwd, err = os.Getwd(); err != nil {
		os.Exit(1)
		return
	}
	log.Println("[LoadConfig] Morty Is Running, Execute Path", pwd)

	//if cli has two parameters -f  and -t
	absoluteConf, tokenFile = conf()
	if absoluteConf != "" {
		viper.AddConfigPath(absoluteConf)
		log.Println("[LoadConfig] Get Config Path ", absoluteConf)
	}
	if tokenFile != "" {
		token.AddConfigPath(tokenFile)
		log.Println("[LoadConfig] Get Token Path ", tokenFile)
	}

	//Find and read the config and token file
	if err = viper.ReadInConfig(); err != nil {
		log.Printf("[LoadConfig] Fatal Error Config File: %s \n", err)
		err = fmt.Errorf("[LoadConfig] Fatal Error Config File: %s \n", err)
		return
	}
	if err = token.ReadInConfig(); err != nil {
		log.Printf("[LoadConfig] Fatal Error Token File: %s \n", err)
		err = fmt.Errorf("[LoadConfig] Fatal Error Token File: %s \n", err)
		return
	}

	return
}
