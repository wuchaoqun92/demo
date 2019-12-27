package main

import (
	"bytes"
	"demo-person/demo-OpenApiTest/faceRecognition/faceRecog/SmallComponents"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	//todo("/Users/wuchaoqun/Desktop/codeMaterial/pic/7.jpg","/Users/wuchaoqun/Desktop/codeMaterial/pic")
	MultiShearImage("/Users/wuchaoqun/Desktop/codeMaterial/pic")
	cost := time.Since(start)
	fmt.Println("cost time:", cost)
}

func todo(path, dir string) {
	img := SmallComponents.ReadImage(path)
	if img.Format == "" {
		log.Println("read image err")
		return
	}

	imgByte := SmallComponents.FaceRec(img.ImageByte)

	faceDet := SmallComponents.ResultUnmarshal(imgByte)
	if len(faceDet.Result.List) == 0 {
		log.Println("ummarshal err")
		return
	}

	x := SmallComponents.ImageCut(faceDet, img)

	checkList := picList(dir)
	fmt.Println("人脸检测结束，开始校对人像")
	resList := make([]SmallComponents.FinalRes, 0)
	for _, v := range checkList {
		if v == path {
			continue
		}
		go SmallComponents.FaceCheck(x, v)
		a := <-SmallComponents.Res
		resList = append(resList, a)
	}

	fmt.Println("----------------对比结果--------------")
	for _, v := range resList {
		t := SmallComponents.CheckResUnmarshal(v.Res, v.Path)
		fmt.Printf("对比图：%s，对比值：%f\n", t.ComparisonImg, t.Result)
	}

}

//Users/wuchaoqun/Desktop/codeMaterial/pic
func picList(path string) (a []string) {

	list, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("readDir failed,err is", err)
		return
	}
	a = make([]string, 0)

	for _, v := range list {
		if v.IsDir() {
			continue
		} else if m := strings.LastIndex(v.Name(), ".DS_Store"); m == -1 { //mac会自动.DS_Store文件，此处进行剔除
			a = append(a, path+"/"+v.Name())
		}
	}

	//fmt.Println(a)
	return
}

func ShearImage(inPath, outPath string) {
	img := SmallComponents.ReadImage(inPath)
	if img.Format == "" {
		log.Println("read image err")
		return
	}

	imgByte := SmallComponents.FaceRec(img.ImageByte)
	faceDet := SmallComponents.ResultUnmarshal(imgByte)
	if len(faceDet.Result.List) == 0 {
		log.Println("ummarshal err")
		return
	}

	x := SmallComponents.ImageCut(faceDet, img)
	y, _ := os.Create(outPath)
	wrilen, _ := io.Copy(y, bytes.NewReader(x))
	fmt.Println("writen len", wrilen)

}

func MultiShearImage(listPath string) {
	cutList := picList(listPath)
	for _, v := range cutList {
		output := fmt.Sprintf("%s-0.jpg", v)
		fmt.Printf("input：%s,output:%s\n", v, output)
		ShearImage(v, output)
	}
	fmt.Println("切图完毕。存储路径为同文件夹下，文件名-0.jpg 文件")
}

var testChan = make(chan int, 10)

func t(i int, wg *sync.WaitGroup) {
	i = i * i
	time.Sleep(time.Second * 3)
	testChan <- i
	wg.Done()
}

func testGoStack() {
	var list = make([]int, 10)
	for i := 0; i < 10; i++ {
		list[i] = i
	}

	var l = make([]int, 10)
	s := SmallComponents.CreateStack()
	var wg sync.WaitGroup
	for _, v := range list {
		wg.Add(1)
		s.Push(v)
		if s.IsFull() {
			for !s.IsEmpty() {
				go t(int(reflect.ValueOf(s.Pop()).Int()), &wg)
				l[0] = <-testChan
				fmt.Println("计数：", s.Pop())
			}
		} else {
			continue
		}
	}
	wg.Wait()
	close(testChan)

	for v := range testChan {
		fmt.Println("result", v)
	}
}
