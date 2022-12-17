package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/izzanzahrial/jojonomic-interview/check-saldo-service/repository"
	"github.com/izzanzahrial/jojonomic-interview/common"
	"github.com/teris-io/shortid"
)

type CheckSaldo struct {
	repository *repository.CheckSaldo
	sID        *shortid.Shortid
}

func NewCheckSaldo(repository *repository.CheckSaldo, sid *shortid.Shortid) *CheckSaldo {
	return &CheckSaldo{
		repository: repository,
		sID:        sid,
	}
}

func (cs *CheckSaldo) GetSaldo(w http.ResponseWriter, r *http.Request) {
	var Input struct {
		Norek string `json:"norek"`
	}

	if err := common.ReadJSON(w, r, &Input); err != nil {
		common.BadRequestResponse(w, r, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	reffID, err := cs.sID.Generate()
	if err != nil {
		common.GenerateReffIDErrorResponse(w, r, err)
		return
	}

	saldo, err := cs.repository.GetSaldo(ctx, Input.Norek)
	if err != nil {
		fmt.Println(err)
		common.KafkaErrorResponse(w, r, reffID, common.ErrKafkaNotReady)
		return
	}

	var data struct {
		Norek string  `json:"norek"`
		Saldo float32 `json:"saldo"`
	}

	data.Norek = Input.Norek
	data.Saldo = saldo

	common.SuccessResponse(w, r, http.StatusOK, common.Envolope{"error": false, "data": data})
}
