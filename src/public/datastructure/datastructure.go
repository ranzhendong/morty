package datastructure

//定义数据类型，如果没有获取，默认为空

//RequestBody数据结构
type Request struct {
	Deployment      string        `json:"deployment"`
	NameSpace       string        `json:"namespace"`
	DeploymentApi   string        `json:"deploymentApi"`
	JavaProject     string        `json:"javaProject"`
	Version         string        `json:"version"`
	Image           string        `json:"image"`
	MinReadySeconds int           `json:"minReadySeconds"`
	Replicas        int           `json:"replicas"`
	Paused          string        `json:"paused"`
	SendFormat      string        `json:"sendFormat"`
	RollingUpdate   RollingUpdate `json:"rollingUpdate"`
	Info            Info          `json:"info"`
}

type Info struct {
	RequestMan    string `json:"requestMan"`
	UpdateSummary string `json:"updateSummary"`
	PhoneNumber   string `json:"phoneNumber"`
}

type RollingUpdate struct {
	MaxUnavailable string `json:"maxUnavailable"`
	MaxSurge       string `json:"maxSurge"`
}

//配置文件数据结构
type Config struct {
	Userlist  []Userlist `yaml:"userlist"`
	Kubenetes Kubenetes  `yaml:"kubenetes"`
	DingDing  DingDing   `yaml:"dingding"`
}

type Userlist struct {
	Name        string `yaml:"name"`
	ChineseName string `yaml:"chinesename"`
	PhoneNumber string `yaml:"phonenumber"`
}

type Kubenetes struct {
	Host          string `yaml:"host"`
	TokenFile     string `yaml:"tokenfile"`
	DeploymentApi string `yaml:"deploymentapi"`
}

type Token struct {
	Token string `json:"token"`
}

type DingDing struct {
	Robotsurl string `yaml:"robotsurl"`
}

//钉钉消息提示数据结构
//text文本提醒
type DingText struct {
	Msgtype string `json:"msgtype"`
	Text    Text   `json:"text"`
	At      At     `json:"at"`
}

type Text struct {
	Content string `json:"content"`
}

//markdown文本提醒
type DingMarkDown struct {
	Msgtype  string   `json:"msgtype"`
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
