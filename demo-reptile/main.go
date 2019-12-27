package main

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"github.com/unidoc/unioffice/document"
	"io/ioutil"
	"os"
	"regexp"
)

func main() {

	try3()
}

type SConfig struct {
	XMLName xml.Name `xml:"document"` // 指定最外层的标签为config
	a       string   `xml:"t"`        // 读取smtpServer配置项，并将结果保存到SmtpServer变量中
}

func try3() {

	f, _ := os.Open("/Users/wuchaoqun/Desktop/1/beijing/word/document.xml")
	defer f.Close()
	fi, _ := f.Stat()

	by := make([]byte, fi.Size())
	f.Read(by)
	fmt.Println(string(by))

	value := []SConfig{}
	err := xml.Unmarshal(by, &value)
	if err != nil {
		fmt.Println("err is", err)
		return
	}
	fmt.Println(value)

	for _, v := range value {
		fmt.Println(v.a)
	}

}

func try2() {
	f, _ := os.Open("/Users/wuchaoqun/Desktop/1/beijing/word/document.xml")
	defer f.Close()
	fi, _ := f.Stat()

	by := make([]byte, fi.Size())
	f.Read(by)

	//commentCount := `<span style="(.*?)">(.*?)</span>`
	commentCount := `r:embed="(.*?)"`
	reg := regexp.MustCompile(commentCount)

	for _, d := range reg.FindAllString(string(by), -1) {
		fmt.Println(d)
	}

}

func try1() {
	f, _ := os.Open("/Users/wuchaoqun/Desktop/1/beijing.docx")

	fi, _ := os.Stat("/Users/wuchaoqun/Desktop/1/beijing.docx")

	zr, _ := zip.NewReader(f, fi.Size())

	files := []*zip.File{}
	files = append(files, zr.File...)

	for _, f := range files {
		fmt.Println()
		if f.FileHeader.Name == "word/document.xml" {
			content, err := ioutil.ReadFile(f.FileHeader.Name)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(content)
		}
	}

	dox, _ := document.Open("/Users/wuchaoqun/Desktop/1/beijing.docx")
	fmt.Println(dox.Images[0].Path())
}
