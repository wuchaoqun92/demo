package main

import (
	"bufio"
	"bytes"
	"demo-person/demo-miniprogramme/method"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {

	http.HandleFunc("/checkIn", Panics(method.CheckIn))
	http.HandleFunc("/GetContentList", Panics(method.ChoiceMethodToList))

	http.HandleFunc("/upload", Panics(download))

	err := http.ListenAndServe(":8888", nil)
	fmt.Println("shshshshshs", err)
	if err == nil {
		fmt.Println("shshshshshs", err)
	}

}

//panic预处理
func Panics(handle http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if x := recover(); x != nil {
				log.Printf("[%v]caught panic:%v", r.RemoteAddr, x)
			}
		}()
		handle(w, r)
	}
}

func download(w http.ResponseWriter, r *http.Request) {
	//var buf [30]byte
	//str := r.FormValue("body")
	//x := r.Body
	//bytes, _ := ioutil.ReadAll(x)
	//fmt.Println("recv data is", string(bytes))

	path := "/Users/wuchaoqun/Desktop/" + strconv.Itoa(123) + ".mp4"
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()

	defer r.Body.Close()

	bytes, err := ioutil.ReadAll(r.Body)
	err = ioutil.WriteFile(path, bytes, 0666)
	if err != nil {
		fmt.Println("write failed ,err is", err)
		return
	}

	//var result string
	//
	//buf := make([]byte, 4096)
	//res := make([]byte, 0)
	//for {
	//	n, err := r.Body.Read(buf)
	//	fmt.Println("读取长度：",n)
	//	fmt.Println("对应长度内容",buf[:n])
	//	if err != nil {
	//		break
	//	}
	//	result += string(buf[:n])
	//	f.Write(buf[:n])
	//	res = BytesCombine(res,buf[:n])
	//}
	//err = ioutil.WriteFile(path, res, 0666)
	//if err != nil {
	//	fmt.Println("write failed ,err is",err)
	//	return
	//}

}

func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

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
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)

	return bytes, err
}
