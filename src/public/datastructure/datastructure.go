package datastructure

//定义数据类型，如果没有获取，默认为空
type Request struct {
	Deployment    string `json:"deployment"`
	NameSpace     string `json:"namespace"`
	DeploymentApi string `json:"deploymentapi"`
	JavaProject   string `json:"javaProject"`
	Version       string `json:"version"`
	Image         string `json:"image"`
	SendFormat    string `json:"sendFormat"`
	Info          Info   `json:"info"`
}

//详细类型
type Info struct {
	RequestMan    string `json:"requestMan"`
	UpdateSummary string `json:"updateSummary"`
	PhoneNumber   string `json:"phoneNumber"`
}

//配置文件
type Config struct {
	Userlist  []Userlist `yaml:"userlist"`
	Kubenetes Kubenetes  `yaml:"kubenetes"`
	DingDing  DingDing   `yaml:"dingding"`
}

//用户列表
type Userlist struct {
	Name        string `yaml:"name"`
	ChineseName string `yaml:"chinesename"`
	PhoneNumber string `yaml:"phonenumber"`
}

//k8s配置
type Kubenetes struct {
	Host          string `yaml:"host"`
	TokenFile     string `yaml:"tokenfile"`
	DeploymentApi string `yaml:"deploymentapi"`
}

// token file
type Token struct {
	Token string `json:"token"`
}

//钉钉配置
type DingDing struct {
	Robotsurl string `yaml:"robotsurl"`
}

//text文本提醒
type DingText struct {
	Msgtype string `json:"msgtype"`
	Text    Text   `json:"text"`
	At      At     `json:"at"`
}

//markdown文本提醒
type DingMarkDown struct {
	Msgtype  string   `json:"msgtype"`
	MarkDown MarkDown `json:"markdown"`
	At       At       `json:"at"`
}

type Text struct {
	Content string `json:"content"`
}

type MarkDown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type At struct {
	AtMobiles [1]string `json:"atMobiles"`
	IsAtAll   string    `json:"isAtAll"`
}

type F interface {
	Sfs()
}

func (r *Request) Sfs() *Request {
	return r
}

func (r *Info) Sfs() *Info {
	return r
}
