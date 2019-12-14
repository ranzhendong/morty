package alert

import (
	"datastructure"
	"strconv"
)

var (
	requestMux = make(map[string]func(datastructure.Request) (content string, f [1]string))
)

func Main(requestUrl string, a datastructure.Request) (content string, f [1]string) {
	// content route
	route()

	// get the func
	if h, ok := requestMux[requestUrl]; ok {
		return h(a)
	}

	//if func doesn't exist
	content = "DingAlertContent " + requestUrl + " Do Not defined"
	f[0] = a.Info.PhoneNumber.String()
	return
}

func route() {
	requestMux["/dpupdate"] = dpUpdate
	requestMux["/graydpupdate"] = grayDpUpdate
}

func dpUpdate(a datastructure.Request) (content string, f [1]string) {
	var subject = "        乐湃事件通知\n" +
		"即时更新即将完成.....\n"
	// date into struck
	if a.SendFormat == "text" {
		content = subject +
			"\n" + "{ " + a.JavaProject + " } 滚动更新进行中" +
			"\n" + "项目版本：" + a.Version +
			"\n" + "镜像版本：" + a.Image +
			"\n" + "更新备注：" + a.Info.UpdateSummary +
			"\n" + "执行人：" + a.Info.RequestMan +
			"\n@" + a.Info.PhoneNumber.String()
	} else {
		content = "## " + subject +
			"\n" + "### **" + a.JavaProject + "**项目滚动更新完成 \n" +
			"\n" + "1. 项目版本：" + a.Version +
			"\n" + "2. 镜像版本：" + a.Image +
			"\n" + "3. 更新备注：" + a.Info.UpdateSummary +
			"\n" + "4. 执行人：" + a.Info.RequestMan +
			"\n@" + a.Info.PhoneNumber.String()
	}

	// @somebody
	f[0] = a.Info.PhoneNumber.String()
	return
}

func grayDpUpdate(a datastructure.Request) (content string, f [1]string) {
	subject := "        乐湃事件通知\n" +
		"灰度发布更新将持续大约" + strconv.Itoa(a.PausedSecond+60) + "s.....\n"
	// date into struck
	if a.SendFormat == "text" {
		content = subject +
			"\n" + "{ " + a.JavaProject + " } 滚动更新完成" +
			"\n" + "项目版本：" + a.Version +
			"\n" + "镜像版本：" + a.Image +
			"\n" + "更新备注：" + a.Info.UpdateSummary +
			"\n" + "执行人：" + a.Info.RequestMan +
			"\n@" + a.Info.PhoneNumber.String()
	} else {
		content = "## " + subject +
			"\n" + "### **" + a.JavaProject + "**项目滚动更新完成 \n" +
			"\n" + "1. 项目版本：" + a.Version +
			"\n" + "2. 镜像版本：" + a.Image +
			"\n" + "3. 更新备注：" + a.Info.UpdateSummary +
			"\n" + "4. 执行人：" + a.Info.RequestMan +
			"\n@" + a.Info.PhoneNumber.String()
	}

	// @somebody
	f[0] = a.Info.PhoneNumber.String()
	return
}
