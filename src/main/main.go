package main

import (
	"dpimageupdate"
	"io"
	"log"
	"net/http"
	"time"
)

//定义map来实现路由转发
var mux map[string]func(http.ResponseWriter, *http.Request)

type myHandler struct{}

func main() {
	server := http.Server{
		Addr:        ":8080",
		Handler:     &myHandler{},
		ReadTimeout: 5 * time.Second,
	}
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	route(mux)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

// 路由
func route(mux map[string]func(http.ResponseWriter, *http.Request)) {
	//mux["/tmp"] = Tmp
	//镜像更新
	mux["/dpupdate"] = Dpupdate
}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//实现路由的转发
	if h, ok := mux[r.URL.String()]; ok {
		//用这个handler实现路由转发，相应的路由调用相应func
		h(w, r)
		return
	}
	io.WriteString(w, "URL:"+r.URL.String())
}

//func Tmp(w http.ResponseWriter, r *http.Request) {
//	io.WriteString(w, "version 3")
//}

func Dpupdate(w http.ResponseWriter, r *http.Request) {
	dpimageupdate.Main(r)
	io.WriteString(w, "version 3")
}
