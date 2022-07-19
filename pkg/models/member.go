package models

import (
	"fmt"
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

func (m Member) Print() {
	fmt.Printf("Member{%s: %d, %s: %s, %s: %d, %s: %s, %s: %t, %s: %s, %s: %s}\n",
		"ID", m.Id,
		"Name", m.Name,
		"Borrowed Items", m.BorrowedItemCount,
		"Comment", m.Comment,
		"Active", m.Active,
		"Created", m.Created.String(),
		"Modified", m.Modified.String())
}
