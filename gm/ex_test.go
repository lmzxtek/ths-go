package gm

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/go-gota/gota/dataframe"
)

// var gmURL = "http://45.154.14.186:5000" // locVPS-kr
var gmURL = "http://localhost:5000"
var gmCSV = "http://localhost:5002"

// TestReadCSV is a test function to read CSV file and print the dataframe
func TestCSV(t *testing.T) {
	jsonStr := `{
		"columns":["Symbol","Time","Price","Volume"],
		"data":[
			["AAPL","2025-03-01",100,200],
			["AAPL","2025-03-01",101,210],
			["AAPL","2025-03-01",102,220],
			["AAPL","2025-03-01",103,230],
			["AAPL","2025-03-01",104,240]]
	}`
	df, err := ParseJsonToDataframe(jsonStr)
	if err != nil {
		fmt.Printf("获取数据失败: %s\n", err)
	}
	fmt.Println(df)
}

// TestReadCSV is a test function to read CSV file and print the dataframe
func TestReadCSV(t *testing.T) {
	fmt.Println("\n >>> Start read dataframe... ")

	csvStr := `
Country,Date,Age,Amount,Id
"United States",2012-02-01,50,112.1,01234
"United States",2012-02-01,32,321.31,54320
"United Kingdom",2012-02-01,17,18.2,12345
"United States",2012-02-01,32,321.31,54320
"United Kingdom",2012-02-01,NA,18.2,12345
"United States",2012-02-01,32,321.31,54320
"United States",2012-02-01,32,321.31,54320
Spain,2012-02-01,66,555.42,00241
`
	df := dataframe.ReadCSV(strings.NewReader(csvStr))
	// fmt.Println(df.Col("Country"))
	fmt.Println(df.Types())
	fmt.Println(df.Describe())
	fmt.Println(df)
}

// TestReadJSON is a test function to read JSON file and print the dataframe
func TestReadJSON(t *testing.T) {
	jsonStr := `[{"COL.2":1,"COL.3":3},{"COL.1":5,"COL.2":2,"COL.3":2},{"COL.1":6,"COL.2":3,"COL.3":1}]`
	df := dataframe.ReadJSON(strings.NewReader(jsonStr))
	fmt.Println(df)

}

// TestReadJSON is a test function to read JSON file and print the dataframe
func TestParseJsonToDataframe(t *testing.T) {
	jsonStr := `{"columns":["Symbol","Time","Price","Volume"],"data":[["AAPL","2025-05-01",100,200],["AAPL","2025-05-01",101,210],["AAPL","2025-05-01",102,220],["AAPL","2025-05-01",103,230],["AAPL","2025-05-01",104,240]]}`
	df, err := ParseJsonToDataframe(jsonStr)
	if err != nil {
		fmt.Printf("获取数据失败: %s\n", err)
	}
	fmt.Println(df)
}

func TestTimestamp(t *testing.T) {
	// 测试时间字符串
	testTimeStr := "2024-01-13 10:30:00"

	fmt.Printf("原始时间字符串: %s\n", testTimeStr)
	fmt.Println("=" + fmt.Sprintf("%*s", 50, "="))

	// 1. 转换为Unix时间戳（秒）- UTC时区
	timestamp, err := TimeStringToTimestamp(testTimeStr)
	if err != nil {
		fmt.Printf("转换失败: %v\n", err)
	} else {
		fmt.Printf("Unix时间戳（秒，UTC）: %d\n", timestamp)
		// 验证转换
		backToString := TimestampToTimeString(timestamp)
		fmt.Printf("转换回字符串: %s\n", backToString)
	}

	fmt.Println()

	// 2. 转换为Unix时间戳（秒）- 中国时区
	timestampChina, err := TimeStringToTimestampWithLocation(testTimeStr, "Asia/Shanghai")
	if err != nil {
		fmt.Printf("转换失败: %v\n", err)
	} else {
		fmt.Printf("Unix时间戳（秒，中国时区）: %d\n", timestampChina)
		// 验证转换
		backToStringChina, _ := TimestampToTimeStringWithLocation(timestampChina, "Asia/Shanghai")
		fmt.Printf("转换回字符串（中国时区）: %s\n", backToStringChina)
	}

	fmt.Println()

	// 3. 转换为毫秒时间戳
	timestampMillis, err := TimeStringToTimestampMillis(testTimeStr)
	if err != nil {
		fmt.Printf("转换失败: %v\n", err)
	} else {
		fmt.Printf("Unix时间戳（毫秒）: %d\n", timestampMillis)
	}

	// 4. 转换为纳秒时间戳
	timestampNano, err := TimeStringToTimestampNano(testTimeStr)
	if err != nil {
		fmt.Printf("转换失败: %v\n", err)
	} else {
		fmt.Printf("Unix时间戳（纳秒）: %d\n", timestampNano)
	}

	fmt.Println()
	fmt.Println("时间戳差异说明:")
	fmt.Printf("UTC与中国时区时间戳差异: %d 秒（8小时 = %d秒）\n",
		timestampChina-timestamp, 8*3600)

	// 测试多个时间格式
	fmt.Println()
	fmt.Println("测试其他时间:")
	testTimes := []string{
		"2024-01-01 00:00:00",
		"2024-12-31 23:59:59",
		"2024-06-15 12:00:00",
	}

	for _, timeStr := range testTimes {
		if ts, err := TimeStringToTimestampWithLocation(timeStr, "Asia/Shanghai"); err == nil {
			fmt.Printf("%s -> %d\n", timeStr, ts)
		}
	}
}

func TestTimestamp2(tt *testing.T) {
	// 测试的毫秒时间戳
	tint := int64(1136214245000)

	fmt.Printf("原始毫秒时间戳: %d\n", tint)
	fmt.Println("=" + fmt.Sprintf("%*s", 60, "="))

	// 1. 基本转换
	t := MillisToTime(tint)
	fmt.Printf("转换为time.Time: %v\n", t)
	fmt.Printf("UTC时间: %s\n", t.UTC().Format("2006-01-02 15:04:05"))
	fmt.Printf("本地时间: %s\n", t.Local().Format("2006-01-02 15:04:05"))

	fmt.Println()

	// 2. 转换为不同时区
	locations := []string{"UTC", "Asia/Shanghai", "America/New_York", "Europe/London"}

	fmt.Println("不同时区的时间表示:")
	for _, loc := range locations {
		if timeInLoc, err := MillisToTimeInLocation(tint, loc); err == nil {
			fmt.Printf("%-20s: %s\n", loc, timeInLoc.Format("2006-01-02 15:04:05 MST"))
		} else {
			fmt.Printf("%-20s: 转换失败 - %v\n", loc, err)
		}
	}

	fmt.Println()

	// 3. 不同格式的字符串输出
	fmt.Println("不同格式的时间字符串:")
	formats := map[string]string{
		"标准格式":    "2006-01-02 15:04:05",
		"日期格式":    "2006-01-02",
		"时间格式":    "15:04:05",
		"RFC3339": time.RFC3339,
		"自定义格式1":  "2006年01月02日 15时04分05秒",
		"自定义格式2":  "Jan 02, 2006 3:04:05 PM",
	}

	for name, layout := range formats {
		formatted := MillisToFormattedString(tint, layout)
		fmt.Printf("%-20s: %s\n", name, formatted)
	}

	fmt.Println()

	// 4. 显示详细信息
	fmt.Println("时间戳详细信息:")
	info := GetTimestampInfo(tint)
	for key, value := range info {
		fmt.Printf("%-18s: %v\n", key, value)
	}

	fmt.Println()

	// 5. 验证转换（往返转换）
	fmt.Println("验证转换正确性:")
	convertedBack := TimeToMillis(t)
	fmt.Printf("原始毫秒时间戳: %d\n", tint)
	fmt.Printf("转换后再转回:   %d\n", convertedBack)
	fmt.Printf("转换是否正确:   %t\n", tint == convertedBack)

	fmt.Println()

	// 6. 测试其他毫秒时间戳
	fmt.Println("测试其他时间戳:")
	testTimestamps := []int64{
		0,                      // Unix纪元开始
		1000000000000,          // 2001-09-09
		1609459200000,          // 2021-01-01 00:00:00 UTC
		time.Now().UnixMilli(), // 当前时间
	}

	for _, ts := range testTimestamps {
		t := MillisToTime(ts)
		fmt.Printf("%13d -> %s\n", ts, t.Format("2006-01-02 15:04:05"))
	}
}

func TestMarketOpenTime(t *testing.T) {
	// 测试当前时间
	isOpen := IsChineseStockMarketOpen()
	fmt.Printf("当前中国股市是否开市: %t\n", isOpen)

	// 获取当前中国时间
	chinaLocation, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(chinaLocation)
	fmt.Printf("当前中国时间: %s\n", now.Format("2006-01-02 15:04:05 Monday"))

	// 获取下一个交易时间
	nextTradingTime := GetNextTradingTime()
	fmt.Printf("下一个交易时间: %s\n", nextTradingTime.Format("2006-01-02 15:04:05 Monday"))

	// 测试一些特定时间
	testTimes := []string{
		"2024-01-15 09:25:00", // 开市前
		"2024-01-15 10:30:00", // 上午交易时间
		"2024-01-15 12:00:00", // 午休时间
		"2024-01-15 14:30:00", // 下午交易时间
		"2024-01-15 16:00:00", // 收市后
		"2024-01-13 10:30:00", // 周六
	}

	fmt.Println("\n测试特定时间:")
	for _, timeStr := range testTimes {
		testTime, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr, chinaLocation)
		isTestOpen := IsChineseStockMarketOpenAt(testTime)
		fmt.Printf("%s (%s): %t\n",
			testTime.Format("2006-01-02 15:04:05"),
			testTime.Weekday().String()[:3],
			isTestOpen)
	}
}

func TestTickRecord(t *testing.T) {
	// 创建示例分时数据
	now := time.Now()
	records := []TickRecord{
		{Timestamp: now.Add(-2 * 24 * time.Hour), Price: 100.0, Volume: 1000},
		{Timestamp: now.Add(-2*24*time.Hour + time.Hour), Price: 101.5, Volume: 1500},
		{Timestamp: now.Add(-2*24*time.Hour + 2*time.Hour), Price: 99.8, Volume: 800},
		{Timestamp: now.Add(-2*24*time.Hour + 3*time.Hour), Price: 102.0, Volume: 1200},

		{Timestamp: now.Add(-24 * time.Hour), Price: 102.5, Volume: 900},
		{Timestamp: now.Add(-24*time.Hour + time.Hour), Price: 103.0, Volume: 1100},
		{Timestamp: now.Add(-24*time.Hour + 2*time.Hour), Price: 101.0, Volume: 700},
		{Timestamp: now.Add(-24*time.Hour + 3*time.Hour), Price: 104.0, Volume: 1300},

		{Timestamp: now, Price: 104.5, Volume: 1000},
		{Timestamp: now.Add(time.Hour), Price: 105.0, Volume: 1400},
		{Timestamp: now.Add(2 * time.Hour), Price: 103.5, Volume: 600},
		{Timestamp: now.Add(3 * time.Hour), Price: 106.0, Volume: 1600},
	}

	// 转换为日K线
	dailyKLines, err := ConvertTicksToDailyWithValidation(records)
	if err != nil {
		fmt.Printf("转换失败: %v\n", err)
		return
	}

	// 输出结果
	fmt.Println("日K线数据:")
	fmt.Printf("%-12s %-8s %-8s %-8s %-8s %-8s\n", "日期", "开盘", "最高", "最低", "收盘", "成交量")
	fmt.Println(strings.Repeat("-", 60))

	for _, kline := range dailyKLines {
		fmt.Printf("%-12s %-8.2f %-8.2f %-8.2f %-8.2f %-8d\n",
			kline.Date.Format("2006-01-02"),
			kline.Open,
			kline.High,
			kline.Low,
			kline.Close,
			kline.Volume,
		)
	}
}

func TestKBarType(t *testing.T) {
	fmt.Println(" -=> Start test KBar type ... ")

	kb := KBar{
		Open:   100.0,
		High:   110.0,
		Low:    95.0,
		Close:  105.0,
		Volume: 1000,
	}
	fmt.Println(kb)

	now := time.Now()
	fmt.Println(now.Format("2006-01-02 15:04:05 Monday"))

	stringTime := "2006-01-02 15:04:05"
	fmt.Println(stringTime)

	tt, err := time.Parse("2006-01-02 15:04:05", stringTime)
	if err != nil {
		fmt.Printf("解析时间失败: %v", err)
	}
	fmt.Println(tt)
	fmt.Println(tt.Local())

	fmt.Println(kb.ToList(now.Local().Format("2006-01-02T15:04:05")))
	fmt.Println(kb.ToList(now.Unix()))
	fmt.Println(kb.ToList(tt.UnixMilli()))
	fmt.Println(kb.ToRecords(now.UnixMilli()))
}

func TestKBData(t *testing.T) {
	// 示例1: 单个K线数据（字符串时间戳）
	jsonStr1 := `{
		"timestamp": "2025-06-11T09:30:00+08:00",
		"open": 100.5,
		"high": 102.0,
		"low": 99.8,
		"close": 101.2,
		"volume": 50000
	}`

	kbar1, err := ParseKBarFromJSON(jsonStr1)
	if err != nil {
		fmt.Printf("解析失败: %v\n", err)
	} else {
		fmt.Printf("K线1: %+v\n", kbar1)
	}

	// 示例2: K线数组（Unix时间戳）
	jsonStr2 := `[
		{
			"timestamp": 1704067800,
			"open": 100.5,
			"high": 102.0,
			"low": 99.8,
			"close": 101.2,
			"volume": 50000
		},
		{
			"timestamp": 1704154200,
			"open": 101.2,
			"high": 103.5,
			"low": 100.9,
			"close": 102.8,
			"volume": 60000
		}
	]`

	kbars2, err := ParseKBarArrayFromJSON(jsonStr2)
	if err != nil {
		fmt.Printf("解析数组失败: %v\n", err)
	} else {
		fmt.Println("\nK线数组:")
		for i, kbar := range kbars2 {
			fmt.Printf("K线%d: 时间=%s, 开盘=%.2f, 最高=%.2f, 最低=%.2f, 收盘=%.2f, 成交量=%d\n",
				i+1, kbar.Timestamp.Format("2006-01-02 15:04:05"),
				kbar.Open, kbar.High, kbar.Low, kbar.Close, kbar.Volume)
		}
	}

	// 示例3: API响应格式
	apiJsonStr := `{
		"code": 0,
		"message": "success",
		"data": [
			{
				"timestamp": "2024-01-01 09:30:00",
				"open": 100.5,
				"high": 102.0,
				"low": 99.8,
				"close": 101.2,
				"volume": 50000
			}
		]
	}`

	kbars3, err := ParseKBarFromAPIResponse(apiJsonStr)
	if err != nil {
		fmt.Printf("解析API响应失败: %v\n", err)
	} else {
		fmt.Println("\nAPI响应K线数据:")
		for _, kbar := range kbars3 {
			fmt.Printf("时间=%s, OHLCV=(%.2f,%.2f,%.2f,%.2f,%d)\n",
				kbar.Timestamp.Format("2006-01-02 15:04:05"),
				kbar.Open, kbar.High, kbar.Low, kbar.Close, kbar.Volume)
		}
	}

	// 示例4: 毫秒时间戳
	jsonStr4 := `{
		"timestamp": 1704067800000,
		"open": 100.5,
		"high": 102.0,
		"low": 99.8,
		"close": 101.2,
		"volume": 50000
	}`

	kbar4, err := ParseKBarFromJSON(jsonStr4)
	if err != nil {
		fmt.Printf("解析毫秒时间戳失败: %v\n", err)
	} else {
		fmt.Printf("\n毫秒时间戳K线: 时间=%s\n", kbar4.Timestamp.Format("2006-01-02 15:04:05"))
	}

	// 批量处理示例
	fmt.Println("\n=== 批量处理示例 ===")
	jsonStrings := []string{
		`{"timestamp": "2024-01-01T09:30:00", "open": 100, "high": 101, "low": 99, "close": 100.5, "volume": 1000}`,
		`{"timestamp": 1704067800, "open": 100.5, "high": 102, "low": 99.5, "close": 101, "volume": 1500}`,
		`{"timestamp": "2024-01-01 10:30:00", "open": 101, "high": 103, "low": 100, "close": 102, "volume": 2000}`,
	}

	for i, jsonStr := range jsonStrings {
		if kbar, err := ParseKBarFromJSON(jsonStr); err != nil {
			fmt.Printf("批量处理%d失败: %v\n", i+1, err)
		} else {
			fmt.Printf("批量%d: %s -> OHLC(%.1f,%.1f,%.1f,%.1f)\n",
				i+1, kbar.Timestamp.Format("15:04:05"),
				kbar.Open, kbar.High, kbar.Low, kbar.Close)
		}
	}
}

func TestOHLCV(t *testing.T) {
	// 创建示例数据
	now := time.Now()
	kbars := []KBarData{
		{
			Timestamp: now.Add(-2 * time.Hour),
			Open:      100.5,
			High:      102.0,
			Low:       99.8,
			Close:     101.2,
			Volume:    50000,
		},
		{
			Timestamp: now.Add(-1 * time.Hour),
			Open:      101.2,
			High:      103.5,
			Low:       100.9,
			Close:     102.8,
			Volume:    60000,
		},
		{
			Timestamp: now,
			Open:      102.8,
			High:      104.0,
			Low:       102.0,
			Close:     103.5,
			Volume:    70000,
		},
	}

	fmt.Println("=== 原始数据 ===")
	for i, kbar := range kbars {
		fmt.Printf("K线%d: %s -> OHLCV(%.2f,%.2f,%.2f,%.2f,%d)\n",
			i+1, kbar.Timestamp.Format("15:04:05"),
			kbar.Open, kbar.High, kbar.Low, kbar.Close, kbar.Volume)
	}

	// 方法1: RFC3339格式时间戳为键
	fmt.Println("\n=== 方法1: RFC3339格式时间戳为键 ===")
	json1, err := ConvertToTimestampKeyedJSON(kbars)
	if err != nil {
		fmt.Printf("转换失败: %v\n", err)
	} else {
		fmt.Println(json1)
	}

	// 方法2: Unix时间戳为键
	fmt.Println("\n=== 方法2: Unix时间戳为键 ===")
	json2, err := ConvertToUnixTimestampKeyedJSON(kbars)
	if err != nil {
		fmt.Printf("转换失败: %v\n", err)
	} else {
		fmt.Println(json2)
	}

	// 方法3: 自定义格式时间戳为键
	fmt.Println("\n=== 方法3: 自定义格式时间戳为键 ===")
	json3, err := ConvertToCustomTimestampKeyedJSON(kbars, "2006-01-02 15:04:05")
	if err != nil {
		fmt.Printf("转换失败: %v\n", err)
	} else {
		fmt.Println(json3)
	}

	// 方法4: 毫秒时间戳为键
	fmt.Println("\n=== 方法4: 毫秒时间戳为键 ===")
	json4, err := ConvertToMillisTimestampKeyedJSON(kbars)
	if err != nil {
		fmt.Printf("转换失败: %v\n", err)
	} else {
		fmt.Println(json4)
	}

	// 方法5: 紧凑格式
	fmt.Println("\n=== 方法5: 紧凑格式（数组形式） ===")
	json5, err := ConvertToCompactTimestampKeyedJSON(kbars)
	if err != nil {
		fmt.Printf("转换失败: %v\n", err)
	} else {
		fmt.Println(json5)
	}

	// 方法6: 嵌套结构（按日期分组）
	fmt.Println("\n=== 方法6: 按日期分组 ===")
	json6, err := ConvertToNestedTimestampKeyedJSON(kbars)
	if err != nil {
		fmt.Printf("转换失败: %v\n", err)
	} else {
		fmt.Println(json6)
	}

	// 通用转换函数示例
	fmt.Println("\n=== 通用转换函数示例 ===")

	// 使用Unix时间戳，不美化
	options1 := ConvertOptions{
		UseUnixTime: true,
		PrettyPrint: false,
	}
	compactJson, err := ConvertKBarToTimestampJSON(kbars, options1)
	if err != nil {
		fmt.Printf("转换失败: %v\n", err)
	} else {
		fmt.Printf("紧凑Unix时间戳: %s\n", compactJson)
	}

	// 使用自定义格式，美化输出
	options2 := ConvertOptions{
		TimeFormat:  "01-02 15:04",
		PrettyPrint: true,
	}
	prettyJson, err := ConvertKBarToTimestampJSON(kbars, options2)
	if err != nil {
		fmt.Printf("转换失败: %v\n", err)
	} else {
		fmt.Printf("自定义格式:\n%s\n", prettyJson)
	}

	// 演示从JSON中获取特定时间戳的数据
	fmt.Println("\n=== 从JSON中获取特定数据 ===")
	if len(kbars) > 0 {
		targetTime := kbars[0].Timestamp.Format(time.RFC3339)
		if ohlcv, err := GetKBarByTimestamp(json1, targetTime); err != nil {
			fmt.Printf("获取数据失败: %v\n", err)
		} else {
			fmt.Printf("时间戳 %s 对应的数据: %+v\n", targetTime, ohlcv)
		}
	}
}

func TestJsonDF(t *testing.T) {
	jsonData := []byte(`{
			"columns": ["Symbol","Time","Price","Volume"],
			"data": [["AAPL","2025-05-01",100,200],["AAPL","2025-05-01",101,210]]
		}`)

	// 解析原始数据
	var raw RawData
	if err := json.Unmarshal(jsonData, &raw); err != nil {
		panic(err)
	}

	fmt.Println("表头:\n", raw.Columns)
	fmt.Println("数据:\n", raw.Data)

	// df := dataframe.LoadRecords(raw.Data)
	// df := dataframe.LoadRecords(
	// 	[][]string{
	// 		{"A", "B", "C", "D"},
	// 		{"a", "4", "5.1", "true"},
	// 		{"k", "5", "7.0", "true"},
	// 		{"k", "4", "6.0", "true"},
	// 		{"a", "2", "7.1", "false"},
	// 	},
	// )
	// fmt.Println(df)

	// 转换为列式存储
	df, err := ConvertToColumnar(raw)
	if err != nil {
		panic(err)
	}

	// 示例输出
	fmt.Println("Symbol列:", df["Symbol"])
	fmt.Println("Price列 :", df["Price"])
}

// TestloadRecords is a test function to load records and print the dataframe
// 测试Dataframe的LoadRecords方法
func TestLoadRecords(t *testing.T) {
	rec := [][]string{
		{"A", "B", "C", "D"},
		{"a", "4", "5.1", "true"},
		{"k", "5", "7.0", "true"},
		{"k", "4", "6.0", "true"},
		{"a", "2", "7.1", "false"},
	}
	fmt.Println(rec)

	df := dataframe.LoadRecords(rec)
	fmt.Println(df)
}

// 测试Dataframe的LoadStructs方法
// 测试结构体
func TestLoadStructs(t *testing.T) {
	type User struct {
		Name     string
		Age      int
		Accuracy float64
	}
	users := []User{
		{"Aram", 17, 0.2},
		{"Juan", 18, 0.8},
		{"Ana", 22, 0.5},
	}
	df := dataframe.LoadStructs(users)
	fmt.Println(df)
}

func TestDFGetTest(t *testing.T) {
	fmt.Println(" -=> Start fetch df from url(Test) ... ")
	df, err := DfGetTest(gmURL)
	if err != nil {
		fmt.Printf("获取数据失败: %s\n", err)
	}
	fmt.Println(df)
	// cxz.SaveDataframeToCSV(&df, "data.csv")
	// cxz.SaveDataframeToCSVxz(&df, "data.csv.xz")
}

func TestDFGetTest2(t *testing.T) {
	fmt.Println(" -=> Start fetch df from url(Test2) ... ")
	df, err := DfGetTest2(gmURL)
	if err != nil {
		fmt.Printf("获取数据失败: %s\n", err)
	}
	fmt.Println(df)
	// cxz.SaveDataframeToCSV(&df, "data.csv")
	// cxz.SaveDataframeToCSVxz(&df, "data.csv.xz")
}

func TestGetCalendar(t *testing.T) {
	fmt.Println(" -=> Start fetch calendar using GM-api ... ")

	timeoutSeconds := 10
	url := gmURL
	syear := "2025"
	eyear := "2025"
	exchange := "" // "SHSE"

	resp, err := GetCalendar(url, syear, eyear, exchange, timeoutSeconds)
	if err != nil {
		fmt.Printf(" 获取数据失败(gm.GetCalendar): %v\n", err)
	}
	fmt.Println(string(resp))
}

func TestGetPrevN(t *testing.T) {
	fmt.Println(" -=> Start fetch prev_n GM-api ... ")

	timeoutSeconds := 10
	url := gmURL

	date := "2025-05-01"
	count := 5
	include := true

	resp, err := GetPrevNByte(url, date, count, timeoutSeconds, include)
	if err != nil {
		fmt.Printf(" 获取数据失败(gm.GetPrevN): %v\n", err)
	}
	fmt.Println(string(resp))
}
func TestGetNextN(t *testing.T) {
	fmt.Println(" -=> Start fetch next_n GM-api ... ")

	timeoutSeconds := 10
	url := gmURL

	date := "2025-05-11"
	count := 5
	include := true

	resp, err := GetNextNByte(url, date, count, timeoutSeconds, include)
	if err != nil {
		fmt.Printf(" 获取数据失败(gm.GetNextN): %v\n", err)
	}
	fmt.Println(string(resp))
}
func TestGetCurrent(t *testing.T) {
	fmt.Println(" -=> Start fetch current snap data using GM-api ... ")

	url := gmURL
	timeoutSeconds := 10

	symbols := "SHSE.601088,SZSE.300917"
	// resp, err := GetCurrent(url, symbols, timeoutSeconds, false)
	resp, err := GetCurrentByte(url, symbols, timeoutSeconds, true)
	if err != nil {
		fmt.Printf("获取数据失败: %s\n", err)
	}
	fmt.Println(string(resp))
}
func TestGetCurrent2(t *testing.T) {
	fmt.Println(" -=> Start fetch current snap data using GM-api ... ")

	url := gmURL
	timeoutSeconds := 10

	symbols := "SHSE.601088,SHSE.000001"
	resp, err := GetCurrentByte(url, symbols, timeoutSeconds, false)
	// resp, err := GetCurrent(url, symbols, timeoutSeconds, true)
	if err != nil {
		fmt.Printf("获取数据失败: %s\n", err)
	}
	fmt.Println(string(resp))
}

func TestGetDFKbars(t *testing.T) {
	fmt.Println(" -=> Start fetch kbars history data using GM-api ... ")

	url := gmURL
	timeoutSeconds := 10

	symbols := "SHSE.601088,SZSE.300917"
	sdate := "2025-05-01"
	edate := "2025-05-12"
	tag := "1d"

	df, err := DfGetKbars(symbols, tag, sdate, edate, url, timeoutSeconds)
	if err != nil {
		fmt.Printf("获取数据失败: %s\n", err)
	}
	fmt.Println(df)
}

func TestFetchData(t *testing.T) {
	fmt.Println(" -=> Start download test.txt file ... ")

	url := "http://localhost:5002/download/test.txt"

	rsp, err := FetchData(url)
	if err != nil {
		fmt.Printf("获取数据失败: %s\n", err)
	}
	fmt.Println(string(rsp))
}

func TestDownloadCSV(t *testing.T) {
	fmt.Println(" -=> Start download csv.xz file month... ")

	url := "http://localhost:5002/download/kbars-month/month-2025/month-2025-05--SH-60/kbars-1m--SHSE.601088--2025-05-.csv.xz"
	istime := true

	rsp, err := DownloadAndConvertToJSON(url, istime, "timestamp")
	if err != nil {
		fmt.Printf("获取数据失败: %s\n", err)
	}
	fmt.Println(rsp)
}

func TestGetCSVMonth(t *testing.T) {
	fmt.Println(" -=> Start fetch csv.xz file month... ")

	url := gmCSV

	symbol := "SHSE.601088"
	// tag := "1m"
	month := 5
	year := 2025
	istime := true

	rsp, err := GetCSVMonthJson(url, symbol, month, year, istime, 10)
	if err != nil {
		fmt.Printf("获取数据失败: %s\n", err)
	}
	fmt.Println(string(rsp))
}

func TestGetCSVYear(t *testing.T) {
	fmt.Println(" -=> Start fetch csv.xz file year... ")

	url := gmCSV

	symbol := "SHSE.601088"
	// tag := "1m"
	// tag := "vv"
	tag := "pe"
	year := 2025
	istime := true

	rsp, err := GetCSVYearJson(url, symbol, tag, year, istime, "timestamp", 10)
	if err != nil {
		fmt.Printf("获取数据失败: %s\n", err)
	}
	fmt.Println(string(rsp))
}

func TestMap(t *testing.T) {
	fmt.Println(" -=> Start TestMap ... ")
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	dic := make(map[string]any, 1)
	dic["m"] = m
	dic["n"] = map[string]int{"a": 5, "b": 3, "c": 7}
	jsonBytes, _ := json.Marshal(dic)
	// fmt.Println(dic)
	fmt.Println(string(jsonBytes))
	// for k, v := range m {
	// 	fmt.Println(k, v)
	// }
}

func TestGetCSV1m(t *testing.T) {
	fmt.Println(" -=> Start fetch csv.xz file year... ")

	url := gmCSV

	symbol := "SHSE.601088"
	// tag := "1m"
	// tag := "vv"
	// tag := "pe"
	// year := 2025
	istime := true
	isclip := true
	// istime = false
	// isclip = false

	rsp, err := GetCSV1m(url, symbol, "2025-05-29", "2025-05-29", istime, isclip, 10)
	if err != nil {
		fmt.Printf("获取数据失败: %s\n", err)
	}

	jsonData, _ := json.Marshal(rsp[len(rsp)-5:])
	fmt.Println(string(jsonData))
	// fmt.Println(string(rsp))
	// fmt.Println(rsp[:5])
}

func TestGetCSVTag(t *testing.T) {
	fmt.Println(" -=> Start fetch csv.xz file by tag... ")

	url := gmCSV

	symbol := "SHSE.601088"
	// tag := "1m"
	tag := "vv"
	tag = "pe"
	// year := 2025
	istime := true
	isclip := true
	istime = false
	// isclip = false

	rsp, err := GetCSVTag(url, tag, symbol, "2025-06-01", "2025-06-13", istime, isclip, 10)
	if err != nil {
		fmt.Printf("获取数据失败: %s\n", err)
		return
	}

	fmt.Println("=" + strings.Repeat("=", 50))
	// jsonData, _ := json.Marshal(rsp[len(rsp)-5:])
	// jsonData, _ := json.Marshal(rsp[:5])
	jsonData, _ := json.Marshal(rsp)
	fmt.Println(string(jsonData))
	// fmt.Println(string(rsp))
	// fmt.Println(rsp[:5])
	fmt.Println("=" + strings.Repeat("=", 50))
}

func TestGetGM1m(t *testing.T) {
	fmt.Println(" -=> Start fetch gm1m ... ")

	symbol := "SHSE.601088"
	// year := 2025
	istime := true
	isclip := true
	// istime = false
	// isclip = false

	rsp, err := GetGM1m(gmCSV, gmURL, symbol, "2025-06-05", "2025-06-09", istime, isclip, 10)
	if err != nil {
		fmt.Printf("获取数据失败: %s\n", err)
	}

	fmt.Println("=" + strings.Repeat("=", 50))
	jsonData, _ := json.Marshal(rsp[:5])
	fmt.Println(string(jsonData))
	fmt.Println("=" + strings.Repeat("=", 50))
	jsonData, _ = json.Marshal(rsp[len(rsp)-5:])
	fmt.Println(string(jsonData))
	fmt.Println("=" + strings.Repeat("=", 50))
	// fmt.Println(string(rsp))
	// fmt.Println(rsp[:5])
}

func TestGetKbarsByte(t *testing.T) {
	fmt.Println(" -=> Start fetch kbars history data using GM-api ... ")

	url := gmURL
	timeoutSeconds := 30

	// symbols := "SHSE.601088,SZSE.300917"
	symbols := "SHSE.601088"
	sdate := "2025-05-29"
	edate := "2025-05-29"
	// tag := "1d"
	tag := "1m"
	ists := true
	// ists = false

	resp, err := GetKbarsHisByte(url, symbols, tag, sdate, edate, timeoutSeconds)
	if err != nil {
		fmt.Printf("获取数据失败: %s\n", err)
	}
	// fmt.Println(string(resp))

	var rcd RawColData
	unmarshalErr := json.Unmarshal(resp, &rcd)
	if unmarshalErr != nil {
		fmt.Printf("将字节解组为 RawColData 结构体失败: %v\n", unmarshalErr)
	} else {
		// 处理获取到的 JSON 数据为 records 形式
		records, transformErr := rcd.ToRecords()
		if transformErr != nil {
			fmt.Printf("转换为 records 格式失败: %v\n", transformErr)
		} else {
			fmt.Println("\n--- 转换后的 records 格式 ---")
			recordsJSON, _ := json.MarshalIndent(records[:3], "", "  ") // 格式化输出 JSON
			fmt.Printf("%s\n", recordsJSON)

			fmt.Println("=" + strings.Repeat("=", 50))
			records2 := ConvertEob2Timestamp(records, ists)
			recordsJSON, _ = json.MarshalIndent(records2[:3], "", "  ") // 格式化输出 JSON
			fmt.Printf("%s\n", recordsJSON)
		}
	}
}

func TestGetKbars2(t *testing.T) {
	fmt.Println(" -=> Start fetch kbars history data using GM-api ... ")

	url := gmURL
	timeoutSeconds := 30

	// symbols := "SHSE.601088,SZSE.300917"
	symbols := "SHSE.601088"
	sdate := "2025-05-01"
	// sdate = "2025-05-29"
	edate := "2025-05-29"
	tag := "1d"
	// tag := "1m"
	ists := true
	ists = false

	resp, err := GetKbarsHis(url, symbols, tag, sdate, edate, ists, timeoutSeconds)
	if err != nil {
		fmt.Printf("获取数据失败: %s\n", err)
	}
	recordsJSON, _ := json.MarshalIndent(resp[max(0, len(resp)-3):], "", "  ") // 格式化输出 JSON
	// fmt.Printf("%s\n", recordsJSON)
	fmt.Println(string(recordsJSON))
	// fmt.Println(resp[:5])
}

func TestGetFinance(t *testing.T) {
	fmt.Println(" -=> Start fetch data using GM-api ... ")

	url := gmURL
	timeoutSeconds := 30

	// symbols := "SHSE.601088,SZSE.300917"
	symbols := "SHSE.601088"
	sdate := "2024-05-01"
	// sdate = "2025-05-29"
	edate := "2025-05-29"
	rpt_type := "12"
	data_type := "101"

	// resp, err := GetFinancePrime(url, symbols, sdate, edate, "", rpt_type, data_type, timeoutSeconds)
	resp, err := GetFinanceDeriv(url, symbols, sdate, edate, "", rpt_type, data_type, timeoutSeconds)

	if err != nil {
		fmt.Printf("获取数据失败: %s\n", err)
	}
	recordsJSON, _ := json.MarshalIndent(resp[max(0, len(resp)-3):], "", "  ") // 格式化输出 JSON
	// fmt.Printf("%s\n", recordsJSON)
	fmt.Println(string(recordsJSON))
	// fmt.Println(resp[:5])
}

func TestGetFinancePt(t *testing.T) {
	fmt.Println(" -=> Start fetch data using GM-api ... ")

	url := gmURL
	timeoutSeconds := 10

	symbols := "SHSE.601088,SZSE.300917"
	// symbols := "SHSE.601088"
	sdate := "2025-06-03"
	// sdate = "2025-05-29"
	// edate := "2025-05-29"
	rpt_type := "12"
	data_type := "101"

	// resp, err := GetFinancePrimePt(url, symbols, sdate, "", rpt_type, data_type, timeoutSeconds)
	resp, err := GetFinanceDerivPt(url, symbols, sdate, "", rpt_type, data_type, timeoutSeconds)
	if err != nil {
		fmt.Printf("获取数据失败: %s\n", err)
	}
	recordsJSON, _ := json.MarshalIndent(resp[max(0, len(resp)-3):], "", "  ") // 格式化输出 JSON
	// fmt.Printf("%s\n", recordsJSON)
	fmt.Println(string(recordsJSON))
	// fmt.Println(resp[:5])
}

func TestGetDaily(t *testing.T) {
	fmt.Println(" -=> Start fetch data using GM-api ... ")

	url := gmURL
	timeoutSeconds := 30

	// symbols := "SHSE.601088,SZSE.300917"
	symbols := "SHSE.601088"
	sdate := "2025-05-01"
	// sdate = "2025-05-29"
	edate := "2025-05-29"

	// resp, err := GetDailyValuation(url, symbols, sdate, edate, "", timeoutSeconds)
	// resp, err := GetDailyMktvalue(url, symbols, sdate, edate, "", timeoutSeconds)
	resp, err := GetDailyBasic(url, symbols, sdate, edate, "", timeoutSeconds)

	if err != nil {
		fmt.Printf("获取数据失败: %s\n", err)
	}
	recordsJSON, _ := json.MarshalIndent(resp[max(0, len(resp)-3):], "", "  ") // 格式化输出 JSON
	// fmt.Printf("%s\n", recordsJSON)
	fmt.Println(string(recordsJSON))
	// fmt.Println(resp[:5])
}

func TestGetDailyPt(t *testing.T) {
	fmt.Println(" -=> Start fetch data using GM-api ... ")

	url := gmURL
	timeoutSeconds := 10

	symbols := "SHSE.601088,SZSE.300917"
	// symbols := "SHSE.601088"
	sdate := "2025-06-03"
	// sdate = "2025-05-29"
	// edate := "2025-05-29"

	// resp, err := GetDailyValuationPt(url, symbols, sdate, "", timeoutSeconds)
	resp, err := GetDailyMktvaluePt(url, symbols, sdate, "", timeoutSeconds)
	// resp, err := GetDailyBasicPt(url, symbols, sdate, "", timeoutSeconds)
	if err != nil {
		fmt.Printf("获取数据失败: %s\n", err)
	}
	recordsJSON, _ := json.MarshalIndent(resp[max(0, len(resp)-3):], "", "  ") // 格式化输出 JSON
	// fmt.Printf("%s\n", recordsJSON)
	fmt.Println(string(recordsJSON))
	// fmt.Println(resp[:5])
}

func TestGetMarketInfo(t *testing.T) {
	fmt.Println(" -=> Start fetch data using GM-api ... ")

	url := gmURL
	timeoutSeconds := 10

	symbols := "SHSE.601088,SZSE.300917"
	// symbols = "SHSE.601088"
	// symbols = ""
	sec := "stock"
	exchange := "SHSE,SZSE"

	resp, err := GetMarketInfo(url, symbols, sec, exchange, timeoutSeconds)
	if err != nil {
		fmt.Printf("获取数据失败: %s\n", err)
	}
	recordsJSON, _ := json.MarshalIndent(resp[max(0, len(resp)-3):], "", "  ") // 格式化输出 JSON
	fmt.Println(string(recordsJSON))
}

func TestGetSymbolsInfo(t *testing.T) {
	fmt.Println(" -=> Start fetch data using GM-api ... ")

	url := gmURL
	timeoutSeconds := 30

	symbols := "SHSE.601088,SZSE.300917"
	// symbols = "SHSE.601088"
	// symbols = ""
	sec := "stock"
	exchange := "SHSE,SZSE"
	trade_date := "2025-05-29"

	resp, err := GetSymbolsInfo(url, symbols, sec, exchange, trade_date, timeoutSeconds)
	if err != nil {
		fmt.Printf("获取数据失败: %s\n", err)
	}
	recordsJSON, _ := json.MarshalIndent(resp[max(0, len(resp)-3):], "", "  ") // 格式化输出 JSON
	fmt.Println(string(recordsJSON))
}

func TestGetHistoryInfo(t *testing.T) {
	fmt.Println(" -=> Start fetch data using GM-api ... ")

	url := gmURL
	timeoutSeconds := 30

	// symbols := "SHSE.601088,SZSE.300917"
	// symbols = ""
	symbol := "SHSE.601088"
	sdate := "2025-06-01"
	edate := "2025-06-15"

	resp, err := GetHistoryInfo(url, symbol, sdate, edate, timeoutSeconds)
	if err != nil {
		fmt.Printf("获取数据失败: %s\n", err)
	}
	recordsJSON, _ := json.MarshalIndent(resp[max(0, len(resp)-3):], "", "  ") // 格式化输出 JSON
	fmt.Println(string(recordsJSON))
}

func TestConvertTimestamp(t *testing.T) {
	fmt.Println(" -=> Start TestConvertTimestamp ... ")
	timestamp := " 2025-05-29 09:34:00+08:00 "
	// timestamp := "2025-05-29 09:34:00"
	// tt, _ := time.Parse("2006-01-02T15:04:05+08:00", timestamp)
	timestamp = strings.TrimSpace(timestamp)
	timestamp = strings.TrimSuffix(timestamp, "+08:00")
	timestamp = strings.Replace(timestamp, " ", "T", 1)
	// if timestamp[len(timestamp)-6:] == "+08:00" {
	// 	timestamp = timestamp[:len(timestamp)-6]
	// }
	// tt, _ := time.Parse(time.RFC3339, timestamp)
	tz, _ := time.LoadLocation("Asia/Shanghai")
	tt, _ := time.ParseInLocation("2006-01-02T15:04:05", timestamp, tz)
	// tt, _ := time.Parse("2006-01-02T15:04:05", timestamp)
	// fmt.Println(ConvertTimeFormatByReplace(timestamp))
	fmt.Println(tt)
	fmt.Println("=" + strings.Repeat("=", 50))
	fmt.Println(tt.Local())
	fmt.Println(tt.Unix())
	fmt.Println(tt.UnixMilli())
}
