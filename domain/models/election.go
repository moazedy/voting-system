package models

import (
	"time"

	"github.com/google/uuid"
)

type Election struct {
	Id           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	CreationTime time.Time `json:"creation_time"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	HasEnded     bool      `json:"has_ended"`
}
