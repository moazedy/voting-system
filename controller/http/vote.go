package httpEngine

import (
	"net/http"
	"voting-system/constants"
	"voting-system/domain/models"
	"voting-system/logic"

	"github.com/gin-gonic/gin"
)

type VoteController interface {
	// SaveNewVote is a handler function to voting in an election
	SaveNewVote(c *gin.Context)
	// ReadVoteData reads data of a specific vote and returns it to requester
	ReadVoteData(c *gin.Context)
	// DeleteVote deletes data of a specific vote
	DeleteVote(c *gin.Context)
	// UpdateVoteData updates a specific votes data
	UpdateVoteData(c *gin.Context)
	// GetCandidatePositiveVotesCount gets number of positive votes of a candidate
	GetCandidatePositiveVotesCount(c *gin.Context)
}

type vote struct {
	Logic logic.VoteLogic
}

func NewVoteController(logic logic.VoteLogic) VoteController {
	return vote{Logic: logic}
}

func (v vote) SaveNewVote(c *gin.Context) {
	// requesterId should be extracted from access token
	requesterId := ""
	var voteData models.Vote
	err := c.BindJSON(&voteData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": constants.InvalidVoteData,
		})
		return
	}
	id, err := v.Logic.SaveNewVote(c, voteData, requesterId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, id)
}

func (v vote) ReadVoteData(c *gin.Context) {
	// requesterId should be extracted from access token
	requesterId := ""
	voteId := c.Param("vote_id")
	err := logic.IdValidation(voteId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	voteData, err := v.Logic.ReadVoteData(c, voteId, requesterId, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, voteData)
}

func (v vote) DeleteVote(c *gin.Context) {
	// requesterId should be extracted from access token
	requesterId := ""
	voteId := c.Param("vote_id")
	err := logic.IdValidation(voteId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = v.Logic.DeleteVote(c, voteId, requesterId, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (v vote) UpdateVoteData(c *gin.Context) {
	// requesterId should be extracted from access token
	requesterId := ""
	voteId := c.Param("vote_id")
	err := logic.IdValidation(voteId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var voteData models.Vote
	err = c.BindJSON(&voteData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": constants.InvalidVoteData,
		})
		return
	}

	err = v.Logic.UpdateVoteData(c, voteId, requesterId, voteData, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (v vote) GetCandidatePositiveVotesCount(c *gin.Context) {
	// requesterId should be extracted from access token
	requesterId := ""
	candidateId := c.Param("candidate_id")
	err := logic.IdValidation(candidateId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	votes, err := v.Logic.AgregateOfCandidatePositiveVotes(c, candidateId, requesterId, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, votes)
}
