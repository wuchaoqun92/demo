package fromDocument

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/nfnt/resize"
	"github.com/unidoc/unioffice/document"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
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

//var docChan2 chan res
var docChan = make(chan res, 1024) //存储文件解析结果
var chanFlag = true                //记录管道的关闭状态
var errcode = ""                   //记录错误信息

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

	//判断管道是否关闭，若关闭则重新开启
	if !chanFlag {
		fmt.Println("chan closed,now open")
		Open(docChan)
	}

	//检查文件是否存在非docx后缀的文件
	list, errcode := docList(path)
	if errcode != "" {
		errcode += fmt.Sprintf("以上文件后缀错误，请使用.doxc\n")
		io.WriteString(w, errcode)
		return
	}
	start := time.Now()
	//some func or operation
	var wg sync.WaitGroup
	for _, v := range list { //根据文件个数创建goroutine，用于解析文件
		wg.Add(1)
		go docParse(&wg, v.Path, v.Title)
	}
	wg.Wait()
	close(docChan)
	chanFlag = false
	cost := time.Since(start) //文件解析花费时间

	fmt.Println("All document parse finished executing，now back data")
	fmt.Printf("cost time=[%s]\n", cost)

	var dataBack string

	for result := range docChan {
		errcode += backDataCheck(result)
		res, _ := json.Marshal(result)
		dataBack += string(res)
	}

	io.WriteString(w, errcode+string(dataBack))
	fmt.Println("success back data")
	return
}

//此函数用于检查部分文件图片解析情况，unioffice包存在部分文档无法解析出图片的情况，即doc.Image的切片无内容。
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

//获取对应路径下的文件列表
func docList(path string) (result []documentList, errcode string) {

	list, _ := ioutil.ReadDir(path)
	a := make([]string, 0)

	for _, v := range list {
		if v.IsDir() {
			continue
		} else if m := strings.LastIndex(v.Name(), ".DS_Store"); m == -1 { //mac会自动.DS_Store文件，此处进行剔除
			a = append(a, v.Name())
		}
	}

	flag := true //记录是否存在非.docx后缀的文件，用于自查
	result = make([]documentList, 0)

	for _, v := range a {
		n := strings.LastIndex(v, ".docx") //检查文件是否已.docx结尾，若存在非后缀的文件则拼接错误文件用于返回
		if n == -1 {
			errcode += fmt.Sprintf("文件名：【%s】\n", v)
			flag = false
		} else {
			res := documentList{
				Title: v,
				Path:  path + "/" + v,
			}
			result = append(result, res) //拼接实际可解析的文件列表
		}

	}

	fmt.Println("实际待解析列表", result)
	if !flag {
		fmt.Println("所有文件列表：")
		for _, v := range a {
			fmt.Println(v)
		}
		return
	}

	return
}

//检查图片类型，由于借助开源包解析 word 中的图片，故此处无需使用
func isPic(path string) (format string) {
	l := strings.Split(path, ".")
	a := len(l)
	switch l[a-1] {
	case "jpg":
		fallthrough
	case "jpeg":
		format = "jpg"
	case "png":
		format = "png"
	default:
		format = "noPic"
	}
	return
}

func docParse(wg *sync.WaitGroup, path, docTitle string) {
	fmt.Println("开始解析文件", path)
	dox, _ := document.Open(path) //使用unioffice包解析文件

	contentList := make([]DocContent, len(dox.Paragraphs()))

	i := 0
	//unioffice解析文本内容，图片位置使用对应图片id进行替换
	for j, para := range dox.Paragraphs() { //摸索发现，换行位置的para为空切片即为nil，故后续使用runs遍历是会直接结束。
		value := ""
		var b strings.Builder
		for _, run := range para.Runs() { //摸索发现，runs函数可以作为对应到word的document.xml文件中<w:r>标签，即每行的内容
			if run.Text() == "" { //摸索发现，图片位置text的内容为空
				value = strconv.Itoa(i)
				contentList[j].ContentType = "pic"
				i++
			} else {
				b.WriteString(run.Text())
				value = b.String()
				contentList[j].ContentType = "text"
			}
		}
		contentList[j].RowNum = j
		contentList[j].Values = value

	}

	//unioffice解析图片内容，根据文本解析结果中的图片id，将图片信息base64编码后进行替换（unioffice，图片列表根据文中从上到下的顺序进行排列）
	for i, v := range dox.Images {
		start := time.Now()
		imageValue := ImageCompression(v.Path())
		cost := time.Since(start) //文件解析花费时间
		fmt.Printf(" %s,pic  conpress cost time=[%s]\n", path, cost)
		for j, _ := range contentList {
			if contentList[j].Values == strconv.Itoa(i) {
				b := base64.StdEncoding.EncodeToString(imageValue)
				contentList[j].Values = b
				break
			}
		}
	}

	//配置返回结果并插入管道
	result := res{
		Title:  docTitle,
		DocCon: contentList,
	}
	docChan <- result
	wg.Done()
}

func ImageCompression(imagePath string) (backImage []byte) {
	imageValue, format := getPic(imagePath)
	var pic image.Image
	var err error

	switch format {
	case "jpg":
		pic, err = jpeg.Decode(bytes.NewReader(imageValue))
		if err != nil {
			fmt.Println("decode err:", err)
			return
		}
	case "png":
		pic, err = png.Decode(bytes.NewReader(imageValue))
		if err != nil {
			fmt.Println("decode err:", err)
			return
		}
	default:
		fmt.Println("picture wrong err:", err)
		return
	}

	newdx := 410
	dx := pic.Bounds().Dx()
	dy := pic.Bounds().Dy()
	op := jpeg.Options{Quality: 30}

	a := bytes.NewBuffer(backImage)

	m := resize.Resize(uint(newdx), uint(newdx*dy/dx), pic, resize.Lanczos3)

	jpeg.Encode(a, m, &op)

	backImage = a.Bytes()

	return

}

func getPic(path string) (imageValue []byte, format string) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("open file failed,err is", err)
		return
	}
	defer f.Close()

	imageValue, _ = ioutil.ReadAll(f)
	picSOI := fmt.Sprintf("%#x", imageValue[:2])
	png := "0x8950"
	jpg := "0xffd8"
	switch picSOI {
	case png:
		format = "png"
	case jpg:
		format = "jpg"
	default:
		format = "noPic"
	}
	return
}
