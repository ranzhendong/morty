package alert

import (
	"datastructure"
	"fmt"
	"strconv"
	"time"
)

var (
	requestMux = make(map[string]func(datastructure.Request) (content string, f [1]string))
	t          time.Duration
)

func Main(requestUrl string, a datastructure.Request, time time.Duration) (content string, f [1]string) {
	if requestUrl == "/goroutinesend" {
		t = time
		fmt.Println("t = time", t, time)
	}
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
	requestMux["/goroutinesend"] = goroutineSend
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
	var paused int64
	paused, _ = a.Gray.DurationOfStay.Int64()
	subject := "        乐湃事件通知\n" +
		"混合灰度发布更新将持续大约" + strconv.Itoa(int(paused)+60) + "s.....\n"
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

func goroutineSend(a datastructure.Request) (content string, f [1]string) {
	var subject = "        乐湃事件通知\n" +
		"混合灰度发布更新已经完成.....\n" +
		"总共耗时：" + t.String()

	fmt.Println("总共耗时:", t)
	// date into struck
	if a.SendFormat == "text" {
		content = subject +
			"\n" + "{ " + a.JavaProject + " } 滚动更新已经完成" +
			"\n" + "项目版本：" + a.Version +
			"\n" + "镜像版本：" + a.Image +
			"\n" + "更新备注：" + a.Info.UpdateSummary +
			"\n" + "执行人：" + a.Info.RequestMan +
			"\n@" + a.Info.PhoneNumber.String()
	} else {
		content = "## " + subject +
			"\n" + "### **" + a.JavaProject + "**滚动更新已经完成 \n" +
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
