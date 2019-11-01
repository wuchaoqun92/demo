package fromDocument

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/unidoc/unioffice/common"
	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unioffice/measurement"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Content struct {
	RowNum      int
	ContentType string
	Values      string
}

type List struct {
	Title string
	text  []Content
}

type ChanTest struct {
	id   int
	name string
}

var testChan = make(chan ChanTest, 50)

func main1() {

	//allDo("/Users/wuchaoqun/Desktop/beijing.docx") 对应地址的文件解析后返回json
	// fileDirectory("/Users/wuchaoqun/Desktop/1")
	var x []ChanTest = make([]ChanTest, 10)
	for i := 0; i < 10; i++ {
		x[i].id = i
		x[i].name = fmt.Sprintf("test%d", i)
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go chanTest(&wg, x[i])
	}
	wg.Wait()
	close(testChan)

	fmt.Println("All go routines finished executing")

	for result := range testChan {
		//for i,v := range result{
		fmt.Printf("id is %d, name is %s\n", result.id, result.name)
		//}

	}

}

func chanTest(wg *sync.WaitGroup, x ChanTest) {

	//b := make([]ChanTest,0)
	var b ChanTest

	b.id = x.id
	b.name = x.name

	testChan <- b
	wg.Done()
}

func fileDirectory(path string) string {
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
			fmt.Printf("文件名：%s，后缀包含错误，请使用.docx", v)
			flag = false
		}
	}
	if !flag {
		return "error"
	}
	result := make([]List, len(a))
	for i, v := range a {
		result[i].Title = v
		result[i].text = allDo(path + "/" + v)
	}

	y, _ := json.Marshal(result)
	return string(y)
}

func allDo(path string) []Content {

	dox, _ := document.Open(path)

	content_list := make([]Content, len(dox.Paragraphs()))

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
				content_list[j].Values = base64.StdEncoding.EncodeToString(a)
			}
		}
	}

	//contentChan <- content_list

	//for _,v := range content_list  {
	fmt.Println(content_list[0])

	//y,_ := json.Marshal(content_list)
	//fmt.Println(string(y))
	return content_list

}

func test2(path string) {

	fi, _ := os.Open(path)
	defer fi.Close()
	stat, _ := fi.Stat()
	a := make([]byte, stat.Size())
	fi.Read(a)

	f, _ := os.Create("/Users/wuchaoqun/Desktop/123.png")
	defer f.Close()
	f.Write(a)
}

func test() {
	doc, _ := document.Open("/Users/wuchaoqun/Desktop/image.docx")

	img := doc.Images[0]

	fi, _ := os.Open(img.Path())
	stat, _ := fi.Stat()
	a := make([]byte, stat.Size())
	fi.Read(a)

	_, err := os.Create("/Users/wuchaoqun/Desktop/111.png")
	if err != nil {
		return
	}
	ioutil.WriteFile("/Users/wuchaoqun/Desktop/111.png", a, 0666)

}

func insertDoc() {

	dox, _ := document.Open("/Users/wuchaoqun/Desktop/image.docx")

	a := dox.Images
	fmt.Println(a)
	b, _ := common.ImageFromFile("/Users/wuchaoqun/Desktop/1.png")
	fmt.Println(b)
	img1ref, _ := dox.AddImage(b)
	fmt.Println(img1ref)

	para := dox.AddParagraph()

	run := para.AddRun()

	for i := 0; i < 16; i++ {
		run.AddText("lorem")

		// drop an inline image in
		if i == 13 {
			inl, err := run.AddDrawingInline(img1ref)
			if err != nil {
				log.Fatalf("unable to add inline image: %s", err)
			}
			inl.SetSize(1*measurement.Inch, 1*measurement.Inch)
			x, _ := inl.GetImage()
			fmt.Println("www:", x)
		}

	}
	dox.SaveToFile("/Users/wuchaoqun/Desktop/image1.docx")
}

func readDoc() {

	dox, _ := document.Open("/Users/wuchaoqun/Desktop/beijing.docx")

	cp := dox.CoreProperties
	// You can read properties from the document
	fmt.Println("Title:", cp.Title())
	fmt.Println("Author:", cp.Author())
	fmt.Println("Description:", cp.Description())
	fmt.Println("Last Modified By:", cp.LastModifiedBy())
	fmt.Println("Category:", cp.Category())
	fmt.Println("Content Status:", cp.ContentStatus())
	fmt.Println("Created:", cp.Created())
	fmt.Println("Modified:", cp.Modified())

	a := dox.Images
	fmt.Println(a)

	for i, para := range dox.Paragraphs() {
		//run为每个段落相同格式的文字组成的片段
		fmt.Println("-----------第", i, "段-------------")
		for _, run := range para.Runs() {
			//fmt.Print("\t-----------第", j, "格式片段-------------")
			fmt.Print(run.Text())
			if run.Text() == "" {
				fmt.Print(run.X().RsidRAttr)
			}

		}
		fmt.Println()
	}
}
