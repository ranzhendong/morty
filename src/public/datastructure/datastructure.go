package datastructure

//定义数据类型，如果没有获取，默认为空
type Request struct {
	Deployment    string `json:"deployment"`
	NameSpace     string `json:"namespace"`
	DeploymentApi string `json:"deploymentapi"`
	Image         string `json:"image"`
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
	Userlist []MyList `yaml:"userlist"`
}

//用户列表
type MyList struct {
	Name        string `yaml:"name"`
	ChineseName string `yaml:"chinesename"`
	PhoneNumber string `yaml:"phonenumber"`
}

type DingTalk struct {
	Msgtype string
	Text    Text
	At      At
}

type Text struct {
	Content string
}

type At struct {
	AtMobiles [1]string
	IsAtAll   string
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
