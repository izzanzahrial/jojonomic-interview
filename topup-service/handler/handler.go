package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/izzanzahrial/jojonomic-interview/common"
	"github.com/izzanzahrial/jojonomic-interview/entities"
	"github.com/izzanzahrial/jojonomic-interview/topup-service/usecase"
	"github.com/teris-io/shortid"
)

type Topup struct {
	service *usecase.Topup
	sID     *shortid.Shortid
}

func NewTopup(service *usecase.Topup, sid *shortid.Shortid) *Topup {
	return &Topup{
		service: service,
		sID:     sid,
	}
}

func (t *Topup) Publish(w http.ResponseWriter, r *http.Request) {
	var Input struct {
		Gram  string `json:"gram"`
		Harga string `json:"harga"`
		Norek string `json:"norek"`
	}

	if err := common.ReadJSON(w, r, &Input); err != nil {
		common.BadRequestResponse(w, r, err)
		return
	}

	if err := common.CheckGram(Input.Gram); err != nil {
		common.BadRequestResponse(w, r, err)
		return
	}

	harga, err := strconv.ParseInt(Input.Harga, 0, 32)
	if err != nil {
		common.BadRequestResponse(w, r, err)
	}

	gram, err := strconv.ParseFloat(Input.Gram, 32)
	if err != nil {
		common.BadRequestResponse(w, r, err)
	}

	topup := entities.NewTopup(Input.Norek, int32(harga), float32(gram))

	if err := topup.Validate(); len(err) != 0 {
		common.FailedValidationResponse(w, r, err)
		return
	}

	reffID, err := t.sID.Generate()
	if err != nil {
		common.GenerateReffIDErrorResponse(w, r, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = t.service.Publish(ctx, reffID, topup)
	if err != nil {
		common.KafkaErrorResponse(w, r, reffID, common.ErrKafkaNotReady)
		return
	}

	common.ReffIDResponse(w, r, reffID)
}
