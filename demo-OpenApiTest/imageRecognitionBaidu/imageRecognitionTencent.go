package main

import (
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tiia "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tiia/v20190529"
)

func main() {

	credential := common.NewCredential(
		"your id",
		"your key",
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "tiia.tencentcloudapi.com"
	client, _ := tiia.NewClient(credential, "ap-shanghai", cpf)

	request := tiia.NewRecognizeCarRequest()

	params := "{\"ImageUrl\":\"https://car2.autoimg.cn/cardfs/product/g29/M06/1E/A1/1024x0_1_q95_autohomecar__ChcCSF12PzyAXbH6AAcY1qcikxs857.jpg\"}"
	err := request.FromJsonString(params)
	if err != nil {
		panic(err)
	}
	response, err := client.RecognizeCar(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", response.ToJsonString())

}
