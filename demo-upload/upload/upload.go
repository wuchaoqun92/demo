package upload

import (
	"crypto/rand"
	"demo-person/demo-upload/circleBuffer"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func randomBoundary() string {
	var buf [30]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", buf[:])
}

/*
 * url:上传地址
 * flpath:上传文件地址
 */
func Upload(url, flpath string) {
	body := circleBuffer.NewCircleByteBuffer(1024 * 2)
	boundary := randomBoundary()
	boundarybytes := []byte("\r\n--" + boundary + "\r\n")
	endbytes := []byte("\r\n--" + boundary + "--\r\n")

	reqest, err := http.NewRequest("POST", url, body)
	if err != nil {
		panic(err)
	}
	reqest.Header.Add("Connection", "keep-alive")
	reqest.Header.Add("Content-Type", "multipart/form-data; charset=utf-8; boundary="+boundary)
	go func() {
		//defer ruisRecovers("upload.run")
		f, err := os.OpenFile(flpath, os.O_RDONLY, 0666) //其实这里的 O_RDWR应该是 O_RDWR|O_CREATE，也就是文件不存在的情况下就建一个空文件，但是因为windows下还有BUG，如果使用这个O_CREATE，就会直接清空文件，所以这里就不用了这个标志，你自己事先建立好文件。
		if err != nil {
			panic(err)
		}
		stat, err := f.Stat() //获取文件状态
		if err != nil {
			panic(err)
		}
		defer f.Close()

		header := fmt.Sprintf("Content-Disposition: form-data; name=\"upfile\"; filename=\"%s\"\r\nContent-Type: application/octet-stream\r\n\r\n", stat.Name())
		body.Write(boundarybytes)
		body.Write([]byte(header))

		fsz := float64(stat.Size())
		fupsz := float64(0)
		buf := make([]byte, 1024)
		for {
			time.Sleep(10 * time.Microsecond) //减缓上传速度，看进度效果
			n, err := f.Read(buf)
			if n > 0 {
				nz, _ := body.Write(buf[0:n])
				fupsz += float64(nz)
				progress := strconv.Itoa(int((fupsz/fsz)*100)) + "%"
				fmt.Println("upload:", progress, "|", strconv.FormatFloat(fupsz, 'f', 0, 64), "/", stat.Size())
			}
			if err == io.EOF {
				break
			}
		}
		body.Write(endbytes)
		body.Write(nil) //输入EOF,表示数据写完了
	}()
	resp, err := http.DefaultClient.Do(reqest)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		fmt.Println("上传成功")
	} else {
		fmt.Println("上传失败,StatusCode:", resp.StatusCode)
	}
}
