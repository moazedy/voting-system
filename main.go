package main

import (
	"voting-system/constants"
	"voting-system/logic"
	"voting-system/repository"
)

func main() {
	repository.Init()
	electionManager := logic.NewElectionManager()
	electionManager.Run(constants.ElectionManagerWorkerWorkPeriod)
}
