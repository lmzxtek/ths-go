package gm

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

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

	resp, err := GetPrevN(url, date, count, timeoutSeconds)
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

	resp, err := GetNextN(url, date, count, timeoutSeconds)
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
	resp, err := GetCurrent(url, symbols, timeoutSeconds, true)
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
	resp, err := GetCurrent(url, symbols, timeoutSeconds, false)
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

	rsp, err := DownloadAndConvertToJSON(url)
	if err != nil {
		fmt.Printf("获取数据失败: %s\n", err)
	}
	fmt.Println(rsp)
}

func TestGetCSVMonth(t *testing.T) {
	fmt.Println(" -=> Start fetch csv.xz file month... ")

	url := gmCSV

	symbol := "SHSE.601088"
	tag := "1m"
	month := 5
	year := 2025

	rsp, err := GetCSVMonth(url, symbol, tag, month, year, 10)
	if err != nil {
		fmt.Printf("获取数据失败: %s\n", err)
	}
	fmt.Println(string(rsp))
}

func TestGetKbars(t *testing.T) {
	fmt.Println(" -=> Start fetch kbars history data using GM-api ... ")

	url := gmURL
	timeoutSeconds := 10

	// symbols := "SHSE.601088,SZSE.300917"
	symbols := "SHSE.601088"
	sdate := "2025-05-29"
	edate := "2025-05-29"
	// tag := "1d"
	tag := "1m"

	resp, err := GetKbarsHis(symbols, tag, sdate, edate, url, timeoutSeconds)
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
		records, transformErr := rcd.TransformToRecords()
		if transformErr != nil {
			fmt.Printf("转换为 records 格式失败: %v\n", transformErr)
		} else {
			fmt.Println("\n--- 转换后的 records 格式 ---")
			recordsJSON, _ := json.MarshalIndent(records[:5], "", "  ") // 格式化输出 JSON
			fmt.Printf("%s\n", recordsJSON)
		}
	}
}
