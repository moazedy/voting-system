package httpEngine

import (
	"net/http"
	"voting-system/constants"
	"voting-system/domain/models"
	"voting-system/logic"

	"github.com/gin-gonic/gin"
)

type CandidateController interface {
	// AddNewCandidate adds new candidate to candidates list of an election
	AddNewCandidate(ctx *gin.Context)
	// ReadCandidateData gets candidate data of given candidate id
	ReadCandidateData(ctx *gin.Context)
}

type candidate struct {
	Logic logic.CandidateLogic
}

func NewCandidateController(logic logic.CandidateLogic) CandidateController {
	return candidate{Logic: logic}
}

func (c candidate) AddNewCandidate(ctx *gin.Context) {
	// TODO : requester id should be extracted from user claims
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
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, id)
}

func (c candidate) ReadCandidateData(ctx *gin.Context) {
	// TODO : requester id should be extracted from user claims
	requesterId := ""
	candidateId := ctx.Param("candidate_id")

	err := logic.IdValidation(candidateId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	candidateData, err := c.Logic.ReadCandidateData(ctx, candidateId, requesterId, false)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, candidateData)
}
