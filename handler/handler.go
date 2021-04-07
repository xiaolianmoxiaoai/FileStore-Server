package handler

import (
	"encoding/json"
	"filestore-server/meta"
	"filestore-server/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
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
		fmt.Println(fileMeta.FileSha1)
		//meta.UpdateFileMeta(fileMeta)
		meta.UpdateFileDb(fileMeta)
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}

}

func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w,"Upload finished!")

}

// 获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form["filehash"][0]
	//fMeta := meta.GetFileMeta(filehash)
	fMeta, err := meta.GetFileMetaDb(filehash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(fMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)


}

func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	limitCnt,_ := strconv.Atoi(r.Form.Get("limit"))
	fileMetas := meta.GetLastFileMetas(limitCnt)
	data, err := json.Marshal(fileMetas)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(data)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileSha1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fileSha1)
	f, err := os.Open(fm.Location)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Disposition", "attachment;filename=\""+fm.FileName+"\"")

	w.Write(data)
}

func UpdateFileMeta(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	op := r.Form.Get("op")
	newFileName := r.Form.Get("filename")
	if op != "0"{
		w.WriteHeader(http.StatusForbidden)
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	currFile := meta.GetFileMeta(fsha1)
	currFile.FileName = newFileName
	meta.UpdateFileMeta(currFile)
	w.WriteHeader(http.StatusOK)
	data, err := json.Marshal(currFile)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(data)
}

func RemoveFileMeta(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	fMeta := meta.GetFileMeta(fsha1)
	os.Remove(fMeta.Location)
	meta.DeletaFileMeta(fsha1)
	w.WriteHeader(http.StatusOK)

}