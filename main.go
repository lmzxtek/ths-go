package main

import (
	"fmt"
	"os"
	"ths/cxz"
	"ths/gm"
)

// var gmURL = "http://45.154.14.186:5000"  // locVPS-kr
var gmURL = "http://111.67.205.166:5000" // uDouYun-bj

func main() {
	args := os.Args[1:] // os.Args[0] 是脚本名，后面是参数

	// 示例用法
	filePath := "data.csv" // 替换为你的 CSV 文件路径

	if len(args) > 0 && args[0] == "test" {
		fmt.Println(" >>> Start to test script ...")
		// 使用第一个版本
		records, err := cxz.ReadCSVFile(filePath)
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}
		fmt.Println("CSV 文件内容:")
		for i, record := range records {
			fmt.Printf("%d: %v\n", i+1, record)
		}

	} else if len(args) > 0 && args[0] == "update" {
		fmt.Println(" >>> Start update THS daily indications ...")
		// 或者使用更高效的版本
		records, err := cxz.ReadCSVFileEfficient(filePath)
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}

		fmt.Println("CSV 文件内容:")
		for i, record := range records {
			fmt.Printf("%d: %v\n", i+1, record)
		}

		// df := dataframe.ReadCSV(records)
		// fmt.Println(df)

	} else if len(args) > 0 && args[0] == "url" {
		fmt.Println(" -=> Start fetch url... ")

		// url := "http://45.154.14.186:5000/test"
		// makeRequest(url)

	} else if len(args) > 0 && args[0] == "kline" {
		fmt.Println(" -=> Start fetch url... ")
		// url := fmt.Sprintf("%s/get_current", gmURL)
		// url := fmt.Sprintf("%s/get_his", gmURL)
		// url := fmt.Sprintf("%s/test2", gmURL)
		// url := fmt.Sprintf("%s/test", gmURL)

		// fmt.Println(" -=> Start fetch data from url... ")
		// klineData, err := gm.GetTest(gmURL)
		// 获取历史K线数据
		// klineData, err := gm.FetchURLData(url, map[string]string{
		// 	"symbols": "SHSE.000001,SZSE.300917",
		// 	"tag":     "1d",
		// 	"sdate":   "2025-05-01",
		// 	"edate":   "2025-05-12",
		// })
		// if err != nil {
		// 	fmt.Printf("获取数据失败: %s\n", err)
		// }
		// fmt.Printf("\n获取到的数据: \n%s\n", klineData)
		// fmt.Printf("获取到的数据: \n%s\n", klineData)

		// df := dataframe.ReadJSON(strings.NewReader(klineData))
		// newreader := strings.NewReader(string(klineData))
		// df := dataframe.ReadJSON(strings.NewReader(string(klineData)))
		fmt.Println(" -=> Start fetch df from url... ")
		df, err := gm.DfGetTest(gmURL)
		if err != nil {
			fmt.Printf("获取数据失败: %s\n", err)
		}
		fmt.Println(df)

		cxz.SaveDataframeToCSV(&df, "data.csv")
		cxz.SaveDataframeToCSVxz(&df, "data.csv.xz")

		// // 解析数据
		// var klines dfSplitData
		// if err := parseJSON(klineData, &klines); err != nil {
		// 	fmt.Printf("数据解析失败: %s\n", err)
		// }
		// fmt.Printf("\n解析到的数据: \n%s\n", klines.Columns[:])
		// fmt.Printf("\n解析到的数据: \n%s\n", klines.Data[0])
		// fmt.Printf("\n解析到的数据: \n%s\n", klines.Data[1][2])

		// df := dataframe.ReadCSV(klines)
		// fmt.Println(df)

	} else if len(args) > 0 && args[0] == "df" {
		// testReadCSV()
		// testloadRecords()
		// testReadJSON()
		// TestloadStructs()

		jsonStr := `{"columns":["Symbol","Time","Price","Volume"],"data":[["AAPL","2025-05-01",100,200],["AAPL","2025-05-01",101,210],["AAPL","2025-05-01",102,220],["AAPL","2025-05-01",103,230],["AAPL","2025-05-01",104,240]]}`
		df, err := gm.ParseJsonToDataframe(jsonStr)
		if err != nil {
			fmt.Printf("数据解析失败: %s\n", err)
		}
		fmt.Println(df)

	} else {
		fmt.Println()
		fmt.Println(" >>> Opps: cmd options or params error:", args, "...")
		fmt.Println("           > go run main.go test ")
		fmt.Println("           > go run main.go update ")
		fmt.Println()
	}
}
