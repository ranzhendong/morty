package dpimageupdate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8sapi"
	"net/http"
)

type Request struct {
	Deployment    string `json:"deployment"`
	NameSpace     string `json:"namespace"`
	DeploymentApi string `json:"deploymentapi"`
	Image         string `json:"image"`
}

const DeploymentApi = "/apis/extensions/v1beta1"

func Main(r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("read body err, %v\n", err)
		return
	}
	var a Request
	if err = json.Unmarshal(body, &a); err != nil {
		fmt.Printf("Unmarshal err, %v\n", err)
		return
	}
	fmt.Println("Deployment", a.Deployment)
	fmt.Println("NameSpace", a.NameSpace)
	if a.DeploymentApi == "" {
		a.DeploymentApi = DeploymentApi
	}
	fmt.Println("DeploymentApi", a.DeploymentApi)
	_, bodyContentByte := k8sapi.APIServerGet(a.Deployment, a.NameSpace, a.DeploymentApi)
	imageReplace(bodyContentByte)
}

func imageReplace(bodyContentByte []byte) {
	var deploymentMap map[string]interface{}
	if err := json.Unmarshal(bodyContentByte, &deploymentMap); err != nil {
		fmt.Println("deploymentMap json change err: ", err)
		return
	}
	//fmt.Println(deploymentMap)
	fmt.Println(deploymentMap["spec"].(map[string]interface{}))
}
