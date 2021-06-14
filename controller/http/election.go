package httpEngine

import (
	"net/http"
	"voting-system/constants"
	"voting-system/domain/models"
	"voting-system/logic"

	"github.com/gin-gonic/gin"
)

// ElectionController is interface of election entity in controller  layer
type ElectionController interface {
	// CreateElection gets election data from client and passes it to logic layer for creating a new election
	CreateNewElection(c *gin.Context)
}

// election is a struct ot hold controller methods for election entity in controller
type election struct {
	Logic logic.ElectionLogic
}

// NewElectionController is construction function for ElectionController
func NewElectionController(logic logic.ElectionLogic) ElectionController {
	return election{Logic: logic}
}

// TODO : error handling should be done, using standard package
func (e election) CreateNewElection(c *gin.Context) {
	// extracting requester id
	requesterId := "" // TODO : need to be extracted from user claimes
	var electionData models.Election
	err := c.BindJSON(&electionData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": constants.InvalidElectionData,
		})
		return
	}

	// passing data to logic layer
	id, err := e.Logic.CreateNewElection(c, requesterId, electionData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// returning created election's id to http client
	c.JSON(http.StatusOK, id)
}
