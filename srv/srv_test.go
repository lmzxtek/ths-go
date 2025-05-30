package srv

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

// 示例用法
func TestSmartURLHandler(t *testing.T) {
	// 测试用例
	testURLs := []string{
		"example.com",
		"example.com:443",
		"example.com:80",
		"example.com:8080",
		"example.com:3000",
		"https://example.com",
		"http://example.com:8000",
		"",
	}

	fmt.Println("测试 preferHTTPS = false:")
	for _, url := range testURLs {
		result := SmartURLHandler(url, false)
		fmt.Printf("输入: %-25s -> 输出: %s\n", fmt.Sprintf("'%s'", url), result)
	}

	fmt.Println("\n测试 preferHTTPS = true:")
	for _, url := range testURLs {
		result := SmartURLHandler(url, true)
		fmt.Printf("输入: %-25s -> 输出: %s\n", fmt.Sprintf("'%s'", url), result)
	}
}

// TestReadCSV is a test function to read CSV file and print the dataframe
func TestGetURLRetry(t *testing.T) {
	fmt.Println("\n >>> Start Test GetURL... ")

	// 示例 URL 和参数
	// 注意：这些是示例URL，实际使用时请替换为可访问的API端点。
	// 例如，一个返回JSON的公共API，如 JSONPlaceholder。
	// 假设我们有一个返回 {"message": "Hello, world!"} 的端点
	exampleURL := "https://jsonplaceholder.typicode.com/posts/1" // 这是一个返回JSON的公共测试API
	exampleParams := map[string]string{
		"userId": "1",
		"id":     "1",
	}

	// 1. 使用 getURLWithoutRetry 函数的例子
	fmt.Println("--- 调用 getURLWithoutRetry ---")
	// connectTimeout 设置为 5 秒，dataTimeout 默认为 connectTimeout
	dataWithoutRetry, err := GetURLWithoutRetry(exampleURL, exampleParams, 5*time.Second, 0)
	if err != nil {
		fmt.Printf("getURLWithoutRetry 发生错误: %v\n", err)
	} else {
		fmt.Printf("getURLWithoutRetry 成功响应: %+v\n", dataWithoutRetry)
	}

	fmt.Println("\n--- 调用 getURLWithoutRetry (模拟失败情况) ---")
	// 尝试一个不存在的 URL 或端口，模拟连接失败或HTTP错误
	// 注意：这个URL可能无法访问，从而触发错误
	badURL := "http://localhost:9999/nonexistent"
	_, err = GetURLWithoutRetry(badURL, nil, 2*time.Second, 0)
	if err != nil {
		fmt.Printf("getURLWithoutRetry 模拟失败错误: %v\n", err)
	} else {
		fmt.Println("getURLWithoutRetry 模拟失败情况意外成功。")
	}

	// 2. 使用 getURLWithRetry 函数的例子
	fmt.Println("\n--- 调用 getURLWithRetry ---")
	// connectTimeout 设置为 5 秒，dataTimeout 默认为 connectTimeout
	dataWithRetry, err := GetURLWithRetry(exampleURL, exampleParams, 5*time.Second, 0)
	if err != nil {
		fmt.Printf("getURLWithRetry 发生错误: %v\n", err)
	} else {
		fmt.Printf("getURLWithRetry 成功响应: %+v\n", dataWithRetry)
	}

	fmt.Println("\n--- 调用 getURLWithRetry (模拟重试失败情况) ---")
	// 尝试一个会持续失败的 URL，观察重试机制
	// 这个URL很可能在多次重试后仍然失败
	const alwaysFailURL = "http://localhost:12345/alwaysfail" // 假设这是一个永远无法连接的地址
	_, err = GetURLWithRetry(alwaysFailURL, nil, 1*time.Second, 0)
	if err != nil {
		fmt.Printf("getURLWithRetry 模拟重试失败错误: %v\n", err)
	} else {
		fmt.Println("getURLWithRetry 模拟重试失败情况意外成功。")
	}
}

// TestReadCSV is a test function to read CSV file and print the dataframe
func TestGetGMURL(t *testing.T) {
	fmt.Println("\n >>> Start TestGetGMURL ... ")

	// gmapi := "http://localhost:5000/test"
	gmapi := "http://localhost:5000/test2"
	// gmapi := "http://localhost:5002/test3"
	pars := map[string]string{
		"userId": "1",
	}

	// 1. 使用 getURLWithoutRetry 函数的例子
	fmt.Println("--- 调用 getURLWithoutRetry ---")
	// connectTimeout 设置为 5 秒，dataTimeout 默认为 connectTimeout
	rawData, err := GetURLWithoutRetry(gmapi, pars, 5*time.Second, 0)
	if err != nil {
		fmt.Printf("getURLWithoutRetry 发生错误: %v\n", err)
	}

	fmt.Printf("getURLWithoutRetry 成功响应: %+v\n", rawData)

	// 将获取到的原始数据转换为 InputData 结构体
	jsonBytes, marshalErr := json.Marshal(rawData)
	if marshalErr != nil {
		fmt.Printf("将原始数据编组为字节失败: %v\n", marshalErr)
	} else {
		var inputData RawColData
		unmarshalErr := json.Unmarshal(jsonBytes, &inputData)
		if unmarshalErr != nil {
			fmt.Printf("将字节解组为 InputData 结构体失败: %v\n", unmarshalErr)
		} else {
			// 处理获取到的 JSON 数据为 records 形式
			records, transformErr := inputData.TransformToRecords()
			if transformErr != nil {
				fmt.Printf("转换为 records 格式失败: %v\n", transformErr)
			} else {
				fmt.Println("\n--- 转换后的 records 格式 (getURLWithRetry) ---")
				recordsJSON, _ := json.MarshalIndent(records, "", "  ") // 格式化输出 JSON
				fmt.Printf("%s\n", recordsJSON)
			}
		}
	}
}

// TestReadCSV is a test function to read CSV file and print the dataframe
func TestGetDatesByYear(t *testing.T) {
	fmt.Println("\n >>> Start TestGetDatesByYear ... ")

	year := "2025"
	gmapi := "http://localhost:5000/get_dates_by_year"
	pars := map[string]string{
		"syear": year,
	}

	// 1. 使用 getURLWithoutRetry 函数的例子
	fmt.Println("--- 调用 getURLWithoutRetry ---")
	// connectTimeout 设置为 5 秒，dataTimeout 默认为 connectTimeout
	rawData, err := GetURLWithoutRetry(gmapi, pars, 5*time.Second, 0)
	if err != nil {
		fmt.Printf("getURLWithoutRetry 发生错误: %v\n", err)
	}

	// fmt.Printf("getURLWithoutRetry 成功响应: %+v\n", rawData)

	// 将获取到的原始数据转换为 InputData 结构体
	jsonBytes, marshalErr := json.Marshal(rawData)
	if marshalErr != nil {
		fmt.Printf("将原始数据编组为字节失败: %v\n", marshalErr)
	} else {
		var inputData RawColData
		unmarshalErr := json.Unmarshal(jsonBytes, &inputData)
		if unmarshalErr != nil {
			fmt.Printf("将字节解组为 InputData 结构体失败: %v\n", unmarshalErr)
		} else {
			// 处理获取到的 JSON 数据为 records 形式
			records, transformErr := inputData.TransformToRecords()
			if transformErr != nil {
				fmt.Printf("转换为 records 格式失败: %v\n", transformErr)
			} else {
				fmt.Println("\n--- 转换后的 records 格式 (getURLWithRetry) ---")
				recordsJSON, _ := json.MarshalIndent(records[:5], "", "  ") // 格式化输出 JSON
				fmt.Printf("%s\n", recordsJSON)
			}
		}
	}
}
