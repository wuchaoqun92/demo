package fromDocument

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/unidoc/unioffice/document"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type documentList struct {
	Title string
	Path  string
}
type DocContent struct {
	RowNum      int
	ContentType string
	Values      string
}
type res struct {
	Title  string
	DocCon []DocContent
}

var docListChan = make(chan []documentList)
var docChan = make(chan res, 10)
var chanFlag = true
var errcode = ""

func Main2(w http.ResponseWriter, r *http.Request) {

	path := r.FormValue("path")
	if path == "" {
		errcode = "未填写文件路径"
		io.WriteString(w, errcode)
		return
	}
	_, err := ioutil.ReadDir(path)
	if err != nil {
		errcode = err.Error()
		io.WriteString(w, errcode)
		return
	}

	if !chanFlag {
		fmt.Println("chan closed,now open")
		Open(docChan)
	}

	list, errcode := docList(path)
	if errcode != "" {
		errcode += fmt.Sprintf("以上文件后缀错误，请使用.doxc\n")
		io.WriteString(w, errcode)
		return
	}
	start := time.Now()
	//some func or operation
	var wg sync.WaitGroup
	for _, v := range list {
		wg.Add(1)
		go docParse(&wg, v.Path, v.Title)
	}
	wg.Wait()
	close(docChan)
	chanFlag = false
	cost := time.Since(start)

	fmt.Println("All document parse finished executing，now back data")
	fmt.Printf("cost time=[%s]\n", cost)

	var dataBack string

	for result := range docChan {
		errcode += backDataCheck(result)
		res, _ := json.Marshal(result)
		dataBack += errcode + string(res)
	}

	io.WriteString(w, errcode+string(dataBack))
	fmt.Println("success back data")
	return
}

func backDataCheck(result res) (errcode string) {

	for _, v := range result.DocCon {
		_, err := strconv.Atoi(v.Values)
		if v.ContentType == "pic" && err == nil {
			errcode = fmt.Sprintf("文件名:【%s】,图片获取失败，请检查文件.\n", result.Title)
			break
		}
	}
	return
}

func docList(path string) (result []documentList, errcode string) {

	list, _ := ioutil.ReadDir(path)
	a := make([]string, 0)

	for _, v := range list {
		if v.IsDir() {
			continue
		} else {
			a = append(a, v.Name())
		}
	}

	flag := true
	for _, v := range a {
		n := strings.LastIndex(v, ".docx")
		if n == -1 {
			errcode += fmt.Sprintf("文件名：【%s】\n", v)
			flag = false
		}
	}

	if !flag {
		for _, v := range a {
			fmt.Println(v)
		}
		return
	}
	result = make([]documentList, len(a))
	for i, v := range a {
		result[i].Title = v
		result[i].Path = path + "/" + v
	}
	return
}

func docParse(wg *sync.WaitGroup, path, docTitle string) {
	fmt.Println("开始解析文件", path)
	dox, _ := document.Open(path)

	content_list := make([]DocContent, len(dox.Paragraphs()))

	i := 0
	for j, para := range dox.Paragraphs() {
		value := ""
		var b strings.Builder
		for _, run := range para.Runs() {
			if run.Text() == "" {
				value = strconv.Itoa(i)
				content_list[j].ContentType = "pic"
				i++
			} else {
				b.WriteString(run.Text())
				value = b.String()
				content_list[j].ContentType = "text"
			}
		}
		content_list[j].RowNum = j
		content_list[j].Values = value

	}
	for i, v := range dox.Images {
		fi, _ := os.Open(v.Path())
		stat, _ := fi.Stat()
		a := make([]byte, stat.Size())
		fi.Read(a)
		for j, _ := range content_list {
			if content_list[j].Values == strconv.Itoa(i) {
				b := base64.StdEncoding.EncodeToString(a)
				content_list[j].Values = b[:10]
			}
		}
	}
	result := res{
		Title:  docTitle,
		DocCon: content_list,
	}
	docChan <- result
	wg.Done()
}
