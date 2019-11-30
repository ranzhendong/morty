package k8sapi

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func APIServerGet(url, name, nameSpace, endPointApi, tokenFile string) (err error, bodyContent string) {
	endPointApi = strings.Replace(endPointApi, "myNameSpaces", nameSpace, -1)
	endPointApi = strings.Replace(endPointApi, "myEndPoints", name, -1)
	requestUrl := url + endPointApi
	// 忽略证书校验
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	requestGet, _ := http.NewRequest("GET", requestUrl, nil)
	requestGet.Header.Add("Authorization", "Bearer "+tokenFile)

	resp, err := client.Do(requestGet)
	if err != nil {
		fmt.Printf("get request failed, err:[%s]", err.Error())
		return
	}
	defer resp.Body.Close()
	// 读取请求体
	bodyContentByte, err := ioutil.ReadAll(resp.Body)
	bodyContent = string(bodyContentByte)
	StatusCode := resp.StatusCode
	if StatusCode != 200 {
		bodyContent = ""
		return
	}
	return
}

func APIServerPut(url, name, nameSpace, endPointApi, tokenFile, contentType, yamlConverter string) (err error, bodyContent string) {
	endPointApi = strings.Replace(endPointApi, "myNameSpaces", nameSpace, -1)
	endPointApi = strings.Replace(endPointApi, "myEndPoints", name, -1)
	requestUrl := url + endPointApi
	// 忽略证书校验
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	//两种str2bytes方式，下面的一种是利用缓存实现
	//body := []byte(yamlConverter)
	body := new(bytes.Buffer)
	body.ReadFrom(strings.NewReader(yamlConverter))

	requestGet, _ := http.NewRequest("PUT", requestUrl, body)
	requestGet.Header.Set("Content-Type", contentType)
	requestGet.Header.Add("Authorization", "Bearer "+tokenFile)

	resp, err := client.Do(requestGet)
	if err != nil {
		fmt.Printf("get request failed, err:[%s]", err.Error())
		return
	}
	defer resp.Body.Close()

	bodyContentByte, err := ioutil.ReadAll(resp.Body)
	bodyContent = string(bodyContentByte)
	//StatusCode := resp.StatusCode
	return
}
