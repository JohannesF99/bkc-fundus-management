package models

import "time"

type Export struct {
	Id       int
	Name     string
	Capacity int
	Date     time.Time
}

type ExpandedExport struct {
	MemberId     int
	MemberName   string
	ItemId       int
	ItemName     string
	ItemCapacity int
	EntryDate    time.Time
}
