package configanalysis

import (
	"datastructure"
	"flag"
	"reflect"
	"runtime"
	"strings"
)

//结构体判断空值
func IsEmpty(a datastructure.Config) bool {
	return reflect.DeepEqual(a, datastructure.Config{})
}

//运行环境判断
func changePath(pwd string) string {
	operating := runtime.GOOS
	if operating == "windows" {
		pwd = strings.Replace(pwd, "\\", "/", -1)
		return pwd
	}
	return pwd
}

//支持指定配置文件路径
func conf(confFilePath string) (confpath, tokenpah string) {
	var (
		absoluteconf, tokenfile string
	)
	flag.StringVar(&absoluteconf, "f", confFilePath, "default absolute conf is 'relative path+config.yaml'")
	flag.StringVar(&tokenfile, "t", "", "default absolute tokenfile is 'relative path+tokenfile'")
	flag.Parse()
	return
}
