package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/izzanzahrial/jojonomic-interview/buyback-service/usecase"
	"github.com/izzanzahrial/jojonomic-interview/common"
	"github.com/izzanzahrial/jojonomic-interview/entities"
	"github.com/teris-io/shortid"
)

type Buyback struct {
	service *usecase.Buyback
	sID     *shortid.Shortid
}

func NewBuyback(service *usecase.Buyback, sid *shortid.Shortid) *Buyback {
	return &Buyback{
		service: service,
		sID:     sid,
	}
}

func (b *Buyback) Publish(w http.ResponseWriter, r *http.Request) {
	var Input struct {
		Gram  float32 `json:"gram"`
		Harga int32   `json:"harga"`
		Norek string  `json:"norek"`
	}

	if err := common.ReadJSON(w, r, &Input); err != nil {
		common.BadRequestResponse(w, r, err)
		return
	}

	gram := strconv.FormatFloat(float64(Input.Gram), 'f', -1, 32)

	if err := common.CheckGram(gram); err != nil {
		common.BadRequestResponse(w, r, err)
		return
	}

	gramFloat, err := strconv.ParseFloat(gram, 32)
	if err != nil {
		common.BadRequestResponse(w, r, err)
	}

	buyback := entities.NewBuyback(Input.Norek, Input.Harga, float32(gramFloat))

	if err := buyback.Validate(); len(err) != 0 {
		common.FailedValidationResponse(w, r, err)
		return
	}

	reffID, err := b.sID.Generate()
	if err != nil {
		common.GenerateReffIDErrorResponse(w, r, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = b.service.Publish(ctx, reffID, buyback)
	if err != nil {
		common.KafkaErrorResponse(w, r, reffID, common.ErrKafkaNotReady)
		return
	}

	common.ReffIDResponse(w, r, reffID)
}
