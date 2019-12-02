package dpimageupdate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8sapi"
	"log"
	"net/http"
)

type Request struct {
	Deployment    string `json:"deployment"`
	NameSpace     string `json:"namespace"`
	DeploymentApi string `json:"deploymentapi"`
	Image         string `json:"image"`
}

const DeploymentApi = "/apis/extensions/v1beta1"

func Main(r *http.Request) (err error) {
	var (
		a                                  Request
		bodyContentByte, newDeploymentByte []byte
		deploymentUrl                      string
	)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Read Body ERR: %v\n", err)
		return
	}
	if err = json.Unmarshal(body, &a); err != nil {
		log.Printf("Unmarshal Body ERR: %v", err)
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
	// get deployment info from apiserver
	if err, bodyContentByte, deploymentUrl = k8sapi.APIServerGet(a.Deployment, a.NameSpace, a.DeploymentApi); err != nil {
		return
	}
	fmt.Println(err)
	// replace the image from old to new
	//err, newDeploymentByte := a.imageReplace(bodyContentByte)
	if err, newDeploymentByte = a.imageReplace(bodyContentByte); err != nil {
		return
	}
	// put the new deployment info to apiserver
	err, _ = k8sapi.APIServerPut(newDeploymentByte, deploymentUrl)
	if err, _ = k8sapi.APIServerPut(newDeploymentByte, deploymentUrl); err != nil {
		return
	}
	return
}

func (r Request) imageReplace(bodyContentByte []byte) (err error, newDeploymentByte []byte) {
	var (
		deploymentMap map[string]interface{}
		//newDeploymentByte []byte
	)
	if err := json.Unmarshal(bodyContentByte, &deploymentMap); err != nil {
		log.Println("Json TO DeploymentMap Json Change ERR: ", err)
		return
	}
	//Containers := deploymentMap["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"]
	//getImage := Containers.([]interface{})[0].(map[string]interface{})["image"].(string)
	//fmt.Println(getImage)
	//fmt.Println("imageReplace:", r.Image)
	deploymentMap["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]interface{})[0].(map[string]interface{})["image"] = r.Image

	if newDeploymentByte, err = json.Marshal(deploymentMap); err != nil {
		log.Println("newDeploymentByte TO Json Change ERR: ", err)
		return
	}
	return
}
