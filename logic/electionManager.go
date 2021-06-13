package logic

import "context"

type ElectionManager interface {
	ManageElections(ctx context.Context) error
}

type electionManager struct {
}

func NewElectionManager() ElectionManager {
	return new(electionManager)
}

func (e electionManager) ManageElections(ctx context.Context) error {
	// TODO : implementation
	return nil
}
