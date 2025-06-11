package gm

import (
	"encoding/json"
	"fmt"
)

// RawColData 结构体用于匹配从 URL 获取的原始 JSON 格式
type RawColData struct {
	Columns []string `json:"columns"`
	Data    [][]any  `json:"data"` // 使用 interface{} 来处理不同类型的数据
}

// transformToRecords 将 InputData 格式转换为 records 格式
// records 格式是一个 map 数组，每个 map 的键是列名，值是对应的数据
func (rcd *RawColData) FromByte(bdata []byte) error {
	unmarshalErr := json.Unmarshal(bdata, rcd)
	if unmarshalErr != nil {
		return fmt.Errorf(" 解析原始数据失败: %s", unmarshalErr)
	}
	return nil
}

// transformToRecords 将 InputData 格式转换为 records 格式
// records 格式是一个 map 数组，每个 map 的键是列名，值是对应的数据
func (rcd *RawColData) ToRecords() ([]map[string]any, error) {
	var records []map[string]any

	if len(rcd.Columns) == 0 && len(rcd.Data) > 0 {
		return nil, fmt.Errorf("当存在数据时，列名不能为空")
	}

	for _, row := range rcd.Data {
		// 检查数据行长度是否与列名长度匹配
		if len(row) != len(rcd.Columns) {
			return nil, fmt.Errorf("数据行长度 (%d) 与列名长度 (%d) 不匹配", len(row), len(rcd.Columns))
		}
		record := make(map[string]any)
		for i, colName := range rcd.Columns {
			// if istimestamp && (colName == "timestamp" || colName == "eob") {
			// 	tt := ConvertString2Time(fmt.Sprintf("%s", row[i]))
			// 	record[colName] = tt.UnixMilli()
			// } else {
			// 	record[colName] = row[i]
			// }
			record[colName] = row[i]
		}
		records = append(records, record)
	}
	return records, nil
}
