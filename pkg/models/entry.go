package models

import (
	"strconv"
	"time"
)

type Entry struct {
	Id       int
	MemberId int
	ItemId   int
	Capacity int
	Created  time.Time
	Modified time.Time
}

type NewEntryInfos struct {
	MemberId int
	ItemId   int
	Capacity int
}

func (e Entry) String() string {
	return "Entry[" +
		"ID: " + strconv.Itoa(e.Id) +
		", Member-ID: " + strconv.Itoa(e.MemberId) +
		", Item-ID: " + strconv.Itoa(e.ItemId) +
		", Capacity: " + strconv.Itoa(e.Capacity) +
		"]"
}
