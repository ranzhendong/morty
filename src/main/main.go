package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

//定义一个map来实现路由转发

var mux map[string]func(http.ResponseWriter, *http.Request)
type myHandler struct{}

func main(){
	server := http.Server{
		Addr: ":8080",
		Handler: &myHandler{},
		ReadTimeout: 5*time.Second,
	}
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/tmp"] = Tmp

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	//实现路由的转发
	fmt.Printf("%T", r.URL)
	a,b := mux[r.URL.String()]
	fmt.Println(a,b)
	fmt.Println(r)
	//k := "sssss"
	//s := []byte(k)
	fmt.Println(r.Method)
	fmt.Println(r.URL)
	fmt.Println(r.Header)
	if h, ok := mux[r.URL.String()];ok{
		//用这个handler实现路由转发，相应的路由调用相应func
		h(w, r)
		return
	}
	io.WriteString(w, "URL:"+r.URL.String())
}

func Tmp(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "version 3")
}

