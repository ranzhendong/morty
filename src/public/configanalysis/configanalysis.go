package configanalysis

import (
	"datastructure"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
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
func conf(confFilePath string) string {
	confFilePath = confFilePath + "config.yaml"
	var absoluteconf string
	var f *flag.Flag
	flag.StringVar(&absoluteconf, "f", confFilePath, "default absolute conf is 'relative path+config.yaml'")
	flag.Parse()
	f = flag.Lookup("f")
	if f.DefValue != f.Value.String() {
		return f.Value.String()
	}
	return ""
}

//主函数读取配置
func LoadConfig() (err error, c datastructure.Config) {
	var (
		yamlContent            []uint8
		confFilePath, confName string
	)
	confName = "config.yaml"
	pwd, err := os.Getwd()
	log.Println("Script Execute Path", pwd)
	if err != nil {
		panic(err)
		os.Exit(1)
	}
	executePath := changePath(pwd)
	confFilePath = executePath + "/"
	NewConfFilePath := conf(confFilePath)
	if NewConfFilePath != "" {
		yamlContent, err = ioutil.ReadFile(NewConfFilePath)
		goto AbsoluteConf
	}
	//第一次尝试读取配置
	yamlContent, err = ioutil.ReadFile(confFilePath + confName)
	if err != nil {
		log.Println(err)
	} else {
		goto AbsoluteConf
	}
AbsoluteConf:
	//判断是否读到文件
	if yamlContent == nil {
		panicInfo := "\nCan't Not Get The file Named 'config.yaml' From The Path \n1." +
			executePath + "/\n2." + executePath + "/conf/\nPlease Check it !\n"
		err = fmt.Errorf(panicInfo)
		log.Println(panicInfo)
		return
	}
	//判断是否可以正常解析
	if err = yaml.Unmarshal(yamlContent, &c); err != nil {
		log.Fatalf("Yaml Unmarshal ERROR: %v", err)
		return
	}
	//判断conf文件内容是否为空
	if IsEmpty(c) == true {
		err = fmt.Errorf("Config Do Not HAVE CONTENT! ")
		log.Println("Config Do Not HAVE CONTENT! ")
		return
	}
	log.Println("THE CONFIG: ", c)
	return
}
