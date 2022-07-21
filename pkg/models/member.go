package models

import (
	"strconv"
	"time"
)

type NewMemberInfos struct {
	Name    string
	Comment string
}

type Member struct {
	Id                int
	Name              string
	BorrowedItemCount int
	Comment           string
	Active            bool
	Created           time.Time
	Modified          time.Time
}

type Fundus interface {
	String() string
}

func (m Member) String() string {
	return "Member[ID: " + strconv.Itoa(m.Id) +
		", Name: " + m.Name +
		", Borrowed Items: " + strconv.Itoa(m.BorrowedItemCount) +
		", Comment: " + m.Comment +
		", Active: " + strconv.FormatBool(m.Active) +
		"]"
}
