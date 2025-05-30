package gm

import "fmt"

type AskBidData struct {
	BidPrice string  `json:"bid_p"`
	BidValue float64 `json:"bid_v"`
	AskPrice float64 `json:"ask_p"`
	AskValue string  `json:"ask_v"`
}

type SnapData struct {
	Symbol      string       `json:"symbol"`
	Open        float64      `json:"open"`
	High        float64      `json:"high"`
	Low         float64      `json:"low"`
	Price       float64      `json:"price"`
	CumVolumn   int64        `json:"cum_volume"`
	CumAmount   float64      `json:"cum_amount"`
	TradeType   int64        `json:"trade_type"`
	CreateAt    string       `json:"create_at"`
	CumPosition string       `json:"cum_position"`
	LastAmount  string       `json:"last_amount"`
	LastVolume  int64        `json:"last_volume"`
	Flag        int64        `json:"flag"`
	Iopv        int64        `json:"iopv"`
	Quotes      []AskBidData `json:"quotes"`
}

// 交易日历日期结构体
type TradeDate struct {
	NextTradeDate string `json:"next_trade_date"`
	PrevTradeDate string `json:"pre_trade_date"`
	TradeDate     string `json:"trade_date"`
}

type TradeCalendar struct {
	Calendar map[string]TradeDate
}

type Kbar struct {
	Timestamp string  `json:"timestamp"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    int64   `json:"volume"`
}

type KbarList struct {
	Kbars map[string][]Kbar
}

// RawColData 结构体用于匹配从 URL 获取的原始 JSON 格式
type RawColData struct {
	Columns []string `json:"columns"`
	Data    [][]any  `json:"data"` // 使用 interface{} 来处理不同类型的数据
}

// transformToRecords 将 InputData 格式转换为 records 格式
// records 格式是一个 map 数组，每个 map 的键是列名，值是对应的数据
func (rcd *RawColData) TransformToRecords() ([]map[string]any, error) {
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
			record[colName] = row[i]
		}
		records = append(records, record)
	}
	return records, nil
}
