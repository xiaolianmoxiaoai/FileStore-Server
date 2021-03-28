package main

import (
	"filestore-server/handler"
	"fmt"
	"net/http"
)

func main() {
	// 静态资源处理
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil{
		fmt.Println("Failed to start server, err:%s", err.Error())
	}
}