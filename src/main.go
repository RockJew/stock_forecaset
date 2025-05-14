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
	"os/exec"
)

const (
	StockGold string = "sh518880"
	Stock300  string = "sh510300"
	StockNSDK string = "sz159632"
	StockYFD  string = "sz159915"

	StockCycle int = 22

	NoticeApptoken string = "AT_r9L2ORNpkVOy9HdoQ5YHra8ag10oDjX8"
)

// 发送通知
func pushNotification(content, summary string) {
	target_user = "wya97@icould.com"
	exec_shell(target_user, summary)
	exec_shell(target_user, content)
	// var targetUrl = "https://wxpusher.zjiecode.com/api/send/message"

	// client := resty.New()
	// body := map[string]interface{}{
	// 	"appToken":    NoticeApptoken,
	// 	"content":     content,
	// 	"summary":     summary,
	// 	"contentType": 1,
	// 	"topicIds":    []int{25109},
	// }

	// resp, _ := client.R().
	// 	SetHeader("Content-Type", "application/json").
	// 	SetBody(body).
	// 	Post(targetUrl)

	// fmt.Print(resp)
}

func exec_shell(target_user string, content string) error {
	// 使用 fmt.Sprintf 格式化 shell_code
	shell_code := fmt.Sprintf(
		`tell application "Messages"
			set targetService to 1st service whose service type = iMessage
			set targetBuddy to buddy "%s" of targetService
			send "%s" to targetBuddy
		end tell`, target_user, content)

	// 创建 osascript 命令
	cmd := exec.Command("osascript", "-e", shell_code)

	// 运行命令并返回错误
	return cmd.Run()
}

// 获取股票当前值
func getCurrentValue(stockCode string) (value float64, err error) {
	currentDataString := getCurrentData(stockCode)

	data := strings.Split(currentDataString, ",")
	if len(data) >= 4 {
		value, err = strconv.ParseFloat(data[3], 64)
		if err == nil {
			//fmt.Println(value)
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
		//fmt.Println(value)
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
type stockCalculateUnit struct {
	Code         string
	Name         string
	currentValue float64
	historyValue float64
	rateOfReturn float64
}

func stockStrategy() (err error) {
	stockStructList := []stockCalculateUnit{
		{Code: StockGold, Name: "黄金"},
		{Code: StockYFD, Name: "创业板ETF"},
		{Code: StockNSDK, Name: "美指"},
	}

	for index, _ := range stockStructList {
		stockStructList[index].currentValue, err = getCurrentValue(stockStructList[index].Code)
		stockStructList[index].historyValue, err = getHistoryValue(stockStructList[index].Code, StockCycle)
		stockStructList[index].rateOfReturn = (stockStructList[index].currentValue - stockStructList[index].historyValue) / stockStructList[index].historyValue
	}

	sort.Slice(stockStructList, func(i, j int) bool {
		return stockStructList[i].rateOfReturn-stockStructList[j].rateOfReturn > 0
	})

	//fmt.Println(stockStructList)
	content, summary := getNoticeContent(stockStructList)
	pushNotification(content, summary)
	return nil
}

func getNoticeContent(stockStructList []stockCalculateUnit) (content, summary string) {
	if stockStructList[0].rateOfReturn < 0 {
		summary = "现金"
	} else {
		summary = fmt.Sprintf("调仓: %s %s", stockStructList[0].Code, stockStructList[0].Name)
	}

	contentModel := "现状\n" +
		"%s %s %.2f%%\n" +
		"%s %s %.2f%%\n" +
		"%s %s %.2f%%\n" +
		"操作\n" +
		"%s"

	content = fmt.Sprintf(contentModel, stockStructList[0].Code, stockStructList[0].Name, stockStructList[0].rateOfReturn*100,
		stockStructList[1].Code, stockStructList[1].Name, stockStructList[1].rateOfReturn*100,
		stockStructList[2].Code, stockStructList[2].Name, stockStructList[2].rateOfReturn*100,
		summary)

	return content, summary
}

func main() {
	stockStrategy()
}
