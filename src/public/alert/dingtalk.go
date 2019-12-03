package alert

import (
	"bytes"
	"crypto/tls"
	"datastructure"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var (
	Charset     = "Charset"
	ContentType = "application/json"
	DingTalkUrl = "https://oapi.dingtalk.com/robot/send?access_token=edaef5c4adce3689970145a797da0a33f1da05b12f3573bc9b8196c2f844f6d6"
)

func Ding(a datastructure.Request) (err error) {
	var (
		b, bodyContentByte []byte
		d                  datastructure.DingTalk
		content, subject   string
		//f                  [1]string
	)
	// 忽略证书校验
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// 数据定义
	subject = "乐湃报警处理"
	content = subject + "\n" + "zhendong go!"
	fmt.Println(content)
	d.Msgtype = "text"
	d.Text.Content = content
	if b, err = json.Marshal(d); err == nil {
		log.Printf("Send %v TO DingTalk", string(b))
	}
	body := new(bytes.Buffer)
	body.ReadFrom(bytes.NewBuffer([]byte(strings.ToLower(string(b)))))

	client := &http.Client{Transport: tr}
	requestGet, _ := http.NewRequest("POST", DingTalkUrl, body)
	requestGet.Header.Add("Charset", Charset)
	requestGet.Header.Add("Content-Type", ContentType)
	resp, err := client.Do(requestGet)
	if err != nil {
		log.Printf("Get Request Failed ERR:[%s]", err.Error())
		err = fmt.Errorf("Get Request Failed ERR:[%s]", err.Error())
		return
	}
	bodyContentByte, err = ioutil.ReadAll(resp.Body)
	StatusCode := resp.StatusCode
	bodyContent := string(bodyContentByte)
	fmt.Println(StatusCode, bodyContent)
	return
}
