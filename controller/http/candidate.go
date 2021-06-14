package httpEngine

import (
	"net/http"
	"voting-system/constants"
	"voting-system/domain/models"
	"voting-system/logic"

	"github.com/gin-gonic/gin"
)

type CandidateController interface {
	AddNewCandidate(ctx *gin.Context)
}

type candidate struct {
	Logic logic.CandidateLogic
}

func NewCandidateController(logic logic.CandidateLogic) CandidateController {
	return candidate{Logic: logic}
}

func (c candidate) AddNewCandidate(ctx *gin.Context) {
	requesterId := ""
	var candidateData models.Candidate
	err := ctx.BindJSON(&candidateData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": constants.InvalidCandidateData,
		})
		return
	}

	id, err := c.Logic.CreateNewCandidate(ctx, requesterId, candidateData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": constants.InternalServerError,
		})
		return
	}

	ctx.JSON(http.StatusCreated, id)
}
