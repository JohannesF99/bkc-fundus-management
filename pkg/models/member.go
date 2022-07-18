package models

import (
	"database/sql"
	"fmt"
	"time"
)

type ApiMember struct {
	Name    string
	Comment string
}

type Member struct {
	Id                int
	Name              string
	BorrowedItemCount int
	Comment           sql.NullString
	Active            bool
	Created           time.Time
	Modified          time.Time
}

func (m Member) Print() {
	fmt.Printf("Member{%s: %d, %s: %s, %s: %d, %s: %s, %s: %t, %s: %s, %s: %s}\n",
		"ID", m.Id,
		"Name", m.Name,
		"Borrowed Items", m.BorrowedItemCount,
		"Comment", m.Comment.String,
		"Active", m.Active,
		"Created", m.Created.String(),
		"Modified", m.Modified.String())
}

func (m ApiMember) ToNormalMember() *Member {
	return &Member{
		Id:                0,
		Name:              m.Name,
		BorrowedItemCount: 0,
		Comment: sql.NullString{
			String: m.Comment,
			Valid:  true,
		},
		Active:   true,
		Created:  time.Time{},
		Modified: time.Time{},
	}
}
