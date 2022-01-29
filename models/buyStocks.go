package models

type BuyStock struct {
	ID     uint
	Ticker string
	Number int
	Price  float32
	Fee    float32
	Type   string
}
