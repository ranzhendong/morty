package k8sapi

import (
	"bytes"
	"crypto/tls"
	"datastructure"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
)

//var (
//	k8sHost     = "https://172.16.0.60:6443"
//	tokenFile   = "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJhZG1pbi10b2tlbi16Z3pidiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJhZG1pbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImNlMzRlYTc0LWY2YmEtNGY0ZS1hMTY3LTQ4MTVjZDlhZjkyZiIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTphZG1pbiJ9.LVu91LXbUvfCCekM2w8qA02m_vAKyXgTvFR1zkn_tjCO9MeODSVt1sqmbUsaqfdIN4lgpyrjw66fBm-lWMlTSeXNZBmAI9DSR-xioKS23JEJjMzN3VRTcgEu22sGSpxbJ15x1qyy9dqFWei07xqYESSP4OzwhO7Qt1nYTYJy8jBXMh_u_ePNyxxSPtwrOzMGXToRnM28YFcsOnJC9brvesq8X8VSOeqmigLshdnczoLoUVkGpeKmLtI4Xj60czr3Wk59rnX18N44szAhRJZ-bYDwqrGOnHZ4j9FIU3eDc3XShIUbStZxxQAscjrD_MwzeXExGneMujEBOLwcbW5qvA"
//)

///apis/extensions/v1beta1/namespaces/default/deployments/nginx-deployment
func APIServerGet(deploymentName, nameSpace string) (err error, bodyContentByte []byte, deploymentUrl string) {
	var (
		a                  datastructure.Request
		c                  datastructure.Config
		config, token      *viper.Viper
		k8sHost, tokenFile string
	)
	// Unmarshal the config
	if err = config.Unmarshal(&c); err != nil {
		log.Fatalf("Unable To Decode Into struct, %v", err)
		return
	}

	//assignment k8sHost and tokenFile
	if c.Kubenetes.Host == "" {
		log.Printf("Config  Kubenetes.Host Is %v ", c.Kubenetes.Host)
		err = fmt.Errorf("Config  Kubenetes.Host Is %v ", c.Kubenetes.Host)
		return
	}
	if c.Kubenetes.TokenFile == "" {
		fmt.Println(token.AllKeys())
	}

	k8sHost = c.Kubenetes.Host
	tokenFile = c.Kubenetes.TokenFile

	// if set DeploymentApi
	if a.DeploymentApi == "" {
		a.DeploymentApi = c.Kubenetes.DeploymentApi
	}

	//url
	deploymentUrl = a.DeploymentApi + "/namespaces/" + nameSpace + "/deployments/" + deploymentName
	requestUrl := k8sHost + deploymentUrl

	// 忽略证书校验
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// 请求
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
	var (
		c                  datastructure.Config
		k8sHost, tokenFile string
		contentType        = "application/json"
	)
	k8sHost = c.Kubenetes.Host
	tokenFile = c.Kubenetes.TokenFile
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
