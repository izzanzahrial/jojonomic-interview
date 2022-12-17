package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/izzanzahrial/jojonomic-interview/check-harga-service/repository"
	"github.com/izzanzahrial/jojonomic-interview/common"
)

type CheckPrice struct {
	repository *repository.CheckPrice
}

func NewCheckPrice(repository *repository.CheckPrice) *CheckPrice {
	return &CheckPrice{
		repository: repository,
	}
}

func (cp *CheckPrice) CheckGoldPrice(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	price, err := cp.repository.CheckGoldPrice(ctx)
	if err != nil {
		common.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	var data struct {
		Buyback int32 `json:"harga_buyback"`
		Topup   int32 `json:"harga_topup"`
	}

	data.Buyback = price.BuybackPrice
	data.Topup = price.TopupPrice

	common.SuccessResponse(w, r, http.StatusOK, common.Envolope{"error": false, "data": data})
}
