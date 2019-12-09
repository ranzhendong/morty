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
		log.Fatalf("[User] Unable To Decode Into Config Struct, %v", err)
		return
	}
	for _, c := range c.Userlist {
		if c.Name == a.Info.RequestMan && c.PhoneNumber == a.Info.PhoneNumber {
			log.Printf("[User] {%v} Is Executing", c.Name)
			return
		}
	}
	log.Printf("[User] {requestMan: %v} AND {phoneNumber: %v} NEED TO BE RIGHT！", a.Info.RequestMan, a.Info.PhoneNumber)
	err = fmt.Errorf("[User] {requestMan: %v} AND {phoneNumber: %v} NEED TO BE RIGHT！", a.Info.RequestMan, a.Info.PhoneNumber)
	return
}
