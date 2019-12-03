package user

import (
	"datastructure"
	"fmt"
)

func User(a datastructure.Request, c datastructure.Config) (err error) {
	fmt.Println("RequestDataStructure:", a)
	fmt.Println("Config:", c)
	//if a.Info.RequestMan ==
	for _, c := range c.Userlist {
		if c.Name == a.Info.RequestMan && c.PhoneNumber == a.Info.PhoneNumber {
			fmt.Println("GET THE MAN", c.Name)
			return
		}
	}
	err = fmt.Errorf("{requestMan: %v} AND {phoneNumber: %v} NEED TO BE RIGHT！", a.Info.RequestMan, a.Info.PhoneNumber)
	return
}
