package httpEngine

import "voting-system/logic"

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

}
