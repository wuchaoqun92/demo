package imageEncode

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
)

func ImageEncodeToBase64(bytes []byte) (string, *url.URL) {

	a := base64.StdEncoding.EncodeToString(bytes)
	image_value, _ := url.Parse(a) //获取图片地址的指针
	//fmt.Println(len(a))
	return a, image_value

}

func LoadImage(path string) (bytes []byte, format string) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	bytes, _ = ioutil.ReadAll(f)
	picSOI := fmt.Sprintf("%#x", bytes[:2])
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

	//stats, err := f.Stat() 不同的方法读取图片内容，此方法适用于分片读取
	//if err != nil {
	//	fmt.Println("get file stat faile,err is", err)
	//	return
	//}
	//
	//var size int64 = stats.Size()
	//fmt.Println("size is", size)
	//bytes = make([]byte, size)
	//bufr := bufio.NewReader(f)
	//_, err = bufr.Read(bytes)

	return
}
