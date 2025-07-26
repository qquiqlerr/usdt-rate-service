package grinex

type DepthResponse struct {
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
