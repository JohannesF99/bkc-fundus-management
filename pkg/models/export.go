package models

import "time"

type Export struct {
	Id       int
	Name     string
	Capacity int
	Date     time.Time
}
