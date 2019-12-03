package dpimageupdate

import (
	"alert"
	"datastructure"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"user"
)

const DeploymentApi = "/apis/extensions/v1beta1"

func Main(r *http.Request, c datastructure.Config) (err error) {
	//var (
	//	a                                        datastructure.Request
	//	bodyContentByte, newDeploymentByte, body []byte
	//	deploymentUrl                            string
	//)
	var (
		a    datastructure.Request
		body []byte
	)
	// if the body exist
	if body, err = ioutil.ReadAll(r.Body); err != nil {
		log.Printf("Read Body ERR: %v\n", err)
		return
	}
	//log.Println(string(body))
	// if the body can be turn to json
	if err = json.Unmarshal(body, &a); err != nil {
		log.Printf("Unmarshal Body ERR: %v", err)
		return
	}
	//judge the user if exist
	if err = user.User(a, c); err != nil {
		return
	}
	// if DeploymentApi is not specified
	if a.DeploymentApi == "" {
		a.DeploymentApi = DeploymentApi
	}
	// log the parameter
	if parameter, err := json.Marshal(a); err == nil {
		log.Println(string(parameter))
	}
	//// get deployment info from apiserver
	//if err, bodyContentByte, deploymentUrl = k8sapi.APIServerGet(a.Deployment, a.NameSpace, a.DeploymentApi); err != nil {
	//	return
	//}
	//// replace the image from old to new
	//if err, newDeploymentByte = imageReplace(a, bodyContentByte); err != nil {
	//	return
	//}
	//// put the new deployment info to apiserver
	//if err, _ = k8sapi.APIServerPut(newDeploymentByte, deploymentUrl); err != nil {
	//	return
	//}
	alert.Ding()
	return
}

func imageReplace(a datastructure.Request, bodyContentByte []byte) (err error, newDeploymentByte []byte) {
	var (
		//a             requestdatastructure.RequestDataStructure
		deploymentMap map[string]interface{}
	)
	if err = json.Unmarshal(bodyContentByte, &deploymentMap); err != nil {
		log.Println("Json TO DeploymentMap Json Change ERR: ", err)
		return
	}
	//Containers := deploymentMap["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"]
	//getImage := Containers.([]interface{})[0].(map[string]interface{})["image"].(string)
	//fmt.Println(getImage)
	//fmt.Println("imageReplace:", r.Image)
	deploymentMap["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]interface{})[0].(map[string]interface{})["image"] = a.Image

	if newDeploymentByte, err = json.Marshal(deploymentMap); err != nil {
		log.Println("newDeploymentByte TO Json Change ERR: ", err)
		return
	}
	return
}
