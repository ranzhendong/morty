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

	for {
		time.Sleep(time.Duration(3) * time.Second)
		break
	}

	log.Println("my second", s)
	if err, bodyContentByte = replaceResourcePaused(bodyContentByte, true); err != nil {
		return
	}
	if err, _ = k8sapi.APIServerPut(bodyContentByte, deploymentUrl, token); err != nil {
		return
	}
	log.Println("replaceResourcePaused:", string(bodyContentByte))

	for {
		time.Sleep(time.Duration(s) * time.Second)
		go secondLoop(bodyContentByte, deploymentUrl, token)
		break
	}

	//obtain the request content and phone number
	content, f = alert.Main(r.URL.String(), a)
	if err = ding(a, content, f); err != nil {
		return
	}
	return
}

func secondLoop(bodyContentByte []byte, deploymentUrl string, token *viper.Viper) (err error) {
	if err, bodyContentByte = replaceResourcePaused(bodyContentByte, false); err != nil {
		return
	}
	if err, _ = k8sapi.APIServerPut(bodyContentByte, deploymentUrl, token); err != nil {
		return
	}
	log.Println("secondLoopReplaceResourcePaused:", string(bodyContentByte))
	return
}

type Resp struct {
	data  int
	error error
}

func handleMsg() {
	resp := make(chan Resp)
	stop := make(chan struct{})
	go func() {
		t := time.Tick(time.Second)
		index := 0
		for {
			select {
			case <-t:
				resp <- Resp{
					data: index,
				}
				index++
			case <-stop:
				resp <- Resp{
					error: fmt.Errorf("time tick stop error"),
				}
			}
		}
	}()

	for {
		select {
		case val := <-resp:
			if val.error != nil {
				log.Fatal(val.error)
			}
			if val.data == 5 {
				stop <- struct{}{}

			}
			fmt.Println("index", val.data)

		}
	}
}
