package kbs

// 定义用于解析 JSON 数据的结构体
type StockData struct {
	Symbol string `json:"Symbol"`
	Time   string `json:"Time"`
	Price  int64  `json:"Price"`
	Volume int64  `json:"Volume"`
}

// 数据结构示例
type KLine struct {
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Volume float64 `json:"volume"`
	Time   string  `json:"time"`
}
