package models

import (
	"strconv"
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

func (i Item) String() string {
	return "Item[" +
		"ID: " + strconv.Itoa(i.Id) +
		", Name: " + i.Name +
		", Capacity: " + strconv.Itoa(i.Capacity) +
		", Available: " + strconv.Itoa(i.Availability) +
		", Description: " + i.Description +
		"]"
}
