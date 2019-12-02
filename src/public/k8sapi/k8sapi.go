package k8sapi

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var k8sHost = "https://172.16.0.60:6443"
var tokenFile = "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJhZG1pbi10b2tlbi16Z3pidiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJhZG1pbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImNlMzRlYTc0LWY2YmEtNGY0ZS1hMTY3LTQ4MTVjZDlhZjkyZiIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTphZG1pbiJ9.LVu91LXbUvfCCekM2w8qA02m_vAKyXgTvFR1zkn_tjCO9MeODSVt1sqmbUsaqfdIN4lgpyrjw66fBm-lWMlTSeXNZBmAI9DSR-xioKS23JEJjMzN3VRTcgEu22sGSpxbJ15x1qyy9dqFWei07xqYESSP4OzwhO7Qt1nYTYJy8jBXMh_u_ePNyxxSPtwrOzMGXToRnM28YFcsOnJC9brvesq8X8VSOeqmigLshdnczoLoUVkGpeKmLtI4Xj60czr3Wk59rnX18N44szAhRJZ-bYDwqrGOnHZ4j9FIU3eDc3XShIUbStZxxQAscjrD_MwzeXExGneMujEBOLwcbW5qvA"
var contentType = "application/json"

///apis/extensions/v1beta1/namespaces/default/deployments/nginx-deployment
func APIServerGet(deploymentName, nameSpace, deploymentApi string) (err error, bodyContentByte []byte, deploymentUrl string) {
	deploymentUrl = deploymentApi + "/namespaces/" + nameSpace + "/deployments/" + deploymentName
	requestUrl := k8sHost + deploymentUrl
	// 忽略证书校验
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	requestGet, _ := http.NewRequest("GET", requestUrl, nil)
	requestGet.Header.Add("Authorization", "Bearer "+tokenFile)

	resp, err := client.Do(requestGet)
	if err != nil {
		log.Printf("Get Request Failed ERR:[%s]", err.Error())
		err = fmt.Errorf("Get Request Failed ERR:[%s]", err.Error())
		return
	}
	defer resp.Body.Close()
	// 读取请求体
	bodyContentByte, err = ioutil.ReadAll(resp.Body)
	StatusCode := resp.StatusCode
	bodyContent := string(bodyContentByte)
	if StatusCode != 200 {
		log.Printf("The StatusCode Is %v Bad Response: %v", StatusCode, bodyContent)
		err = fmt.Errorf("The StatusCode Is %v Bad Response: %v", StatusCode, bodyContent)
		return

	}
	log.Println("GET The Deployment: ", bodyContent)
	return
}

func APIServerPut(newDeploymentByte []byte, deploymentUrl string) (err error, bodyContent string) {
	requestUrl := k8sHost + deploymentUrl
	// 忽略证书校验
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	//body := newDeploymentByte

	body := new(bytes.Buffer)
	body.ReadFrom(bytes.NewBuffer(newDeploymentByte))

	requestGet, _ := http.NewRequest("PUT", requestUrl, body)
	requestGet.Header.Set("Content-Type", contentType)
	requestGet.Header.Add("Authorization", "Bearer "+tokenFile)

	resp, err := client.Do(requestGet)
	if err != nil {
		log.Printf("get request failed, err:[%s]", err.Error())
		err = fmt.Errorf("Get Request Failed ERR:[%s]", err.Error())
		return
	}
	defer resp.Body.Close()

	bodyContentByte, err := ioutil.ReadAll(resp.Body)
	bodyContent = string(bodyContentByte)
	StatusCode := resp.StatusCode
	if StatusCode != 200 {
		log.Printf("The StatusCode Is %v Bad Response: %v", StatusCode, bodyContent)
		err = fmt.Errorf("The StatusCode Is %v Bad Response: %v", StatusCode, bodyContent)
		return
	}
	log.Println("PUT The Deployment: ", bodyContent)
	return
}
