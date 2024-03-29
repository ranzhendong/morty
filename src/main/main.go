package main

import (
	"configanalysis"
	"dpimageupdate"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"time"
)

//定义map来实现路由转发
var (
	mux   map[string]func(http.ResponseWriter, *http.Request)
	token *viper.Viper
	err   error
)

type myHandler struct{}

//初始化log函数
func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	//config init
	if err, token = configanalysis.NewLoadConfig(); err != nil {
		return
	}
	server := http.Server{
		Addr:        ":8080",
		Handler:     &myHandler{},
		ReadTimeout: 5 * time.Second,
	}
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	route(mux)
	if err = server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// 路由
func route(mux map[string]func(http.ResponseWriter, *http.Request)) {
	//镜像更新
	mux["/deployupdate"] = deploymentUpdate
	mux["/graydeployupdate"] = grayDeploymentUpdate
	mux["/rollback"] = rollBack
}

//路由的转发
func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		//用这个handler实现路由转发，相应的路由调用相应func
		h(w, r)
		return
	}
	_, _ = io.WriteString(w, "[ServeHTTP] URL:"+r.URL.String()+"IS NOT EXIST")
}

//instant update
func deploymentUpdate(w http.ResponseWriter, r *http.Request) {
	if err := dpimageupdate.Update(r, token); err != nil {
		_, _ = io.WriteString(w, fmt.Sprint(err))
		return
	}
	_, _ = io.WriteString(w, "[Main.DeploymentUpdate] Deployment Image Update Complete!")
}

//
func grayDeploymentUpdate(w http.ResponseWriter, r *http.Request) {
	//dpimageupd	ate.HandleMsg()
	if err := dpimageupdate.GrayUpdate(r, token); err != nil {
		_, _ = io.WriteString(w, fmt.Sprint(err))
		return
	}
	_, _ = io.WriteString(w, "[Main.GrayDeploymentUpdate] Deployment Image Update Complete!")
}

func rollBack(w http.ResponseWriter, r *http.Request) {
	//dpimageupd	ate.HandleMsg()
	if err := dpimageupdate.RollBack(r, token); err != nil {
		_, _ = io.WriteString(w, fmt.Sprint(err))
		return
	}
	_, _ = io.WriteString(w, "[Main.RollBack] Deployment Rollback Successful!")
}
