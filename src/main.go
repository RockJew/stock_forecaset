package main

import "github.com/go-resty/resty/v2"

func httpPost() {
	appToken := "AT_r9L2ORNpkVOy9HdoQ5YHra8ag10oDjX8"
	content := "测试"

	targetUrl := "https://wxpusher.zjiecode.com/api/send/message"

	client := resty.New()
	body := map[string]interface{} {
		"appToken": appToken,
		"content": content,
		"summary": content,
		"contentType": 1,
		"topicIds": []int{68682},
	}

	resp, err := client.R().
		SetHeader("Content-Type"j, "application/json").

}

func main() {

}
