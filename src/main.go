package main

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

const STOCK_GZ string = "sh000012"
const STOCK_300 string = "sh000300"
const STOCK_NSDK string = "sz159632"
const STOCK_CYCLE int = 22

const NOTICE_APPTOKEN string = "AT_r9L2ORNpkVOy9HdoQ5YHra8ag10oDjX8"

// 发送通知
func pushNotification(content, summary string) {
	var targetUrl = "https://wxpusher.zjiecode.com/api/send/message"

	client := resty.New()
	body := map[string]interface{}{
		"appToken":    NOTICE_APPTOKEN,
		"content":     content,
		"summary":     summary,
		"contentType": 1,
		"topicIds":    []int{25109},
	}

	resp, _ := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(targetUrl)

	fmt.Print(resp)
}

// 获取股票当前值
func getCurrentValue(stockCode string) float64 {
	var url = "https://hq.sinajs.cn/list=%s"
	var targetUrl = fmt.Sprintf(url, stockCode)

	client := resty.New()
	resp, _ := client.R().
		SetHeader("Referer", "https://finance.sina.com.cn").
		Get(targetUrl)

	fmt.Print(resp)
	return 1.0
}

// 获取股票历史数据
func getHistoryValue(stockCode string, daysBefore int) float64 {
	var url = "https://quotes.sina.cn/cn/api/jsonp_v2.php/data/CN_MarketDataService.getKLineData?symbol=%s&scale=240&ma=no&datalen=24"
	var targetUrl = fmt.Sprintf(url, stockCode)

	client := resty.New()
	resp, _ := client.R().Get(targetUrl)

	fmt.Print(resp)
	return 1.0
}

func main() {
	getCurrentValue(STOCK_GZ)
}
