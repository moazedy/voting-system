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
	// ReadElectionData reads an election data using received id from client
	ReadElectionData(c *gin.Context)
	// DeleteElection deletes given election
	DeleteElection(c *gin.Context)
	// UpdateElection updates election data
	UpdateElection(c *gin.Context)
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
	c.JSON(http.StatusCreated, id)
}

func (e election) ReadElectionData(c *gin.Context) {
	electionId := c.Param("election_id")
	// id validation
	err := logic.IdValidation(electionId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// passing received id to logic
	electionData, err := e.Logic.ReadElectionData(c, electionId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// returning results to http client
	c.JSON(http.StatusOK, electionData)
}

func (e election) DeleteElection(c *gin.Context) {
	// extracting requester id
	requesterId := "" // TODO : need to be extracted from user claimes
	electionId := c.Param("election_id")

	// validating receieved election id
	err := logic.IdValidation(electionId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// passing received id to logic for deleting process
	err = e.Logic.DeleteElection(c, electionId, requesterId, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (e election) UpdateElection(c *gin.Context) {
	// extracting requester id
	requesterId := "" // TODO : need to be extracted from user claimes

	var newElectionData models.Election
	err := c.BindJSON(&newElectionData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": constants.InvalidElectionData,
		})
		return
	}

	// passing data to logic for updating election
	err = e.Logic.UpdateElection(c, newElectionData, requesterId, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
