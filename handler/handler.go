package handler

import (
	"filestore-server/meta"
	"filestore-server/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// 处理文件上传
func UploadHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET"{
		//返回上传文件页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil{
			io.WriteString(w, "internel server error")
			return
		}
		io.WriteString(w, string(data))
	}else if r.Method == "POST" {
		//接受文件流及存储到本地目录
		file,head,err := r.FormFile("file")
		if err != nil{
			fmt.Printf("failed to get data ,err:%s\n", err.Error())
			return
		}
		defer file.Close()

		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "E:/GoPath/src/filestore-server/tmp/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 21:14:11"),
		}

		newFile, err := os.Create(fileMeta.Location)
		if err != nil{
			fmt.Printf("file create failed err:%s\n", err.Error())
			return
		}
		defer newFile.Close()

		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("failed to save data, err:%s\n", err.Error())
			return
		}

		newFile.Seek(0,0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		meta.UpdateFileMeta(fileMeta)

		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}

}

func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w,"Upload finished!")

}

