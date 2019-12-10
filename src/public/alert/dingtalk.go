package alert

import (
	"bytes"
	"crypto/tls"
	"datastructure"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	ContentType = "application/json"
)

func Ding(content string, f [1]string, sendFormat string) (err error) {
	var (
		b, bodyContentByte []byte
		subject            string
		c                  datastructure.Config
	)
	// Unmarshal the config and token
	if err = viper.Unmarshal(&c); err != nil {
		log.Fatalf("[APIServerGet] Unable To Decode Into Config Struct, %v", err)
		return
	}

	// 忽略证书校验
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// 数据定义
	subject = "        乐湃事件通知"
	if sendFormat == "text" {
		var d = datastructure.DingText{
			"text",
			datastructure.Text{
				content,
			},
			datastructure.At{
				f,
				"false",
			},
		}
		if b, err = json.Marshal(d); err == nil {
			log.Printf("[DingAlert] Send TO DingTalk %v ", string(b))
		}
	} else {
		var d = datastructure.DingMarkDown{
			"markdown",
			datastructure.MarkDown{
				subject,
				content,
			},
			datastructure.At{
				f,
				"false",
			},
		}
		if b, err = json.Marshal(d); err == nil {
			log.Printf("[DingAlert] Send TO DingTalk %v ", string(b))
		}
	}

	body := new(bytes.Buffer)
	body.ReadFrom(bytes.NewBuffer([]byte(b)))

	client := &http.Client{Transport: tr}
	requestGet, _ := http.NewRequest("POST", c.DingDing.Robotsurl, body)
	requestGet.Header.Add("Content-Type", ContentType)
	resp, err := client.Do(requestGet)
	if err != nil {
		log.Printf("[DingAlert] Post Request Failed ERR:[%s]", err.Error())
		err = fmt.Errorf("[DingAlert] Post Request Failed ERR:[%s]", err.Error())
		return
	}
	bodyContentByte, err = ioutil.ReadAll(resp.Body)
	StatusCode := resp.StatusCode
	bodyContent := string(bodyContentByte)
	if StatusCode != 200 {
		log.Printf("[DingAlert] The StatusCode Is %v Bad Response: %v", StatusCode, bodyContent)
		err = fmt.Errorf("[DingDingAlert] The StatusCode Is %v Bad Response: %v", StatusCode, bodyContent)
		return
	}
	return
}
