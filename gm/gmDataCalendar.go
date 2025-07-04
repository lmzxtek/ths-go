package gm

// 交易日历日期结构体
type TradeDate struct {
	NextTradeDate string `json:"next_trade_date"`
	PrevTradeDate string `json:"pre_trade_date"`
	TradeDate     string `json:"trade_date"`
}

type TradeCalendar struct {
	Calendar map[string]TradeDate
}
