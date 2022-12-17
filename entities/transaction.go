package entities

type Transaction struct {
	Date         int32   `json:"date"`
	Type         string  `json:"type"`
	Gram         float32 `json:"gram"`
	TopupPrice   int32   `json:"harga_topup"`
	BuybackPrice int32   `json:"harga_buyback"`
	Saldo        float32 `json:"saldo"`
}
