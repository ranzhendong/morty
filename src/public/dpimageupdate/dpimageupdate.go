package dpimageupdate

import (
	"alert"
	"datastructure"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"k8sapi"
	"log"
	"net/http"
	"user"
)

func Main(r *http.Request, token *viper.Viper) (err error) {
	var (
		a                                        datastructure.Request
		bodyContentByte, newDeploymentByte, body []byte
		f                                        [1]string
		deploymentUrl, content                   string
	)
	// if the body exist
	if body, err = ioutil.ReadAll(r.Body); err != nil {
		log.Printf("[DpImageDate] Read Body ERR: %v\n", err)
		return
	}
	//log.Println(string(body))
	// if the body can be turn to json
	if err = json.Unmarshal(body, &a); err != nil {
		log.Printf("[DpImageDate] Unmarshal Body ERR: %v", err)
		return
	}

	//judge the user if exist
	if err = user.User(a); err != nil {
		return
	}

	// log the parameter
	if parameter, err := json.Marshal(a); err == nil {
		log.Printf("[DpImageDate] The Request Body: %v", string(parameter))
	}

	// get deployment info from apiserver
	if err, bodyContentByte, deploymentUrl = k8sapi.APIServerGet(a, token); err != nil {
		return
	}

	// replace the image from old to new
	if err, newDeploymentByte = imageReplace(a, bodyContentByte); err != nil {
		return
	}

	// put the new deployment info to apiserver
	if err, _ = k8sapi.APIServerPut(newDeploymentByte, deploymentUrl, token); err != nil {
		return
	}

	//obtain the request content and phone number
	content, f = alert.Main(r.URL.String(), a)
	//dingding alert
	if err = alert.Ding(content, f, a.SendFormat); err != nil {
		log.Printf("[DpImageDate] Dingding ERROR:[%s]", err)
		err = fmt.Errorf("[DpImageDate] Dingding ERROR:[%v] %v", err,
			"\n DingAlert Filed, But Request Has Been Done, Do Not Worry !")
		return
	}
	return
}

func imageReplace(a datastructure.Request, bodyContentByte []byte) (err error, newDeploymentByte []byte) {
	var (
		deploymentMap map[string]interface{}
	)
	//Unmarshal the body
	if err = json.Unmarshal(bodyContentByte, &deploymentMap); err != nil {
		log.Printf("[DpImageDate] Json TO DeploymentMap Json Change ERR: %v", err)
		return
	}

	//exchange the image from body
	deploymentMap["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]interface{})[0].(map[string]interface{})["image"] = a.Image

	//Marshal the new body
	if newDeploymentByte, err = json.Marshal(deploymentMap); err != nil {
		log.Println("[DpImageDate] NewDeploymentByte TO Json Change ERR: ", err)
		return
	}
	return
}
