package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/izzanzahrial/jojonomic-interview/check-mutasi-service/repository"
	"github.com/izzanzahrial/jojonomic-interview/common"
	"github.com/teris-io/shortid"
)

type CheckMutation struct {
	repository *repository.CheckMutation
	sID        *shortid.Shortid
}

func NewCheckMutation(repository *repository.CheckMutation, sid *shortid.Shortid) *CheckMutation {
	return &CheckMutation{
		repository: repository,
		sID:        sid,
	}
}

func (cm *CheckMutation) GetListTransactionByDate(w http.ResponseWriter, r *http.Request) {
	var Input struct {
		Norek     string `json:"norek"`
		StartDate int32  `json:"start_date"`
		EndDate   int32  `json:"end_date"`
	}

	if err := common.ReadJSON(w, r, &Input); err != nil {
		common.BadRequestResponse(w, r, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	reffID, err := cm.sID.Generate()
	if err != nil {
		common.GenerateReffIDErrorResponse(w, r, err)
		return
	}

	transactions, err := cm.repository.GetListTransactionByDate(ctx, Input.Norek, Input.StartDate, Input.EndDate)
	if err != nil {
		common.KafkaErrorResponse(w, r, reffID, common.ErrKafkaNotReady)
		return
	}

	common.SuccessResponse(w, r, http.StatusOK, common.Envolope{"error": false, "data": transactions})
}
