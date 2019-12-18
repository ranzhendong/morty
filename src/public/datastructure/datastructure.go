package datastructure

import "encoding/json"

//定义数据类型，如果没有获取，默认为空

//RequestBody数据结构
type Request struct {
	Name            string        `json:"name"`
	Image           string        `json:"image"`
	Version         string        `json:"version"`
	NameSpace       string        `json:"namespace"`
	SendFormat      string        `json:"sendFormat"`
	Deployment      string        `json:"deployment"`
	JavaProject     string        `json:"javaProject"`
	DeploymentApi   string        `json:"deploymentApi"`
	MinReadySeconds json.Number   `json:"minReadySeconds"`
	Replicas        json.Number   `json:"replicas"`
	Gray            Gray          `json:"gray"`
	Info            Info          `json:"info"`
	RollingUpdate   RollingUpdate `json:"rollingUpdate"`
	DurationOfStay  int64
}

type Gray struct {
	TieredRate           json.Number `json:"tieredRate"`
	DurationOfStay       json.Number `json:"durationOfStay"`
	AVersionStepWiseUp   json.Number `json:"aVersionStepWiseUp"`
	BVersionStepWiseUp   json.Number `json:"bVersionStepWiseUp"`
	BVersionStepWiseDown json.Number `json:"bVersionStepWiseDown"`
}

type Info struct {
	RequestMan    string      `json:"requestMan"`
	PhoneNumber   json.Number `json:"phoneNumber"`
	UpdateSummary string      `json:"updateSummary"`
}

type RollingUpdate struct {
	MaxSurge       string `json:"maxSurge"`
	MaxUnavailable string `json:"maxUnavailable"`
}

type MySpec struct {
	Spec struct {
		Replicas int `json:"replicas"`
	} `json:"spec"`
}

//配置文件数据结构
type Config struct {
	UserList   []UserList `yaml:"userlist"`
	DingDing   DingDing   `yaml:"dingding"`
	Kubernetes Kubernetes `yaml:"kubernetes"`
}

type UserList struct {
	Name        string `yaml:"name"`
	ChineseName string `yaml:"chinesename"`
	PhoneNumber string `yaml:"phonenumber"`
}

type Kubernetes struct {
	Host          string `yaml:"host"`
	TokenFile     string `yaml:"tokenfile"`
	DeploymentApi string `yaml:"deploymentapi"`
}

type Token struct {
	Token string `json:"token"`
}

type DingDing struct {
	RobotsUrl string `yaml:"robotsurl"`
}

//钉钉消息提示数据结构
//text文本提醒
type DingText struct {
	MsgType string `json:"msgtype"`
	Text    Text   `json:"text"`
	At      At     `json:"at"`
}

type Text struct {
	Content string `json:"content"`
}

//markdown文本提醒
type DingMarkDown struct {
	MsgType  string   `json:"msgtype"`
	MarkDown MarkDown `json:"markdown"`
	At       At       `json:"at"`
}

type MarkDown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type At struct {
	AtMobiles [1]string `json:"atMobiles"`
	IsAtAll   string    `json:"isAtAll"`
}
