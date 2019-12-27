package main

import (
	"demo-person/demo-httpCommon/netRequest"
	"encoding/json"
	"fmt"
)

func main() {
	requestUrl2 := "http://0.0.0.0:5050/mypath?a=123&b=234"

	header := make(map[string]string, 0)
	header["Content-Type"] = "application/json222222"

	body := make(map[string]string, 0)
	body["image"] = "imageValue1"

	rBody, _ := json.Marshal(body)

	result := netRequest.PostRequest(requestUrl2, header, string(rBody))

	fmt.Println(string(result))
}
