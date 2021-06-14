package logic

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
)

type ElectionManager interface {
	// ManageElctions is a method for ElectionManager worker to terminate expired elections and save running
	// elections result, in db
	ManageElections(ctx context.Context) error
	// Run is a method for ElectionManager to be executed periodically in a certain time
	Run(periodTime int)
}

type electionManager struct {
	electionLogic ElectionLogic
}

func NewElectionManager() ElectionManager {
	return new(electionManager)
}

func (e *electionManager) ManageElections(ctx context.Context) error {
	// singlton design pattern ...
	if e.electionLogic == nil {
		e.electionLogic = NewElectionLogic()
	}

	// getting list of running elections
	elections, err := e.electionLogic.GetListOfStartedElections(ctx, true)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(2 * len(elections))
	// terminating expired elections and saving elction results in db
	for _, v := range elections {
		// a concurrent func for termination
		go func() {
			if time.Now().After(v.EndTime) {
				v.HasEnded = true
			}
			wg.Done()
		}()

		// another func to saving results
		go func() {
			_, err := e.electionLogic.CalculationElectionResults(ctx, v.Id.String(), "", true)
			if err != nil {
				log.Println("error in calculating "+v.Title+" election with id: "+v.Id.String(), err.Error())
			}

			wg.Done()
		}()
	}

	return nil
}

func (e *electionManager) Run(periodTime int) {
	// using gocron package to create an object like a cronjob manager
	s := gocron.NewScheduler(time.UTC)

	// manager runs every (periodTime) * seconds, the passed function
	s.Every(periodTime).Seconds().Do(e.ManageElections(context.Background()))
	s.StartAsync()
}
