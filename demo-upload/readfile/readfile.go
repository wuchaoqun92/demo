package main

import (
	"bufio"
	"bytes"
	"demo-person/demo-upload/circleBuffer"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

func main() {
	//b,_ := RetrieveROM("/Users/wuchaoqun/Documents/Work/缺陷跟踪/4m.mp4")
	//fmt.Println("b is",b)
	//writefile(b,"/Users/wuchaoqun/Desktop/1.mp4")
	c, _ := RetrieveROM2("/Users/wuchaoqun/Documents/Work/缺陷跟踪/1m.mp4")
	//fmt.Println("c is",c)
	writefile(c, "/Users/wuchaoqun/Desktop/2.mp4")
}

//整个读取文件
func RetrieveROM(filename string) ([]byte, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		return nil, statsErr
	}

	var size int64 = stats.Size()
	fmt.Println("size is", size)
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)

	return bytes, err

}

//分片读取文件
func RetrieveROM2(filename string) ([]byte, error) {
	body := circleBuffer.NewCircleByteBuffer(1024 * 2)

	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		return nil, statsErr
	}
	bytes2 := make([]byte, 0)

	fsz := float64(stats.Size())
	fupsz := float64(0)
	buf := make([]byte, 1024*4)
	for {
		time.Sleep(10 * time.Microsecond) //减缓上传速度，看进度效果
		n, err := file.Read(buf)
		if n > 0 {
			time.Sleep(10 * time.Microsecond)
			nz, _ := body.Write(buf[0:n])
			fupsz += float64(nz)
			progress := strconv.Itoa(int((fupsz/fsz)*100)) + "%"
			fmt.Println("upload:", progress, "|", strconv.FormatFloat(fupsz, 'f', 0, 64), "/", stats.Size())
			bytes2 = BytesCombine(bytes2, buf[:n])
		}
		if err == io.EOF {
			break
		}
	}

	//buf:=make([]byte,1024*4)//每次读取长度
	//for {
	//	n,err := file.Read(buf)
	//	if n > 0 {
	//		bytes2=BytesCombine(bytes2,buf[:n])
	//	}
	//	if err == io.EOF{
	//		break
	//	}
	//}
	return bytes2, err

}

func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

func writefile(b []byte, path1 string) {
	//b,_ := RetrieveROM("/Users/wuchaoqun/Pictures/10.jpg")
	//fmt.Println("b is",b)

	//path1 := "/Users/wuchaoqun/Desktop/" + strconv.Itoa(1) + ".jpg"
	_, err := os.Create(path1)
	if err != nil {
		return
	}
	//ioutil.WriteFile(path1,b, 0666)
	ioutil.WriteFile(path1, b, 0666)
}
