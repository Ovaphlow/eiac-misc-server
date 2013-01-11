package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type TraceHandler struct {
	h http.Handler
}

type Person struct {
	UserName string
	Gender   string
	FileName string
	DestDir  string
}

type File struct {
	FileName string
}

func getFileName(fileName string) string {
	path := strings.Split(fileName, "\\")
	index := len(path)
	names := strings.Split(path[index-1], ".")
	return names[0]
}

func (r TraceHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	println("get", req.URL.Path, " from ", req.RemoteAddr)
	r.h.ServeHTTP(w, req)
}

func testPage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("download.html")
	p := Person{UserName: "Ovaphlow", FileName: "cost.pdf"}
	t.Execute(w, p)
}

func oldTestPage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       //解析参数，默认不解析的
	fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello 用户!") //这个写入到w的是输出到客户端的
}

func logonPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	r.ParseForm()
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login/login.html")
		t.Execute(w, nil)
	} else {
		//请求的是登陆数据，那么执行登陆的逻辑判断
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
	}
}

func indexPage(w http.ResponseWriter, r *http.Request) {

}

func convertToPDF(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	r.ParseForm()
	if r.Method == "GET" {
		fmt.Println("fileName:", r.Form["file"])
		fmt.Println("destDir:", r.Form["dir"])
		var fileName []string = r.Form["file"]
		var destDir []string = r.Form["dir"]
		cmd := exec.Command("soffice.exe", "--headless", "-convert-to", "pdf", fileName[0], "-outdir", destDir[0])
		//cmd := exec.Command("soffice.exe", "--headless", "-convert-to", "pdf", fileName[0], "-outdir", "d:\\ftp_root\\")
		buf, err := cmd.Output()
		fmt.Printf("%s\n%s", buf, err)
		//fmt.Fprintf(w, "下载文件准备完毕，请关闭窗口")
		name := getFileName(fileName[0])
		//http.Redirect(w, r, "ftp://1123:1123@172.19.8.242/"+name+".pdf", http.StatusFound)
		fmt.Fprintf(w, "<a href=ftp://1123:1123@127.0.0.1/"+name+".pdf>合同文件下载链接</a>")
	}
}

func uploadPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("upload.html")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}

func main() {
	fmt.Println("服务器正在运行，端口为8081，请不要关闭程序。")
	fmt.Println("转换PDF文档的路径为/converttopdf")
	//http.HandleFunc("/", testPage)
	http.HandleFunc("/test", testPage)
	http.HandleFunc("/logon", logonPage)
	http.HandleFunc("/index", indexPage)
	http.HandleFunc("/upload", uploadPage)
	http.HandleFunc("/converttopdf", convertToPDF)
	err := http.ListenAndServe(":8081", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
