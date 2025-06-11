package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lmzxtek/ths-go/gm"
)

var gmURL = "http://45.154.14.186:5000" // locVPS-kr
// var gmURL = "http://111.67.205.166:5000" // uDouYun-bj

func main() {
	// fmt.Println(" -=> Start fetch url... ")
	args := os.Args[1:] // os.Args[0] 是脚本名，后面是参数

	if len(args) > 0 && args[0] == "json" {
		fmt.Println(" -=> Test json parse to dataframe ... ")
		jsonStr := `{"columns":["Symbol","Time","Price","Volume"],"data":[["AAPL","2025-05-01",100,200],["AAPL","2025-05-01",101,210],["AAPL","2025-05-01",102,220],["AAPL","2025-05-01",103,230],["AAPL","2025-05-01",104,240]]}`
		df, err := gm.ParseJsonToDataframe(jsonStr)
		if err != nil {
			fmt.Printf("获取数据失败: %s\n", err)
		}
		fmt.Println(df)

	} else if len(args) > 0 && args[0] == "df" {
		fmt.Println(" -=> Start fetch df from url ... ")
		df, err := gm.DfGetTest(gmURL)
		if err != nil {
			fmt.Printf("获取数据失败: %s\n", err)
		}
		fmt.Println(df)
		// cxz.SaveDataframeToCSV(&df, "data.csv")
		// cxz.SaveDataframeToCSVxz(&df, "data.csv.xz")

	} else if len(args) > 0 && args[0] == "test2" {
		fmt.Println(" -=> Start fetch df from url after parse json data {columns, data} ... ")
		df, err := gm.DfGetTest2(gmURL)
		if err != nil {
			fmt.Printf("获取数据失败: %s\n", err)
		}
		fmt.Println(df)

	} else if len(args) > 0 && args[0] == "snap" {
		fmt.Println(" -=> Start fetch current snap data using GM-api ... ")

		url := gmURL
		timeoutSeconds := 10

		symbols := "SHSE.601088,SZSE.300917"

		// df, err := gm.DfGetCurrent(symbols, url, timeoutSeconds)
		// if err != nil {
		// 	fmt.Printf("获取数据失败: %s\n", err)
		// }
		// fmt.Println(df)

		resp, err := gm.GetCurrentByte(symbols, url, timeoutSeconds, true)
		if err != nil {
			fmt.Printf("获取数据失败: %s\n", err)
		}
		fmt.Println(string(resp))

	} else if len(args) > 0 && args[0] == "his" {
		fmt.Println(" -=> Start fetch current his data using GM-api ... ")

		url := gmURL
		timeoutSeconds := 10

		symbols := "SHSE.601088,SZSE.300917"
		sdate := "2025-05-01"
		edate := "2025-05-12"
		tag := "1d"

		df, err := gm.DfGetKbars(symbols, tag, sdate, edate, url, timeoutSeconds)
		if err != nil {
			fmt.Printf("获取数据失败: %s\n", err)
		}
		fmt.Println(df)

		// resp, err := gm.GetKbarsHis(symbols, tag, sdate, edate, url, timeoutSeconds)
		// if err != nil {
		// 	fmt.Printf("获取数据失败: %s\n", err)
		// }
		// fmt.Println(string(resp))

	} else if len(args) > 0 && args[0] == "jsondf" {
		jsonData := []byte(`{
			"columns": ["Symbol","Time","Price","Volume"],
			"data": [["AAPL","2025-05-01",100,200],["AAPL","2025-05-01",101,210]]
		}`)

		// 解析原始数据
		var raw gm.RawData
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
		df, err := gm.ConvertToColumnar(raw)
		if err != nil {
			panic(err)
		}

		// 示例输出
		fmt.Println("Symbol列:", df["Symbol"])
		fmt.Println("Price列 :", df["Price"])

	} else if len(args) > 0 && args[0] == "csv" {
		// test.TestReadCSV()
		// test.TestloadRecords()
		// test.TestReadJSON()
		// TestloadStructs()

		jsonStr := `{"columns":["Symbol","Time","Price","Volume"],"data":[["AAPL","2025-05-01",100,200],["AAPL","2025-05-01",101,210],["AAPL","2025-05-01",102,220],["AAPL","2025-05-01",103,230],["AAPL","2025-05-01",104,240]]}`
		df, err := gm.ParseJsonToDataframe(jsonStr)
		if err != nil {
			fmt.Printf("获取数据失败: %s\n", err)
		}
		fmt.Println(df)

	} else {
		fmt.Println()
		fmt.Println(" >>> Opps: cmd options or params error:", args, "...")
		fmt.Println("           > go run testgm.go json ")
		fmt.Println("           > go run testgm.go test2 ")
		fmt.Println("           > go run testgm.go df ")
		fmt.Println()
	}

}
