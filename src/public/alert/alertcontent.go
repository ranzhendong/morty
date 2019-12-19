package alert

import (
	"datastructure"
	"math"
	"strconv"
	"strings"
	"time"
)

var (
	requestMux = make(map[string]func(datastructure.Request) (content string, f [1]string))
	t          time.Duration
)

func Main(requestUrl string, a datastructure.Request, time time.Duration) (content string, f [1]string) {
	if requestUrl == "/endsend" || requestUrl == "/grayendsend" {
		t = time
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
	requestMux["/deployupdate"] = dpUpdate
	requestMux["/graydeployupdate"] = grayDpUpdate
	requestMux["/rollback"] = rollBack
	requestMux["/endsend"] = endSend
	requestMux["/grayendsend"] = grayEndSend
}

func dpUpdate(a datastructure.Request) (content string, f [1]string) {
	var (
		sum, replicas, minReadySeconds int64
		stringSum                      string
	)
	replicas, _ = a.Replicas.Int64()
	minReadySeconds, _ = a.MinReadySeconds.Int64()
	sum = replicas * minReadySeconds

	//keep two decimal place
	stringSum = strconv.FormatFloat(float64(sum)/60, 'f', 3, 64)
	stringSum = stringSum[0 : strings.Index(stringSum, ".")+2]

	subject := "乐湃事件通知\n" +
		"触发即时更新操作\n\n" +
		"预计耗时约：" + stringSum + "min" +
		"(" + strconv.FormatFloat(float64(sum), 'f', -1, 64) + "s)"
	// date into struck
	if a.SendFormat == "text" {
		content = subject +
			"\n" + "**更新工程：" + a.JavaProject + "**" +
			"\n" + "项目版本：" + a.Version +
			"\n" + "镜像版本：" + a.Image +
			"\n" + "更新备注：" + a.Info.UpdateSummary +
			"\n" + "执行人：" + a.Info.RequestMan +
			"\n@" + a.Info.PhoneNumber.String()
	} else {
		content = "# " + subject +
			"\n" + "## 更新工程：" + a.JavaProject +
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
	var (
		rt, sum, durationOfStay, aVersionStepWiseUp, bVersionStepWiseUp, bVersionStepWiseDown float64
		stringSum                                                                             string
	)
	rt, _ = a.Gray.TieredRate.Float64()
	durationOfStay, _ = a.Gray.DurationOfStay.Float64()
	aVersionStepWiseUp, _ = a.Gray.AVersionStepWiseUp.Float64()
	bVersionStepWiseUp, _ = a.Gray.BVersionStepWiseUp.Float64()
	bVersionStepWiseDown, _ = a.Gray.BVersionStepWiseDown.Float64()
	sum = durationOfStay + (aVersionStepWiseUp+bVersionStepWiseUp+bVersionStepWiseDown)*(math.Floor(1/rt))

	//keep two decimal place
	stringSum = strconv.FormatFloat(sum/60, 'f', 3, 64)
	stringSum = stringSum[0 : strings.Index(stringSum, ".")+2]

	if a.SendFormat == "text" {
		subject := "乐湃事件通知\n" +
			"触发混合灰度更新操作\n" +
			"预计耗时约：" + stringSum + "min" +
			"(" + strconv.FormatFloat(sum, 'f', -1, 64) + "s)"
		content = subject + "\n" +
			"\n更新工程：" + a.JavaProject +
			"\n" + "项目版本：" + a.Version +
			"\n" + "镜像版本：" + a.Image +
			"\n" + "更新备注：" + a.Info.UpdateSummary +
			"\n" + "执行人：" + a.Info.RequestMan +
			"\n@" + a.Info.PhoneNumber.String()
	} else {
		subject := "乐湃事件通知\n" +
			"触发混合灰度更新操作\n" +
			"> 预计耗时约：**" + stringSum + "min" +
			"(" + strconv.FormatFloat(sum, 'f', -1, 64) + "s)**"
		content = "# " + subject + "\n" +
			"\n## 更新工程**" + a.JavaProject + "**" +
			"\n" + "1. 工程版本：" + a.Version +
			"\n" + "2. 镜像版本：" + a.Image +
			"\n" + "3. 更新备注：" + a.Info.UpdateSummary +
			"\n" + "4. 执行人：" + a.Info.RequestMan +
			"\n@" + a.Info.PhoneNumber.String()
	}

	// @somebody
	f[0] = a.Info.PhoneNumber.String()
	return
}

func endSend(a datastructure.Request) (content string, f [1]string) {
	var subject = "乐湃事件通知\n" +
		"即时更新已经完成\n" +
		"最终总共耗时：" + t.String()

	// date into struck
	if a.SendFormat == "text" {
		content = subject +
			"\n" + "{ " + a.JavaProject + " } 更新完成" +
			"\n" + "项目版本：" + a.Version +
			"\n" + "镜像版本：" + a.Image +
			"\n" + "更新备注：" + a.Info.UpdateSummary +
			"\n" + "执行人：" + a.Info.RequestMan +
			"\n@" + a.Info.PhoneNumber.String()
	} else {
		content = "# " + subject +
			"\n" + "## **" + a.JavaProject + "**更新完成" +
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

func grayEndSend(a datastructure.Request) (content string, f [1]string) {
	var (
		subject, stringT string
	)

	//keep two decimal place
	stringT = t.String()[0 : strings.Index(t.String(), ".")+3]

	// date into struck
	if a.SendFormat == "text" {
		subject = "乐湃事件通知\n" +
			"混合灰度更新已经完成\n" +
			"最终总共耗时：" + stringT + "s"
		content = subject +
			"\n" + a.JavaProject + " 更新完成" +
			"\n" + "项目版本：" + a.Version +
			"\n" + "镜像版本：" + a.Image +
			"\n" + "更新备注：" + a.Info.UpdateSummary +
			"\n" + "执行人：" + a.Info.RequestMan +
			"\n@" + a.Info.PhoneNumber.String()
	} else {
		subject = "乐湃事件通知\n" +
			"混合灰度更新已经完成\n" +
			"> 最终总共耗时：**" + stringT + "s**"
		content = "# " + subject +
			"\n" + "## **" + a.JavaProject + "**更新完成" +
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

func rollBack(a datastructure.Request) (content string, f [1]string) {
	var (
		subject string
	)

	// date into struck
	if a.SendFormat == "text" {
		subject = "乐湃事件通知\n" +
			"触发回滚操作\n"
		content = subject +
			"\n" + "回滚项目：" + a.JavaProject +
			"\n" + "回滚项目版本：" + a.Version +
			"\n" + "回滚镜像版本：" + a.Image +
			"\n" + "回滚更新备注：" + a.Info.UpdateSummary +
			"\n" + "回滚执行人：" + a.Info.RequestMan +
			"\n@" + a.Info.PhoneNumber.String()
	} else {
		subject = "乐湃事件通知\n" +
			"触发回滚操作\n"
		content = "# " + subject +
			"\n" + "## 回滚项目：**" + a.JavaProject + "**" +
			"\n" + "- 回滚项目版本：" + a.Version +
			"\n" + "- 回滚镜像版本：" + a.Image +
			"\n" + "- 回滚更新备注：" + a.Info.UpdateSummary +
			"\n" + "- 回滚执行人：" + a.Info.RequestMan +
			"\n@" + a.Info.PhoneNumber.String()
	}

	// @somebody
	f[0] = a.Info.PhoneNumber.String()
	return
}
