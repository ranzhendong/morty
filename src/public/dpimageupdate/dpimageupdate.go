package dpimageupdate

import (
	"alert"
	"datastructure"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"github.com/syyongx/php2go"
	"io"
	"io/ioutil"
	"k8sapi"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"user"
)

type MyError struct {
	MyError error
}

var errors = make(chan int)
var MyErrorChan = make(chan MyError)

func initCheck(rBody io.Reader) (err error, a datastructure.Request) {
	var (
		body []byte
	)
	// if the body exist
	if body, err = ioutil.ReadAll(rBody); err != nil {
		log.Printf("[InitCheck] Read Body ERR: %v\n", err)
		err = fmt.Errorf("[InitCheck] Read Body ERR: %v\n", err)
		return
	}

	// if the body can be turn to json
	if err = json.Unmarshal(body, &a); err != nil {
		log.Printf("[InitCheck] Unmarshal Body ERR: %v", err)
		err = fmt.Errorf("[InitCheck] Unmarshal Body ERR: %v", err)
		return
	}

	//judge the user if exist
	if err = user.User(a); err != nil {
		return
	}

	// log the parameter
	if parameter, err := json.Marshal(a); err == nil {
		log.Printf("[InitCheck] The Request Body: %v", string(parameter))
	}
	return
}

//func replace(a datastructure.Request, bodyContentByte []byte) (err error, newBodyContentByte []byte) {
//	// eliminate the Status from deployment
//
//	if err, newBodyContentByte = eliminateStatus(a, bodyContentByte); err != nil {
//		return
//	}
//
//	//replace resource from deployment, include image, replicas
//	if err, newBodyContentByte = replaceResource(a, newBodyContentByte, pauesd, funcs); err != nil {
//		return
//	}
//	return
//}

func ding(a datastructure.Request, content string, f [1]string) (err error) {
	//dingding alert
	if err = alert.Ding(content, f, a.SendFormat); err != nil {
		log.Printf("[Ding] Dingding ERROR:[%s]", err)
		err = fmt.Errorf("[Ding] Dingding ERROR:[%v] %v", err,
			"\n DingAlert Filed, But Request Has Been Done, Do Not Worry !")
		return
	}
	return
}

func anonymousReplace(a datastructure.Request, bodyContentByte []byte, f func(datastructure.Request, []byte) (err error)) (err error) {
	return f(a, bodyContentByte)
}

func DpUpdate(r *http.Request, token *viper.Viper) (err error) {
	var (
		a                      datastructure.Request
		f                      [1]string
		bodyContentByte        []byte
		deploymentUrl, content string
	)
	//Check if body is right
	if err, a = initCheck(r.Body); err != nil {
		return
	}

	// get deployment info from apiserver
	if err, bodyContentByte, deploymentUrl = k8sapi.APIServerGet(a, token); err != nil {
		return
	}

	//replace the resource
	//the anonymous func is equivalent to func replace
	if err = anonymousReplace(a, bodyContentByte, func(a datastructure.Request, bytes []byte) (err error) {
		// eliminate the Status from deployment
		if err, bodyContentByte = eliminateStatus(bodyContentByte); err != nil {
			return
		}
		//replace resource from deployment, include image, replicas
		if err, bodyContentByte = replaceResource(a, bodyContentByte); err != nil {
			return
		}
		return
	}); err != nil {
		return
	}

	// put the new deployment info to apiserver
	if err, _ = k8sapi.APIServerPut(bodyContentByte, deploymentUrl, token); err != nil {
		return
	}

	//obtain the request content and phone number
	content, f = alert.Main(r.URL.String(), a)
	if err = ding(a, content, f); err != nil {
		return
	}
	return
}

func GrayDpUpdate(r *http.Request, token *viper.Viper) (err error) {
	var (
		a                      datastructure.Request
		f                      [1]string
		bodyContentByte        []byte
		s                      int
		deploymentUrl, content string
	)
	//Check if body is right
	if err, a = initCheck(r.Body); err != nil {
		return
	}

	// get deployment info from apiserver
	if err, bodyContentByte, deploymentUrl = k8sapi.APIServerGet(a, token); err != nil {
		return
	}

	//gray deployment timeline
	if a.Paused == "" {
		s = 60
	} else if strings.Contains(a.Paused, "min") {
		if s, err = strconv.Atoi(a.Paused[0 : len(a.Paused)-3]); err != nil {
			log.Printf("[GrayDpUpdate] {%v} Is Not Number In %v", a.Paused[0:len(a.Paused)-3], a.Paused)
			err = fmt.Errorf("[GrayDpUpdate] {%v} Is Not Number In %v", a.Paused[0:len(a.Paused)-3], a.Paused)
			return
		}
		s = s * 60
	} else if strings.Contains(a.Paused, "s") {
		if s, err = strconv.Atoi(a.Paused[0 : len(a.Paused)-1]); err != nil {
			log.Printf("[GrayDpUpdate] {%v} Is Not Number In %v", a.Paused[0:len(a.Paused)-1], a.Paused)
			err = fmt.Errorf("[GrayDpUpdate] {%v} Is Not Number In %v", a.Paused[0:len(a.Paused)-1], a.Paused)
			return
		}
	} else if php2go.IsNumeric(a.Paused) {
		s, _ = strconv.Atoi(a.Paused)
		log.Printf("[GrayDpUpdate] {%v} Has Not Unit, So Default Is Second", a.Paused)
	} else {
		log.Printf("[GrayDpUpdate] Paused: %v Is Null, "+
			"So GrayDeployment Paused Default Is 1 Minute. \n"+
			"Notice: GrayDeployment Are published Later More Than 1 Minute.", a.Paused)
		err = fmt.Errorf("[GrayDpUpdate] Paused: %v Is Null, "+
			"So GrayDeployment Paused Default Is 1 Minute. \n"+
			"Notice: GrayDeployment Are published Later More Than 1 Minute.", a.Paused)
	}

	//replace the resource
	//the anonymous func is equivalent to func replace
	if err = anonymousReplace(a, bodyContentByte, func(a datastructure.Request, bytes []byte) (err error) {
		// eliminate the Status from deployment
		if err, bodyContentByte = eliminateStatus(bodyContentByte); err != nil {
			return
		}
		//replace resource from deployment, include image, replicas
		if err, bodyContentByte = replaceResource(a, bodyContentByte); err != nil {
			return
		}
		return
	}); err != nil {
		return
	}

	// put the new deployment info to apiserver
	if err, _ = k8sapi.APIServerPut(bodyContentByte, deploymentUrl, token); err != nil {
		return
	}

	//gray deployment controller Goroutine
	go pauseGoroutine(a, bodyContentByte, deploymentUrl, s, token)

	//  handle the err of pauseGoroutine,if err exist
	go errHandle()

	//obtain the request content and phone number
	content, f = alert.Main(r.URL.String(), a)
	if err = ding(a, content, f); err != nil {
		return
	}
	return
}

func grayLoop(bodyContentByte []byte, deploymentUrl string, mybool bool, token *viper.Viper) (err error) {
	if err, bodyContentByte = replaceResourcePaused(bodyContentByte, mybool); err != nil {
		return
	}
	log.Println("[GrayLoop] SecondLoopReplaceResourcePaused", string(bodyContentByte))
	if err, _ = k8sapi.APIServerPut(bodyContentByte, deploymentUrl, token); err != nil {
		return
	}
	return
}

func errHandle() {
	for {
		select {
		case <-errors:
			err := <-MyErrorChan
			log.Println(err.MyError)
		}
	}
}

func pauseGoroutine(a datastructure.Request, bodyContentByte []byte, deploymentUrl string, s int, token *viper.Viper) {
	var (
		minReadySeconds int64
		err             error
		interval        = make(chan int)
		first           = make(chan int)
		second          = make(chan int)
	)

	go func() {
		minReadySeconds, _ = a.MinReadySeconds.Int64()
		if minReadySeconds > 10 {
			minReadySeconds = 10
		}
		log.Printf("[Paused] CoolingTime Need TO %v Gray Update", s+int(minReadySeconds)*2)
		for {
			time.Sleep(time.Duration(minReadySeconds) * time.Second)
			break
		}
		first <- 1
	}()

	go func() {
		for {
			select {
			case <-interval:
				log.Println("[Paused] Interval")
				for {
					time.Sleep(time.Duration(s) * time.Second)
					break
				}
				second <- 1
			}
		}
	}()

	for {
		select {
		case <-first:
			log.Println("[Paused] To Be Paued")
			if err = grayLoop(bodyContentByte, deploymentUrl, true, token); err != nil {
				err = fmt.Errorf("[Paued] First Pauese ERR %v", err)
				MyErrorChan <- MyError{err}
				errors <- 1
			}
			interval <- 1
		case <-second:
			log.Println("[Paused] No Be Paued")
			if err = grayLoop(bodyContentByte, deploymentUrl, false, token); err != nil {
				err = fmt.Errorf("[Paued] First Pauese ERR %v", err)
				MyErrorChan <- MyError{err}
				errors <- 1
			}
		}
	}
}
