package alert

import (
	"datastructure"
)

var requestMux map[string]func(datastructure.Request) (content string, f [1]string)

func Main(requestUrl string, a datastructure.Request) (content string, f [1]string) {
	requestMux = make(map[string]func(datastructure.Request) (content string, f [1]string))
	route()
	if h, ok := requestMux[requestUrl]; ok {
		h(a)
		return h(a)
	}
	return
}

func route() {
	requestMux["/dpupdate"] = dpUpdate
}

func dpUpdate(a datastructure.Request) (content string, f [1]string) {
	var (
		subject string
	)
	// 数据定义
	subject = "        乐湃事件通知"
	if a.SendFormat == "text" {
		content = subject +
			"\n" + "{ " + a.JavaProject + " } 滚动更新完成" +
			"\n" + "项目版本：" + a.Version +
			"\n" + "镜像版本：" + a.Image +
			"\n" + "更新备注：" + a.Info.UpdateSummary +
			"\n" + "执行人：" + a.Info.RequestMan +
			"\n@" + a.Info.PhoneNumber
	} else {
		content = "## " + subject +
			"\n" + "### **" + a.JavaProject + "**项目滚动更新完成 \n" +
			"\n" + "1. 项目版本：" + a.Version +
			"\n" + "2. 镜像版本：" + a.Image +
			"\n" + "3. 更新备注：" + a.Info.UpdateSummary +
			"\n" + "4. 执行人：" + a.Info.RequestMan +
			"\n@" + a.Info.PhoneNumber
	}

	f[0] = a.Info.PhoneNumber
	return
}
