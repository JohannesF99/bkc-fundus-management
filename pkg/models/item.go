package models

import (
	"fmt"
	"time"
)

type Item struct {
	Id           int
	Name         string
	Capacity     int
	Availability int
	Description  string
	Created      time.Time
	Modified     time.Time
}

type NewItemInfos struct {
	Name        string
	Capacity    int
	Description string
}

func (i Item) Print() {
	fmt.Printf("Item{%s: %d, %s: %s, %s: %d, %s: %d, %s: %s, %s: %s, %s: %s}\n",
		"ID", i.Id,
		"Name", i.Name,
		"Capacity", i.Capacity,
		"Available", i.Availability,
		"Description", i.Description,
		"Created", i.Created.String(),
		"Modified", i.Modified.String())
}
