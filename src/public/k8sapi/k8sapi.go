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

func apiServer(a datastructure.Request, token *viper.Viper, bodyContentByte []byte, method string) (err error, returnBodyContentByte []byte) {
	var (
		myRequest                         *http.Request
		myResponse                        *http.Response
		c                                 datastructure.Config
		t                                 datastructure.Token
		k8sHost, tokenFile, deploymentUrl string
		contentType                       = "application/json"
	)

	// Unmarshal the config and token
	if err = viper.Unmarshal(&c); err != nil {
		log.Printf("[APIServer%v] Unable To Decode Into Config Struct, %v", method, err)
		err = fmt.Errorf("[APIServer%v] Unable To Decode Into Config Struct, %v", method, err)
		return
	}
	if err = token.Unmarshal(&t); err != nil {
		log.Printf("[APIServer%v] Unable To Decode Into Token Struct, %v", method, err)
		err = fmt.Errorf("[APIServer%v] Unable To Decode Into Token Struct, %v", method, err)
		return
	}

	//assignment k8sHost and tokenFile
	if c.Kubenetes.Host == "" {
		log.Printf("[APIServer%v] Config  Kubenetes.Host Is %v ", method, c.Kubenetes.Host)
		err = fmt.Errorf("[APIServer%v] Config  Kubenetes.Host Is %v ", method, c.Kubenetes.Host)
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

	// GET
	if method == "GET" {
		myRequest, _ = http.NewRequest(method, requestUrl, nil)
		myRequest.Header.Add("Authorization", "Bearer "+tokenFile)
	} else if method == "PUT" {
		body := new(bytes.Buffer)
		if _, err = body.ReadFrom(bytes.NewBuffer(bodyContentByte)); err != nil {
			return
		}
		myRequest, _ = http.NewRequest(method, requestUrl, body)
		myRequest.Header.Set("Content-Type", contentType)
		myRequest.Header.Add("Authorization", "Bearer "+tokenFile)
	}

	//if response exist
	if myResponse, err = client.Do(myRequest); err != nil {
		log.Printf("[APIServer%v] Get Request Failed ERR:[%s]", method, err.Error())
		err = fmt.Errorf("[APIServer%v] Get Request Failed ERR:[%s]", method, err.Error())
		return
	}

	//close the request
	defer myResponse.Body.Close()

	// 读取请求体
	bodyContentByte, err = ioutil.ReadAll(myResponse.Body)
	returnBodyContentByte = bodyContentByte
	if myResponse.StatusCode != 200 {
		log.Printf("[APIServer%v] The StatusCode Is %v Bad Response: %v", method, myResponse.StatusCode, string(bodyContentByte))
		err = fmt.Errorf("[APIServer%v] The StatusCode Is %v Bad Response: %v", method, myResponse.StatusCode, string(bodyContentByte))
		return
	}
	log.Printf("[APIServer%v] Return The Deployment: %v", method, string(bodyContentByte))
	return
}

//GET Resource
func APIServerGet(a datastructure.Request, token *viper.Viper) (err error, bodyContentByte []byte) {
	// parameter bodyContentByte is nil
	if err, bodyContentByte = apiServer(a, token, bodyContentByte, "GET"); err != nil {
		return
	}
	return
}

//PUT Resource
func APIServerPut(a datastructure.Request, DeploymentByte []byte, token *viper.Viper) (err error) {
	if err, _ = apiServer(a, token, DeploymentByte, "PUT"); err != nil {
		return
	}
	return
}
