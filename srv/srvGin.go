package srv

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lmzxtek/ths-go/gm"
)

var (
	localRand *rand.Rand
	gmapi     string = "http://localhost:5000"
	gmcsv     string = "http://localhost:5002"
)

// init 函数在程序启动时自动执行，设置随机数种子
func init() {
	// 创建一个新的本地随机数生成器
	localRand = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func SetURL(gmAPI, gmCSV string) {
	if gmAPI != "" {
		gmapi = SmartURLHandler(gmAPI, false)
	}
	if gmCSV != "" {
		gmcsv = SmartURLHandler(gmCSV, false)
	}
}

// RandomStock 从预定义的股票代码数组中随机返回一个元素
func randomStock() string {
	stocks := []string{"AAPL", "GOOG", "AMZN", "MSFT", "TSLA"}
	// 使用本地的随机数生成器
	return stocks[localRand.Intn(len(stocks))]
}

// randomIntWithRange 返回一个在 [a-b, a+b] 区间内的随机整数
func randomIntWithRange(a, b int) int {
	// 使用本地的随机数生成器生成一个在 [a-b, a+b] 区间内的随机整数
	return localRand.Intn(2*b+1) + a - b
}

// randomDateInLastYear 返回过去一年之内的随机日期
func randomDateInLastYear() string {
	// 获取当前时间
	now := time.Now()
	// 计算一年前的时间
	oneYearAgo := now.AddDate(-1, 0, 0)
	// 计算一年前和现在的时间差（以秒为单位）
	diff := now.Unix() - oneYearAgo.Unix()
	// 生成一个在过去一年内的随机时间点（以秒为单位）
	randomUnix := oneYearAgo.Unix() + int64(localRand.Intn(int(diff)))
	// 将随机时间点转换为时间对象
	randomTime := time.Unix(randomUnix, 0)
	// 格式化随机时间为 "YYYY-MM-DD" 格式
	return randomTime.Format("2006-01-02")
}

func RouteTest(c *gin.Context) {
	// now := time.Now()
	// 格式化当前日期为 "YYYY-MM-DD" 格式
	// today := now.Format("2006-01-02")
	// fmt.Println(" -=> Today is", today)

	// symbols := ["AAPL", "GOOG", "AMZN", "MSFT", "TSLA"]
	// symbol := c.DefaultQuery("symbol", randomStock())
	// time := c.DefaultQuery("time", today)

	// 定义价格和成交量的范围
	pa, pr := 200, 50
	va, vr := 1000, 100

	data := gin.H{
		"Symbol": []string{randomStock(), randomStock(), randomStock(), randomStock(), randomStock()},
		"Time":   []string{randomDateInLastYear(), randomDateInLastYear(), randomDateInLastYear(), randomDateInLastYear(), randomDateInLastYear()},
		"Price":  []int{randomIntWithRange(pa, pr), randomIntWithRange(pa, pr), randomIntWithRange(pa, pr), randomIntWithRange(pa, pr), randomIntWithRange(pa, pr)},
		"Volume": []int{randomIntWithRange(va, vr), randomIntWithRange(va, vr), randomIntWithRange(va, vr), randomIntWithRange(va, vr), randomIntWithRange(va, vr)},
	}
	c.JSON(http.StatusOK, data)
}

func RouteTest2(c *gin.Context) {
	// 定义价格和成交量的范围
	pa, pr := 200, 50
	va, vr := 1000, 100

	data := gin.H{
		"columns": []string{"Symbol", "Time", "Price", "Volume"},
		"data": [][]any{
			{randomStock(), randomDateInLastYear(), randomIntWithRange(pa, pr), randomIntWithRange(va, vr)},
			{randomStock(), randomDateInLastYear(), randomIntWithRange(pa, pr), randomIntWithRange(va, vr)},
			{randomStock(), randomDateInLastYear(), randomIntWithRange(pa, pr), randomIntWithRange(va, vr)},
			{randomStock(), randomDateInLastYear(), randomIntWithRange(pa, pr), randomIntWithRange(va, vr)},
			{randomStock(), randomDateInLastYear(), randomIntWithRange(pa, pr), randomIntWithRange(va, vr)},
		},
	}

	c.JSON(http.StatusOK, data)
}

// 定义一个record结构体，包含每个字段
type record struct {
	Symbol string `json:"Symbol"`
	Time   string `json:"Time"`
	Price  int    `json:"Price"`
	Volume int    `json:"Volume"`
}

func RouteTest3(c *gin.Context) {
	// 定义价格和成交量的范围
	pa, pr := 200, 50
	va, vr := 1000, 100

	// 创建一个包含5个record的切片
	records := make([]record, 5)
	for i := range records {
		records[i] = record{
			Symbol: randomStock(),
			Time:   randomDateInLastYear(),
			Price:  randomIntWithRange(pa, pr),
			Volume: randomIntWithRange(va, vr),
		}
	}

	// 将records切片转换为JSON格式
	// jsonData, err := json.MarshalIndent(records, "", "  ")
	// if err != nil {
	// 	fmt.Println("Error converting to JSON:", err)
	// 	// return c.JSON(http.StatusOK, err.Error())
	// }

	// 输出JSON数据
	// fmt.Println(string(jsonData))

	// c.JSON(http.StatusOK, string(jsonData))
	c.JSON(http.StatusOK, records)
}

func RouteCalendar2(c *gin.Context) {
	yy := fmt.Sprintf("%d", time.Now().Year())
	syear := c.DefaultQuery("syear", yy)
	eyear := c.DefaultQuery("eyear", yy)
	exchange := c.DefaultQuery("exchange", "")
	// timeoutSeconds := 10

	url := fmt.Sprintf("%s/get_dates_by_year", gmapi)
	pars := map[string]string{
		"syear": syear,
		"eyear": eyear,
	}
	if exchange != "" {
		pars["exchange"] = exchange
	}
	rawData, err := GetURLWithoutRetry(url, pars, 10*time.Second, 0)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}

	jsonBytes, marshalErr := json.Marshal(rawData)
	if marshalErr != nil {
		c.JSON(http.StatusNotFound, marshalErr)
		return
	}

	var rcd RawColData
	unmarshalErr := json.Unmarshal(jsonBytes, &rcd)
	if unmarshalErr != nil {
		c.JSON(http.StatusNotFound, unmarshalErr)
		return
	}

	// 处理获取到的 JSON 数据为 records 形式
	records, transformErr := rcd.TransformToRecords()
	if transformErr != nil {
		c.JSON(http.StatusNotFound, transformErr)
		return
	}

	// fmt.Println("\n--- 转换后的 records 格式 (getURLWithRetry) ---")
	// recordsJSON, _ := json.MarshalIndent(records[:5], "", "  ") // 格式化输出 JSON
	// fmt.Printf("%s\n", recordsJSON)
	c.JSON(http.StatusOK, records)

}

func RouteCalendar(c *gin.Context) {
	yy := fmt.Sprintf("%d", time.Now().Year())
	syear := c.DefaultQuery("syear", yy)
	eyear := c.DefaultQuery("eyear", yy)
	exchange := c.DefaultQuery("exchange", "")
	timeoutSeconds := 10

	rawData, err := gm.GetCalendar(gmapi, syear, eyear, exchange, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.Calendar)": err.Error()})
		return
	}

	var rcd RawColData
	unmarshalErr := json.Unmarshal(rawData, &rcd)
	if unmarshalErr != nil {
		c.JSON(http.StatusNotFound, unmarshalErr)
		return
	}

	// 处理获取到的 JSON 数据为 records 形式
	records, transformErr := rcd.TransformToRecords()
	if transformErr != nil {
		c.JSON(http.StatusNotFound, gin.H{" Err(TransformToRecords)": transformErr.Error()})
		return
	}

	c.JSON(http.StatusOK, records)
}

func RouteDatesPrevN(c *gin.Context) {
	today := time.Now().Format("2006-01-02")
	date := c.DefaultQuery("date", "")
	if date == "" {
		date = today
	}
	scount := c.DefaultQuery("count", "10")
	count, err := strconv.Atoi(scount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{" Err(count)": err.Error()})
		return
	}

	include := c.DefaultQuery("include", "true")
	isinclude := true
	if include == "false" {
		isinclude = false
	}

	timeoutSeconds := 10
	// rawData, err := gm.GetPrevNByte(gmapi, date, count, timeoutSeconds, isinclude)
	rawData, err := gm.GetPrevN(gmapi, date, count, timeoutSeconds, isinclude)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetPrevN)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteDatesNextN(c *gin.Context) {
	today := time.Now().Format("2006-01-02")
	date := c.DefaultQuery("date", "")

	if date == "" {
		date = today
	}
	scount := c.DefaultQuery("count", "10")
	count, err := strconv.Atoi(scount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{" Err(count)": err.Error()})
		return
	}

	include := c.DefaultQuery("include", "true")
	isinclude := true
	if include == "false" {
		isinclude = false
	}

	timeoutSeconds := 10
	// rawData, err := gm.GetNextNByte(gmapi, date, count, timeoutSeconds, isinclude)
	rawData, err := gm.GetNextN(gmapi, date, count, timeoutSeconds, isinclude)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetNextN)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteCurrent(c *gin.Context) {
	symbols := c.DefaultQuery("symbols", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbols 参数为必须"))
		return
	}

	split := c.DefaultQuery("split", "false")
	issplit := false
	if split == "true" {
		issplit = true
	}

	timeoutSeconds := 10

	// rawData, err := gm.GetCurrentByte(gmapi, symbols, timeoutSeconds, issplit)
	rawData, err := gm.GetCurrent(gmapi, symbols, timeoutSeconds, issplit)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetCurrent)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)

	// // 将获取到的字符串数据解析为 JSON 格式
	// var data any
	// if err = json.Unmarshal(rawData, &data); err != nil {
	// 	c.JSON(http.StatusNotAcceptable, gin.H{" Err(unmarshaling JSON)": err.Error()})
	// 	return
	// }
	// c.JSON(http.StatusOK, data)
	// c.JSON(http.StatusOK, records)
}

func RouteKbars(c *gin.Context) {
	symbols := c.DefaultQuery("symbols", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbols 参数为必须"))
		return
	}

	today := time.Now().Format("2006-01-02")
	sdate := c.DefaultQuery("sdate", today)
	edate := c.DefaultQuery("edate", today)
	tag := c.DefaultQuery("tag", "1m")
	timestamp := c.DefaultQuery("timestamp", "false")
	timeoutSeconds := 10

	istimestamp := false
	if timestamp == "true" {
		istimestamp = true
	}
	// rawData, err := gm.GetKbarsHisByte(gmapi, symbols, tag, sdate, edate, timeoutSeconds)
	rawData, err := gm.GetKbarsHis(gmapi, symbols, tag, sdate, edate, istimestamp, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetKbarsHis)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)

	// var rcd RawColData
	// if unmarshalErr := json.Unmarshal(rawData, &rcd); unmarshalErr != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{" Err(unmarshalErr)": unmarshalErr.Error()})
	// 	return
	// }

	// records, transformErr := rcd.TransformToRecords()
	// if transformErr != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{" Err(transformErr)": transformErr.Error()})
	// 	return
	// }

	// c.JSON(http.StatusOK, records)
}

// 获取股票当日的 K 线数据，返回字典json格式，键为代码，值为K线数据
func RouteKBDict(c *gin.Context) {
	symbols := c.DefaultQuery("symbols", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbols 参数为必须"))
		return
	}

	today := time.Now().Format("2006-01-02")
	sdate := c.DefaultQuery("sdate", today)
	edate := c.DefaultQuery("edate", today)
	tag := c.DefaultQuery("tag", "1m")
	timestamp := c.DefaultQuery("timestamp", "false")
	timeoutSeconds := 10

	istimestamp := false
	if timestamp == "true" {
		istimestamp = true
	}

	rawData, err := gm.GetKbarsHis(gmapi, symbols, tag, sdate, edate, istimestamp, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetKbarsHis)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gm.ConvertRecords2Dict(rawData))
}

// 获取股票当日的 K 线数据，返回字典json格式，键为代码，值为K线数据
func RouteKBDictTS(c *gin.Context) {
	symbols := c.DefaultQuery("symbols", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbols 参数为必须"))
		return
	}

	today := time.Now().Format("2006-01-02")
	sdate := c.DefaultQuery("sdate", today)
	edate := c.DefaultQuery("edate", today)
	tag := c.DefaultQuery("tag", "1m")
	timestamp := c.DefaultQuery("timestamp", "false")
	timeoutSeconds := 10

	istimestamp := false
	if timestamp == "true" {
		istimestamp = true
	}

	rawData, err := gm.GetKbarsHis(gmapi, symbols, tag, sdate, edate, istimestamp, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetKbarsHis)": err.Error()})
		return
	}
	if istimestamp {
		c.JSON(http.StatusOK, gm.ConvertRecords2DictTSInt(rawData))
	} else {
		c.JSON(http.StatusOK, gm.ConvertRecords2DictTSString(rawData))
	}
	// c.JSON(http.StatusOK, gm.ConvertRecords2DictTSString(rawData))
}

func RouteKbarsN(c *gin.Context) {
	symbol := c.DefaultQuery("symbol", "")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbols 参数为必须"))
		return
	}

	today := time.Now().Format("2006-01-02")
	edate := c.DefaultQuery("edate", today)
	count := c.DefaultQuery("count", "")
	tag := c.DefaultQuery("tag", "1d")
	timestamp := c.DefaultQuery("timestamp", "false")
	timeoutSeconds := 10

	istimestamp := false
	if timestamp == "true" {
		istimestamp = true
	}

	// rawData, err := gm.GetKbarsHisNByte(gmapi, symbol, tag, count, edate, timeoutSeconds)
	rawData, err := gm.GetKbarsHisN(gmapi, symbol, tag, count, edate, istimestamp, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetKbarsHisN)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteCSVxzMonth(c *gin.Context) {
	symbol := c.DefaultQuery("symbol", "")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbol 参数为必须参数"))
		return
	}

	timeoutSeconds := 10
	now := time.Now()
	year, month, _ := now.Date()
	yearStr := fmt.Sprintf("%d", year)
	monthStr := fmt.Sprintf("%02d", month)

	cyear := c.DefaultQuery("year", yearStr)
	cmonth := c.DefaultQuery("month", monthStr)
	// tag := c.DefaultQuery("tag", "1m")
	imonth, _ := strconv.Atoi(cmonth)
	iyear, _ := strconv.Atoi(cyear)
	timestamp := c.DefaultQuery("timestamp", "false")
	istimestamp := false
	if timestamp == "true" {
		istimestamp = true
	}

	// rawData, err := gm.GetCSVMonthJson(gmcsv, symbol, imonth, iyear, timeoutSeconds)
	rawData, err := gm.GetCSVMonth(gmcsv, symbol, imonth, iyear, istimestamp, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetCSVMonth)": err.Error()})
		return
	}
	// c.JSON(http.StatusOK, string(rawData))
	c.JSON(http.StatusOK, rawData)

	// 将获取到的字符串数据解析为 JSON 格式
	// var data any
	// if err = json.Unmarshal(rawData, &data); err != nil {
	// 	c.JSON(http.StatusNotAcceptable, gin.H{" Err(unmarshaling JSON)": err.Error()})
	// 	return
	// }
	// c.JSON(http.StatusOK, data)
}

func RouteCSVxzYear(c *gin.Context) {
	symbol := c.DefaultQuery("symbol", "")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbol 参数为必须参数"))
		return
	}

	timeoutSeconds := 10
	now := time.Now()
	year := now.Year()
	yearStr := fmt.Sprintf("%d", year)

	cyear := c.DefaultQuery("year", yearStr)
	tag := c.DefaultQuery("tag", "1m")
	timestamp := c.DefaultQuery("timestamp", "false")
	istimestamp := false
	if timestamp == "true" {
		istimestamp = true
	}
	iyear, _ := strconv.Atoi(cyear)

	lookuptab := map[string]string{
		"1m": "timestamp",
		"vv": "timestamp",
		"pe": "trade_date",
	}

	rawData, err := gm.GetCSVYear(gmcsv, symbol, tag, iyear, istimestamp, lookuptab[tag], timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetCSVYear)": err.Error()})
		return
	}
	// c.JSON(http.StatusOK, string(rawData))
	c.JSON(http.StatusOK, rawData)

	// 将获取到的字符串数据解析为 JSON 格式
	// var data any
	// if err = json.Unmarshal(rawData, &data); err != nil {
	// 	c.JSON(http.StatusNotAcceptable, gin.H{" Err(unmarshaling JSON)": err.Error()})
	// 	return
	// }
	// c.JSON(http.StatusOK, data)
}

func RouteCSVxz1m(c *gin.Context) {
	symbol := c.DefaultQuery("symbol", "")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbol 参数为必须参数"))
		return
	}

	timeoutSeconds := 10
	now := time.Now()
	today := now.Format("2006-01-02")

	sdate := c.DefaultQuery("sdate", today)
	edate := c.DefaultQuery("edate", today)
	// tag := c.DefaultQuery("tag", "1m")

	timestamp := c.DefaultQuery("timestamp", "false")
	istimestamp := false
	if timestamp == "true" {
		istimestamp = true
	}

	clip := c.DefaultQuery("clip", "true")
	isclip := false
	if clip == "true" {
		isclip = true
	}

	rawData, err := gm.GetCSV1m(gmcsv, symbol, sdate, edate, istimestamp, isclip, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetCSV1m)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)

	// // jsonData, err := json.MarshalIndent(result, "", "  ")
	// jsonData, err := json.Marshal(rawData)
	// if err != nil {
	// 	c.JSON(http.StatusNotAcceptable, gin.H{" Err(unmarshaling JSON)": err.Error()})
	// }

	// c.JSON(http.StatusOK, string(jsonData))
}

func RouteCSVxzTag(c *gin.Context) {
	symbol := c.DefaultQuery("symbol", "")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbol 参数为必须参数"))
		return
	}

	timeoutSeconds := 10
	now := time.Now()
	today := now.Format("2006-01-02")

	sdate := c.DefaultQuery("sdate", today)
	edate := c.DefaultQuery("edate", today)
	tag := c.DefaultQuery("tag", "vv")

	timestamp := c.DefaultQuery("timestamp", "false")
	istimestamp := false
	if timestamp == "true" {
		istimestamp = true
	}

	clip := c.DefaultQuery("clip", "true")
	isclip := false
	if clip == "true" {
		isclip = true
	}

	rawData, err := gm.GetCSVTag(gmcsv, tag, symbol, sdate, edate, istimestamp, isclip, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetCSVTag)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteGM1m(c *gin.Context) {
	symbol := c.DefaultQuery("symbol", "")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbol 参数为必须参数"))
		return
	}

	timeoutSeconds := 10
	now := time.Now()
	today := now.Format("2006-01-02")
	prday := now.AddDate(0, 0, -1)
	yesterday := prday.Format("2006-01-02")

	// tag := c.DefaultQuery("tag", "1m")

	timestamp := c.DefaultQuery("timestamp", "false")
	istimestamp := false
	if timestamp == "true" {
		istimestamp = true
	}
	include := c.DefaultQuery("include", "true")
	isinclude := true
	if include == "false" {
		isinclude = false
	}
	cday := yesterday
	if isinclude && gm.IsAOpen() {
		cday = today
	}
	sdate := c.DefaultQuery("sdate", cday)
	edate := c.DefaultQuery("edate", cday)

	rawData, err := gm.GetGM1m(gmcsv, gmapi, symbol, sdate, edate, istimestamp, isinclude, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetGM1m)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}
