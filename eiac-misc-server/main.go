package main

import (
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	XMLName   xml.Name `xml:"Server"`
	WS_IP     string   `xml:"ws_ip"`
	WS_PORT   string   `xml:"ws_port"`
	WS_OUTDIR string   `xml:"ws_outdir"`
	FS_IP     string   `xml:"fs_ip"`
	FS_PORT   string   `xml:"fs_port"`
	FS_ROOT   string   `xml:"fs_root"`
}
type File struct {
	Token string
	Name  string
	Type  string
}

var fs_ip, fs_port, out_dir string

func getFileName(fileName string) string {
	path := strings.Split(fileName, "/")
	index := len(path)
	names := strings.Split(path[index-1], ".")
	return names[0]
}

func getFileExt(fileName string) string {
	name := strings.Split(fileName, ".")
	index := len(name)
	return name[index-1]
}

func convertToPDF(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	r.ParseForm()
	if r.Method == "GET" {
		//fmt.Println("file name:", r.Form["file"])
		//fmt.Println("dest dir:", r.Form["dir"])
		file := r.Form["file"]
		//dest_dir := r.Form["dir"]
		file_name := strings.Replace(file[0], `/`, `\`, -1)
		fmt.Printf("file name:%s\n", file_name)
		fmt.Printf("out dir:%s\n", out_dir)
		cmd := exec.Command("soffice.exe", "--headless", "-convert-to", "pdf", file_name, "-outdir", out_dir)
		//cmd := exec.Command("soffice.exe", "--headless", "-convert-to", "pdf", file[0], "-outdir", dest_dir[0])
		buf, err := cmd.Output()
		fmt.Printf("%s\n%s", buf, err)
		//fmt.Fprintf(w, "下载文件准备完毕，请关闭窗口")
		name := getFileName(file[0])
		fmt.Println(name)
		fmt.Fprintf(w, "<a href=http://"+fs_ip+":"+fs_port+"/"+name+".pdf>合同文件下载链接</a>（不能正常下载的时候可以鼠标右键选择[目标另存为]）")
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	r.ParseForm()
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))

		fname := r.Form["fname"]
		ftype := r.Form["ftype"]
		s := fmt.Sprintf("%x", h.Sum(nil))
		fmt.Println(s)
		fmt.Println(fname)
		fmt.Println(ftype)
		file := File{Token: s, Name: fname[0], Type: ftype[0]}

		t, _ := template.ParseFiles("upload.html")
		t.Execute(w, file)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		//fmt.Fprintf(w, "%v", handler.Header)
		file_name := r.Form["fname"][0]
		file_type := r.Form["ftype"]

		var file_dir string

		switch {
		case file_type[0] == "report":
			file_dir = "report"
		case file_type[0] == "eupic1":
			file_dir = "eu_pic"
			file_name = file_name + "_lic"
		case file_type[0] == "eupic2":
			file_dir = "eu_pic"
			file_name = file_name + "_sign"
		}

		f, err := os.OpenFile("./dl/"+file_dir+"/"+file_name+"."+getFileExt(handler.Filename), os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
		fmt.Fprintf(w, "上传完毕，请关闭窗口。")
	}
}

func main() {
	file, err := os.Open("config.xml")
	defer file.Close()
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	v := Config{}
	err = xml.Unmarshal(data, &v)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	fmt.Printf("ws-ip:%q\n", v.WS_IP)
	fmt.Printf("ws_port:%q\n", v.WS_PORT)
	fmt.Printf("ws_outdir:%q\n", v.WS_OUTDIR)
	fmt.Printf("fs-ip:%q\n", v.FS_IP)
	fmt.Printf("fs_port:%q\n", v.FS_PORT)
	fmt.Printf("fs_root:%q\n", v.FS_ROOT)

	fs_ip, fs_port, out_dir = v.FS_IP, v.FS_PORT, v.WS_OUTDIR

	http.Handle("/static/", http.FileServer(http.Dir("public")))
	//ttp.Handle("/js/", http.FileServer(http.Dir("static")))
	//http.Handle("/img/", http.FileServer(http.Dir("static")))
	http.Handle("/dl/", http.FileServer(http.Dir("public")))

	http.HandleFunc("/pdf", convertToPDF)
	http.HandleFunc("/upload", upload)
	err = http.ListenAndServe(":"+v.WS_PORT, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
