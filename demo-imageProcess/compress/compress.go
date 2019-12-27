package main

import (
	"bytes"
	"demo-person/demo-reptile/fromDocument"
	"fmt"
	"github.com/nfnt/resize"
	"image/jpeg"
	"io/ioutil"
	"os"
	"strings"
)

func main() {

	getImgList("/Users/wuchaoqun/Desktop/4")

}

func getImgList(path string) {

	list, _ := ioutil.ReadDir(path)
	dList := make([]string, 0)
	for _, v := range list {
		if m := strings.LastIndex(v.Name(), ".DS_Store"); m == -1 {
			a := path + "/" + v.Name()
			dList = append(dList, a)
		}
	}
	fmt.Println(dList)

	//fmt.Println("jeguo",test3(dList[3]))

	for i, _ := range dList {

		fmt.Println("jeguo", len(test3(dList[i])))
		fmt.Println("压缩后", len(fromDocument.ImageCompression(dList[i])))

	}
}
func test3(path string) (backimg []byte) {
	f, err := os.Open(path)
	defer f.Close()
	a, _ := ioutil.ReadAll(f)
	fmt.Println("a:", len(a))
	b, err := jpeg.Decode(bytes.NewReader(a))
	if err != nil {
		fmt.Println("decode err:", err)
		return
	}
	//imgback = b.Bounds().Dx()

	t := bytes.NewBuffer(backimg)

	m := resize.Resize(400, 300, b, resize.Lanczos3)

	jpeg.Encode(t, m, nil)
	backimg = t.Bytes()

	return
}

func test2() {
	f, err := os.Open("/Users/wuchaoqun/Pictures/1.jpg")
	if err != nil {
		fmt.Println("err is", err)
		return
	}
	defer f.Close()

	bytesf, _ := ioutil.ReadAll(f)
	b, _ := jpeg.Decode(bytes.NewReader(bytesf))

	newdx := 410
	dx := b.Bounds().Dx()
	dy := b.Bounds().Dy()

	file_out, _ := os.Create("/Users/wuchaoqun/Desktop/2.jpg")

	defer file_out.Close()

	m := resize.Resize(uint(newdx), uint(newdx*dy/dx), b, resize.Lanczos3)

	//a := image.NewRGBA(image.Rect(0, 0, newdx, newdx*dy/dx))

	//draw.Draw(a, a.Bounds(), b, image.Pt(0,0), draw.Src)

	jpeg.Encode(file_out, m, nil)
	//fmt.Println(base64.StdEncoding.EncodeToString(a.Bytes()))
	return

}

func test1() {
	f, err := os.Open("/Users/wuchaoqun/Pictures/1.jpg")
	if err != nil {
		fmt.Println("err is", err)
		return
	}
	defer f.Close()

	bytesf, err := ioutil.ReadAll(f)

	//每次需要使用图像数据是都需要重新创建读取器，否则会提示“invalid JPEG format: missing SOI marker”错误，Reader接口在读取数据后并不不会重置读取游标

	b, _ := jpeg.DecodeConfig(bytes.NewReader(bytesf))
	fmt.Println(b.Width)

	imgs, _ := jpeg.Decode(bytes.NewReader(bytesf))

	file_out, err := os.Create("/Users/wuchaoqun/Desktop/2.jpg")
	q := jpeg.Options{Quality: 20}
	defer file_out.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	jpeg.Encode(file_out, imgs, &q)
}
