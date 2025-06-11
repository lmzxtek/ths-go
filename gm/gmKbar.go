package gm

import (
	"fmt"
	"time"
)

type KBar struct {
	// Timestamp string  `json:"timestamp"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume int64   `json:"volume"`
}

// KLineData 结构体( klinechart 图表数据)
type KLineData struct {
	Timestamp int64   `json:"timestamp"` // 时间戳
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    int64   `json:"volume"`
	Turnover  float64 `json:"turnover"` // 成交额, 非必须字段，如果需要展示技术指标'EMV'和'AVP'，则需要为该字段填充数据。
}

type KbarMapAny struct {
	Datalist map[any]KBar `json:"timestamp"`
}

func (kb *KBar) ToList(ts any) []any {
	var records []any
	records = append(records, ts)
	records = append(records, kb.Open)
	records = append(records, kb.High)
	records = append(records, kb.Low)
	records = append(records, kb.Close)
	records = append(records, kb.Volume)
	return records
}

// records 格式是一个 map 数组，每个 map 的键是列名，值是对应的数据
func (kb *KBar) ToRecords(ts any) map[any]any {
	// var records map[any]any
	records := make(map[any]any, 6)
	records["timestamp"] = ts
	records["open"] = kb.Open
	records["high"] = kb.High
	records["low"] = kb.Low
	records["close"] = kb.Close
	records["volume"] = kb.Volume
	return records
}

// records 格式是一个 map 数组，每个 map 的键是列名，值是对应的数据
func (kb *KBar) ToKLineDataFromString(ts string) KLineData {
	chinaLocation, _ := time.LoadLocation("Asia/Shanghai")
	tt, err := time.ParseInLocation("2006-01-02 15:04:05", ts, chinaLocation)
	if err != nil {
		fmt.Printf("解析时间失败: %v", err)
	}
	kl := KLineData{
		Timestamp: tt.UnixMilli(),
		Open:      kb.Open,
		High:      kb.High,
		Low:       kb.Low,
		Close:     kb.Close,
		Volume:    kb.Volume,
	}
	return kl
}

// records 格式是一个 map 数组，每个 map 的键是列名，值是对应的数据
func (kb *KBar) ToKLineDataFromInt(ti int64) KLineData {
	kl := KLineData{
		Timestamp: ti,
		Open:      kb.Open,
		High:      kb.High,
		Low:       kb.Low,
		Close:     kb.Close,
		Volume:    kb.Volume,
	}
	return kl
}
