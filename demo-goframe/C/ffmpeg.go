package main

/*
#cgo CFLAGS: -I/usr/local/Cellar/ffmpeg/4.0/include
#cgo LDFLAGS: -L/usr/local/Cellar/ffmpeg/4.0/lib -lavformat
#include "libavformat/avformat.h"
#include "libavcodec/avcodec.h"
#include "libavutil/avutil.h"
#include "libavutil/opt.h"
#include "libavdevice/avdevice.h"
#include "ffmpeg.h"
*/
import "C"
import "fmt"

func main() {
	//C.test(C.int(12))
	//C.test3()
	time, err := C.test4(C.CString("/Users/wuchaoqun/Documents/Work/缺陷跟踪/1m.mp4"))
	fmt.Printf("video's duration: %d 秒,err:%s\n", time, err)
}
