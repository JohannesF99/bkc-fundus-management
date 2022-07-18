package models

import (
	"github.com/google/uuid"
	"time"
)

type item struct {
	id           uuid.UUID
	name         string
	maxCapacity  int
	availability int
	description  string
	created      time.Time
	modified     time.Time
}
