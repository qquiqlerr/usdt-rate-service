package models

type Rate struct {
	Market    string  `json:"market"`
	AskPrice  float64 `json:"ask"`
	BidPrice  float64 `json:"bid"`
	Timestamp int64   `json:"timestamp"`
}

type Depth struct {
	Timestamp int64   `json:"timestamp"`
	Asks      []Order `json:"asks"`
	Bids      []Order `json:"bids"`
}

type Order struct {
	Price  float64 `json:"price"`
	Volume float64 `json:"volume"`
	Amount float64 `json:"amount"`
	Factor float64 `json:"factor"`
	Type   string  `json:"type"`
}
