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

func APIServerGet(a datastructure.Request, token *viper.Viper) (err error, bodyContentByte []byte, deploymentUrl string) {
	var (
		c                  datastructure.Config
		t                  datastructure.Token
		k8sHost, tokenFile string
	)
	// Unmarshal the config and token
	if err = viper.Unmarshal(&c); err != nil {
		log.Fatalf("[APIServerGet] Unable To Decode Into Config Struct, %v", err)
		return
	}
	if err = token.Unmarshal(&t); err != nil {
		log.Fatalf("[APIServerGet] Unable To Decode Into Token Struct, %v", err)
		return
	}

	//assignment k8sHost and tokenFile
	if c.Kubenetes.Host == "" {
		log.Printf("[APIServerGet] Config  Kubenetes.Host Is %v ", c.Kubenetes.Host)
		err = fmt.Errorf("[APIServerGet] Config  Kubenetes.Host Is %v ", c.Kubenetes.Host)
		return
	}
	if c.Kubenetes.TokenFile == "" {
		tokenFile = t.Token
	}

	// if set DeploymentApi
	if a.DeploymentApi == "" {
		a.DeploymentApi = c.Kubenetes.DeploymentApi
	}

	//url
	deploymentUrl = a.DeploymentApi + "/namespaces/" + a.NameSpace + "/deployments/" + a.Deployment
	k8sHost = c.Kubenetes.Host
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
		log.Printf("[APIServerGet] Get Request Failed ERR:[%s]", err.Error())
		err = fmt.Errorf("[APIServerGet] Get Request Failed ERR:[%s]", err.Error())
		return
	}
	defer resp.Body.Close()
	// 读取请求体
	bodyContentByte, err = ioutil.ReadAll(resp.Body)
	StatusCode := resp.StatusCode
	bodyContent := string(bodyContentByte)
	if StatusCode != 200 {
		log.Printf("[APIServerGet] The StatusCode Is %v Bad Response: %v", StatusCode, bodyContent)
		err = fmt.Errorf("[APIServerGet] The StatusCode Is %v Bad Response: %v", StatusCode, bodyContent)
		return

	}
	log.Println("[APIServerGet] GET The Deployment: ", bodyContent)
	return
}

func APIServerPut(newDeploymentByte []byte, deploymentUrl string, token *viper.Viper) (err error, bodyContent string) {
	var (
		c                  datastructure.Config
		t                  datastructure.Token
		k8sHost, tokenFile string
		contentType        = "application/json"
	)
	// Unmarshal the config and token
	if err = viper.Unmarshal(&c); err != nil {
		log.Fatalf("[APIServerPut] Unable To Decode Into Config Struct, %v", err)
		return
	}
	if err = token.Unmarshal(&t); err != nil {
		log.Fatalf("[APIServerPut] Unable To Decode Into Token Struct, %v", err)
		return
	}

	//if token exist
	if c.Kubenetes.TokenFile == "" {
		tokenFile = t.Token
	}

	k8sHost = c.Kubenetes.Host
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
		log.Printf("[APIServerPut] Put request failed, err:[%s]", err.Error())
		err = fmt.Errorf("[APIServerPut] Put Request Failed ERR:[%s]", err.Error())
		return
	}
	defer resp.Body.Close()

	bodyContentByte, err := ioutil.ReadAll(resp.Body)
	bodyContent = string(bodyContentByte)
	StatusCode := resp.StatusCode
	if StatusCode != 200 {
		log.Printf("[APIServerPut] The StatusCode Is %v Bad Response: %v", StatusCode, bodyContent)
		err = fmt.Errorf("[APIServerPut] The StatusCode Is %v Bad Response: %v", StatusCode, bodyContent)
		return
	}
	log.Println("[APIServerPut] PUT The Deployment: ", bodyContent)
	return
}
