package models

// Rate represents a rate for a specific market.
type Rate struct {
	Market    string `json:"market"`
	AskPrice  string `json:"ask"`
	BidPrice  string `json:"bid"`
	Timestamp int64  `json:"timestamp"`
}

// Depth represents the depth data for a specific market.
type Depth struct {
	Timestamp int64   `json:"timestamp"`
	Asks      []Order `json:"asks"`
	Bids      []Order `json:"bids"`
}

// Order represents an order in the depth data.
type Order struct {
	Price  string `json:"price"`
	Volume string `json:"volume"`
	Amount string `json:"amount"`
	Factor string `json:"factor"`
	Type   string `json:"type"`
}
