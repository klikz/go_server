package model

type Galileo struct {
	Barcode        string  `json:"Barcode"`
	OpCode         string  `json:"OpCode"`
	TypeFreon      string  `json:"TypeFreon"`
	Result         string  `json:"Result"`
	ProgQuantity   float32 `json:"ProgQuantity"`
	Quantity       float32 `json:"Quantity"`
	CycleTotalTime int     `json:"CycleTotalTime"`
	Time           string  `json:"Time"`
}
