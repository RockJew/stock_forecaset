package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	StockGz   string = "sh000012"
	Stock300  string = "sh000300"
	StockNSDK string = "sz159632"

	StockCycle int = 22

	NoticeApptoken string = "AT_r9L2ORNpkVOy9HdoQ5YHra8ag10oDjX8"
)

// 发送通知
func pushNotification(content, summary string) {
	var targetUrl = "https://wxpusher.zjiecode.com/api/send/message"

	client := resty.New()
	body := map[string]interface{}{
		"appToken":    NoticeApptoken,
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
func getCurrentValue(stockCode string) (value float64, err error) {
	currentDataString := getCurrentData(stockCode)

	data := strings.Split(currentDataString, ",")
	if len(data) >= 4 {
		value, err = strconv.ParseFloat(data[3], 64)
		if err == nil {
			fmt.Println(value)
		} else {
			fmt.Println("Error parsing float:", err)
		}
	} else {
		fmt.Println("Invalid data format.")
	}
	return value, err
}

func getCurrentData(stockCode string) string {
	var url = "https://hq.sinajs.cn/list=%s"
	var targetUrl = fmt.Sprintf(url, stockCode)

	client := resty.New()
	resp, _ := client.R().
		SetHeader("Referer", "https://finance.sina.com.cn").
		Get(targetUrl)

	re := regexp.MustCompile(`"(.*?)"`)
	match := re.FindStringSubmatch(resp.String())
	dataString := ""

	if len(match) > 1 {
		dataString = match[1]
		//fmt.Println(dataString)
	} else {
		fmt.Println("No match found.")
	}

	return dataString
}

// 获取股票历史数据
type historyDataUnit struct {
	Day    string `json:"day"`
	Open   string `json:"open"`
	High   string `json:"high"`
	Low    string `json:"low"`
	Close  string `json:"close"`
	Volume string `json:"volume"`
}

func getHistoryValue(stockCode string, daysBefore int) (value float64, err error) {
	historyDataString := getHistoryData(stockCode)

	var data []historyDataUnit
	err = json.Unmarshal([]byte(historyDataString), &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return 0.0, err
	}

	sort.Slice(data, func(i, j int) bool {
		t1, _ := time.Parse("2006-01-02", data[i].Day)
		t2, _ := time.Parse("2006-01-02", data[j].Day)
		return t1.After(t2)
	})

	if len(data) >= daysBefore {
		value, err = strconv.ParseFloat(data[daysBefore-1].Close, 64)
		fmt.Println(value)
	} else {
		fmt.Println("Insufficient data.")
	}

	return value, err
}

func getHistoryData(stockCode string) string {
	var url = "https://quotes.sina.cn/cn/api/jsonp_v2.php/data/CN_MarketDataService.getKLineData?symbol=%s&scale=240&ma=no&datalen=24"
	var targetUrl = fmt.Sprintf(url, stockCode)

	client := resty.New()
	resp, _ := client.R().Get(targetUrl)

	re := regexp.MustCompile(`data\((.*?)\);`)
	match := re.FindStringSubmatch(resp.String())
	dataString := ""

	if len(match) > 1 {
		dataString = match[1]
		//fmt.Println(dataString)
	} else {
		fmt.Println("No match found.")
	}
	return dataString
}

// 策略算法
func stockStrategy() {

}

func main() {
	getHistoryValue(StockNSDK, StockCycle)
}
