package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/izzanzahrial/jojonomic-interview/common"
	"github.com/izzanzahrial/jojonomic-interview/entities"
	"github.com/izzanzahrial/jojonomic-interview/input-harga-service/usecase"
	"github.com/teris-io/shortid"
)

type Inputprice struct {
	service *usecase.Inputprice
	sID     *shortid.Shortid
}

func NewInputPrice(service *usecase.Inputprice, sid *shortid.Shortid) *Inputprice {
	return &Inputprice{
		service: service,
		sID:     sid,
	}
}

func (i *Inputprice) Publish(w http.ResponseWriter, r *http.Request) {
	var Input struct {
		AdminID      string `json:"admin_id"`
		TopupPrice   int32  `json:"harga_topup"`
		BuybackPrice int32  `json:"harga_buyback"`
	}

	err := common.ReadJSON(w, r, &Input)
	if err != nil {
		common.BadRequestResponse(w, r, err)
		return
	}

	price := entities.NewPrice(Input.AdminID, Input.TopupPrice, Input.BuybackPrice)

	if err := price.Validate(); len(err) != 0 {
		common.FailedValidationResponse(w, r, err)
		return
	}

	reffID, err := i.sID.Generate()
	if err != nil {
		common.GenerateReffIDErrorResponse(w, r, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = i.service.Publish(ctx, reffID, price)
	if err != nil {
		common.KafkaErrorResponse(w, r, reffID, common.ErrKafkaNotReady)
		return
	}

	common.ReffIDResponse(w, r, reffID)
}
