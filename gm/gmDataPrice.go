package gm

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

// ========================================
