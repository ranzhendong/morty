package dpimageupdate

import (
	"alert"
	"datastructure"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"k8sapi"
	"log"
	"net/http"
	"user"
)

func initCheck(rBody io.Reader) (err error, a datastructure.Request) {
	var (
		body []byte
	)
	// if the body exist
	if body, err = ioutil.ReadAll(rBody); err != nil {
		log.Printf("[DpImageUpdate] Read Body ERR: %v\n", err)
		return
	}
	// if the body can be turn to json
	if err = json.Unmarshal(body, &a); err != nil {
		log.Printf("[DpImageUpdate] Unmarshal Body ERR: %v", err)
		return
	}

	//judge the user if exist
	if err = user.User(a); err != nil {
		return
	}

	// log the parameter
	if parameter, err := json.Marshal(a); err == nil {
		log.Printf("[DpImageUpdate] The Request Body: %v", string(parameter))
	}
	return
}

func replace(a datastructure.Request, bodyContentByte []byte) (err error, newBodyContentByte []byte) {
	// eliminate the Status from deployment
	if err, newBodyContentByte = eliminateStatus(a, bodyContentByte); err != nil {
		return
	}

	//replace resource from deployment, include image, replicas
	if err, newBodyContentByte = replaceResource(a, newBodyContentByte); err != nil {
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
		if err, bodyContentByte = eliminateStatus(a, bodyContentByte); err != nil {
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

	//replace the resource
	//the anonymous func is equivalent to func replace
	//if err, bodyContentByte = replace(a, bodyContentByte); err != nil {
	//	return
	//}

	// put the new deployment info to apiserver
	if err, _ = k8sapi.APIServerPut(bodyContentByte, deploymentUrl, token); err != nil {
		return
	}

	//obtain the request content and phone number
	content, f = alert.Main(r.URL.String(), a)
	//dingding alert
	if err = alert.Ding(content, f, a.SendFormat); err != nil {
		log.Printf("[DpImageUpdate] Dingding ERROR:[%s]", err)
		err = fmt.Errorf("[DpImageUpdate] Dingding ERROR:[%v] %v", err,
			"\n DingAlert Filed, But Request Has Been Done, Do Not Worry !")
		return
	}
	return
}

//func GrayDpUpdate(r *http.Request, token *viper.Viper) (err error) {
//	var (
//		a                                        datastructure.Request
//		f                                        [1]string
//		deploymentUrl, content                   string
//		bodyContentByte, newDeploymentByte, body []byte
//	)
//	// if the body exist
//	if body, err = ioutil.ReadAll(r.Body); err != nil {
//		log.Printf("[GrayDpUpdate] Read Body ERR: %v\n", err)
//		return
//	}
//	//log.Println(string(body))
//	// if the body can be turn to json
//	if err = json.Unmarshal(body, &a); err != nil {
//		log.Printf("[GrayDpUpdate] Unmarshal Body ERR: %v", err)
//		return
//	}
//
//	//judge the user if exist
//	if err = user.User(a); err != nil {
//		return
//	}
//
//	// log the parameter
//	if parameter, err := json.Marshal(a); err == nil {
//		log.Printf("[GrayDpUpdate] The Request Body: %v", string(parameter))
//	}
//
//	// get deployment info from apiserver
//	if err, bodyContentByte, deploymentUrl = k8sapi.APIServerGet(a, token); err != nil {
//		return
//	}
//
//	// eliminate the Status from deployment
//	if err, bodyContentByte = eliminateStatus(a, bodyContentByte); err != nil {
//		return
//	}
//
//	//replace resource from deployment, include image, replicas
//	if err, bodyContentByte = replaceResource(a, bodyContentByte); err != nil {
//		return
//	}
//
//	// put the new deployment info to apiserver
//	if err, _ = k8sapi.APIServerPut(newDeploymentByte, deploymentUrl, token); err != nil {
//		return
//	}
//
//	//obtain the request content and phone number
//	content, f = alert.Main(r.URL.String(), a)
//	//dingding alert
//	if err = alert.Ding(content, f, a.SendFormat); err != nil {
//		log.Printf("[DpImageDate] Dingding ERROR:[%s]", err)
//		err = fmt.Errorf("[DpImageDate] Dingding ERROR:[%v] %v", err,
//			"\n DingAlert Filed, But Request Has Been Done, Do Not Worry !")
//		return
//	}
//	return
//}
