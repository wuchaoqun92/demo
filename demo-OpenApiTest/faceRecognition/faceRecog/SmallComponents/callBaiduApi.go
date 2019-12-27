package SmallComponents

import (
	"demo-person/demo-OpenApiTest/imageRecognitionBaidu/imageEncode"
	"demo-person/demo-httpCommon/netRequest"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

const (
	accessToken = "24.8327b3342a167999dec081c09db7fa65.2592000.1578204250.282335-17706172"
	recogUrl    = "https://aip.baidubce.com/rest/2.0/face/v3/detect" //人脸检测
	recogUrl2   = "https://aip.baidubce.com/rest/2.0/face/v3/match"  //人脸对比
	requestUrl  = recogUrl + "?access_token=" + accessToken
	requestUrl2 = recogUrl2 + "?access_token=" + accessToken
)

func FaceRec(imageByte []byte) (result []byte) {

	imageValue := base64.StdEncoding.EncodeToString(imageByte)

	header := make(map[string]string, 0)
	header["Content-Type"] = "application/json"

	body := make(map[string]string, 0)
	body["image"] = imageValue
	body["max_face_num"] = "1"
	body["image_type"] = "BASE64"
	body["face_field"] = "age,quality,gender"
	//body["face_field"] = "age,beauty,expression,face_shape,gender,glasses,race,quality,eye_status,emotion,face_type"
	body["face_type"] = "LIVE"

	result = netRequest.PostRequest(requestUrl, header, body)

	//fmt.Println("result is", string(result))
	return
}

func FaceCheck(baseImg []byte, pic2 string) {
	fmt.Println("开始对比", pic2)

	imageValue1 := base64.StdEncoding.EncodeToString(baseImg)

	header := make(map[string]string, 0)
	header["Content-Type"] = "application/json"

	body := make(map[string]string, 0)
	body["image"] = imageValue1
	body["image_type"] = "BASE64"
	body["quality_control"] = "LOW"
	body["face_type"] = "LIVE"
	//body["liveness_control"] = "HIGH"

	imageByte2, _ := imageEncode.LoadImage(pic2)
	imageValue2 := base64.StdEncoding.EncodeToString(imageByte2)

	body1 := make(map[string]string, 0)
	body1["image"] = imageValue2
	body1["image_type"] = "BASE64"
	body1["quality_control"] = "LOW"
	body1["face_type"] = "LIVE"
	//body1["liveness_control"] = "HIGH"

	list := make([]map[string]string, 2)
	list[0] = body
	list[1] = body1

	str, _ := json.Marshal(list)

	result := netRequest.PostRequest(requestUrl2, header, str)
	res := FinalRes{
		Res:  result,
		Path: pic2,
	}

	Res <- res

}

type FinalRes struct {
	Res  []byte
	Path string
}

var Res = make(chan FinalRes, 1024)
