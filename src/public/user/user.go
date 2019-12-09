package user

import (
	"datastructure"
	"fmt"
	"github.com/spf13/viper"
	"log"
)

func User(a datastructure.Request) (err error) {
	var (
		c datastructure.Config
	)
	if err = viper.Unmarshal(&c); err != nil {
		log.Fatalf("Unable To Decode Into Config Struct, %v", err)
		return
	}
	for _, c := range c.Userlist {
		if c.Name == a.Info.RequestMan && c.PhoneNumber == a.Info.PhoneNumber {
			log.Println("GET THE MAN", c.Name)
			return
		}
	}
	log.Printf("{requestMan: %v} AND {phoneNumber: %v} NEED TO BE RIGHT！", a.Info.RequestMan, a.Info.PhoneNumber)
	err = fmt.Errorf("{requestMan: %v} AND {phoneNumber: %v} NEED TO BE RIGHT！", a.Info.RequestMan, a.Info.PhoneNumber)
	return
}
