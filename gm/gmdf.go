package gm

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/go-gota/gota/dataframe"
)

type JsonData struct {
	Columns []string `json:"columns"`
	Data    [][]any  `json:"data"`
}

// SplitJSONFormat 定义 Pandas split 格式的 JSON 结构
type SplitJSONFormat struct {
	Index   []any    `json:"index"`   // 行索引
	Columns []string `json:"columns"` // 列名
	Data    [][]any  `json:"data"`    // 数据
}

// parseSplitJSON 解析 Pandas split 格式的 JSON 数据
func ParseSplitJSON(data []byte) (*SplitJSONFormat, error) {
	var result SplitJSONFormat
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %v", err)
	}
	return &result, nil
}

// 解析 JSON 字符串为 DataFrame
func ParseJsonToDataframe(jsonStr string) (dataframe.DataFrame, error) {
	// 假设从 URL 获取的 JSON 字符串如下：
	// jsonStr := `{
	//   "columns":["Symbol","Time","Price","Volume"],
	//   "data":[["AAPL","2025-05-01",100,200],["AAPL","2025-05-01",101,210],["AAPL","2025-05-01",102,220],["AAPL","2025-05-01",103,230],["AAPL","2025-05-01",104,240]]
	// }`

	// 解析 JSON
	var jd JsonData
	err := json.Unmarshal([]byte(jsonStr), &jd)
	if err != nil {
		log.Fatalf("JSON 解析失败: %v", err)
		return dataframe.DataFrame{}, err
	}

	// 将 data 转换为 []map[string]interface{}
	var records []map[string]any
	for _, row := range jd.Data {
		record := make(map[string]any)
		for i, col := range jd.Columns {
			record[col] = row[i]
		}
		records = append(records, record)
	}

	// 构造 DataFrame
	df := dataframe.LoadMaps(records)

	return df, nil
}

// 定义原始数据结构
type RawData struct {
	Columns []string `json:"columns"`
	Data    [][]any  `json:"data"`
	// Data [][]string `json:"data"`
}

// 转换为按列存储的Map结构
func ConvertToColumnar(raw RawData) (map[string][]any, error) {
	columns := make(map[string][]any)

	// 初始化各列
	for _, col := range raw.Columns {
		columns[col] = make([]any, 0, len(raw.Data))
	}

	// 填充数据
	for _, row := range raw.Data {
		if len(row) != len(raw.Columns) {
			return nil, fmt.Errorf("row length mismatch")
		}
		for i, val := range row {
			columns[raw.Columns[i]] = append(columns[raw.Columns[i]], val)
		}
	}

	return columns, nil
}

// 测试数据1
func DfGetTest(url string) (dataframe.DataFrame, error) {
	// tarurl := fmt.Sprintf("%s/test", url)

	// 获取历史K线数据
	resp, err := GetTest(url, 10)
	if err != nil {
		// fmt.Printf("获取数据失败: %s\n", err)
		return dataframe.DataFrame{}, err
	}

	df := dataframe.ReadJSON(strings.NewReader(string(resp)))
	return df, nil
}

// 测试数据2
func DfGetTest2(url string) (dataframe.DataFrame, error) {
	// tarurl := fmt.Sprintf("%s/test", url)

	// 获取历史K线数据
	resp, err := GetTest2(url, 10)
	if err != nil {
		// fmt.Printf("获取数据失败: %s\n", err)
		return dataframe.DataFrame{}, err
	}

	df, err := ParseJsonToDataframe(string(resp))
	if err != nil {
		// fmt.Printf("解析数据失败: %s\n", err)
		return dataframe.DataFrame{}, err
	}

	// df := dataframe.ReadJSON(strings.NewReader(string(resp)))

	return df, nil
}

// 获取最新行情快照，返回 DataFrame
func DfGetCurrent(symbols string, url string, timeoutSeconds int) (dataframe.DataFrame, error) {
	// tarurl := fmt.Sprintf("%s/get_current", url)

	// 获取历史K线数据
	resp, err := GetCurrent(symbols, url, timeoutSeconds)
	if err != nil {
		// fmt.Printf("获取数据失败: %s\n", err)
		return dataframe.DataFrame{}, err
	}

	df, err := ParseJsonToDataframe(string(resp))
	if err != nil {
		// fmt.Printf("解析数据失败: %s\n", err)
		return dataframe.DataFrame{}, err
	}
	return df, nil
}

// 获取K线数据，返回 DataFrame
func DfGetKbars(symbols string, tag string, sdate string, edate string, url string, timeoutSeconds int) (dataframe.DataFrame, error) {
	// tarurl := fmt.Sprintf("%s/get_current", url)

	// 获取历史K线数据
	resp, err := GetKbarsHis(symbols, tag, sdate, edate, url, timeoutSeconds)
	if err != nil {
		// fmt.Printf("获取数据失败: %s\n", err)
		return dataframe.DataFrame{}, err
	}

	df, err := ParseJsonToDataframe(string(resp))
	if err != nil {
		// fmt.Printf("解析数据失败: %s\n", err)
		return dataframe.DataFrame{}, err
	}
	return df, nil
}
