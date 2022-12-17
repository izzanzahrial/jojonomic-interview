package entities

type Price struct {
	AdminID      string `json:"admin_id"`
	TopupPrice   int32  `json:"harga_topup"`
	BuybackPrice int32  `json:"harga_buyback"`
}

func (p Price) Validate() map[string]string {
	var errors = make(map[string]string)

	if ok := p.validateTopupPrice(); !ok {
		errors["topup"] = "top up price cannot be lower than 0"
	}

	if ok := p.validateBuybackPrice(); !ok {
		errors["buyback"] = "buy back price cannot be lower than 0"
	}

	return errors
}

func (p Price) validateTopupPrice() bool {
	return p.TopupPrice > 0
}

func (p Price) validateBuybackPrice() bool {
	return p.BuybackPrice > 0
}

func NewPrice(id string, topup, buyback int32) *Price {
	return &Price{
		AdminID:      id,
		TopupPrice:   topup,
		BuybackPrice: buyback,
	}
}
