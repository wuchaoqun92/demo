package imageEncode

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
)

func ImageEncodeToBase64(bytes []byte) (string, *url.URL) {

	a := base64.StdEncoding.EncodeToString(bytes)
	image_value, _ := url.Parse(a) //获取图片地址的指针
	//fmt.Println(len(a))
	return a, image_value

}

func LoadImage(path string) (bytes []byte) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}

	stats, err := f.Stat()
	if err != nil {
		fmt.Println("get file stat faile,err is", err)
		return
	}

	var size int64 = stats.Size()
	fmt.Println("size is", size)
	bytes = make([]byte, size)

	bufr := bufio.NewReader(f)
	_, err = bufr.Read(bytes)
	//fmt.Println("bytes is",string(bytes))
	return
}
