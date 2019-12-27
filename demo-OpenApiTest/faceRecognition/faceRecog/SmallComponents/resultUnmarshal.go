package SmallComponents

import (
	"encoding/json"
	"fmt"
	"log"
)

type faceDetect struct {
	Error_code int  `json:"error_code"`
	Result     list `json:"result"`
}

type list struct {
	List []info `json:"face_list"`
}

type info struct {
	Age        int    `json:"age"`
	Face_token string `json:"face_token"`
	Gen        gen    `json:"gender"`
	Loc        loc    `json:"location"`
}
type gen struct {
	Gender string `json:"type"`
}
type loc struct {
	Left   float32 `json:"left"`
	Top    float32 `json:"top"`
	Width  float32 `json:"width"`
	Height float32 `json:"height"`
}

func ResultUnmarshal(res []byte) (t faceDetect) {
	var verd map[string]interface{}
	err := json.Unmarshal(res, &verd)
	if err != nil {
		log.Println("unmarshal map failed,err is", err)
		return
	}
	if verd["error_code"] != float64(0) {
		log.Println("facedetect back error msg:", verd["error_msg"])
		return
	}
	err = json.Unmarshal(res, &t)
	if err != nil {
		log.Println("unmarshal failed,err is", err)
		return
	} else {
		log.Println("人脸识别反序列化结果：", t)
	}
	return
}

type score struct {
	Score float32 `json:"score"`
}

type CheckRes struct {
	ComparisonImg string
	Result        score `json:"result"`
}

func CheckResUnmarshal(res []byte, checkPath string) (t CheckRes) {
	var verd map[string]interface{}
	err := json.Unmarshal(res, &verd)
	if err != nil {
		log.Println("unmarshal map failed,err is", err)
		return
	}

	if verd["error_code"] != float64(0) {
		t.ComparisonImg = fmt.Sprintf("%s check failed;err is %s", checkPath, verd["error_msg"])
		t.Result.Score = -1
		//log.Printf("facecheck back error msg:%s,path is %s",verd["error_msg"],checkPath)
		return
	}

	//fmt.Println("原始校对结果",string(res))
	err = json.Unmarshal(res, &t)
	if err != nil {
		log.Println("unmarshal failed,err is", err)
		return
	}
	t.ComparisonImg = checkPath
	return
}
