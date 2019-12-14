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
		log.Printf("[User] Unable To Decode Into Config Struct, %v", err)
		err = fmt.Errorf("[User] Unable To Decode Into Config Struct, %v", err)
		return
	}
	for _, c := range c.UserList {
		// a.Info.PhoneNumber's type is json.Number, so i can decide the data type
		if c.Name == a.Info.RequestMan && c.PhoneNumber == a.Info.PhoneNumber.String() {
			log.Printf("[User] {%v} Is Executing", c.Name)
			return
		}
	}
	log.Printf("[User] {requestMan: %v} AND {phoneNumber: %v} NEED TO BE RIGHT！", a.Info.RequestMan, a.Info.PhoneNumber)
	err = fmt.Errorf("[User] {requestMan: %v} AND {phoneNumber: %v} NEED TO BE RIGHT！", a.Info.RequestMan, a.Info.PhoneNumber)
	return
}
