package models

import (
	"fmt"
	"time"
)

type EntryInfo struct {
	Id       int
	MemberId int
	ItemId   int
	Capacity int
	Created  time.Time
	Modified time.Time
}

type NewEntry struct {
	MemberId int
	ItemId   int
	Capacity int
}

func (e EntryInfo) Print() {
	fmt.Printf("Entry{%s: %d, %s: %d, %s: %d, %s: %d, %s: %s, %s: %s}\n",
		"ID", e.Id,
		"Member-ID", e.MemberId,
		"Item-ID", e.ItemId,
		"Capacity", e.Capacity,
		"Created", e.Created.String(),
		"Modified", e.Modified.String())
}
