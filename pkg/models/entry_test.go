package models

import (
	"testing"
	"time"
)

func TestEntry_String(t *testing.T) {
	//given
	testEntry := Entry{
		Id:       1,
		MemberId: 2,
		ItemId:   3,
		Capacity: 10,
		Created:  time.Now(),
		Modified: time.Now(),
	}
	expected := "Entry[ID: 1, Member-ID: 2, Item-ID: 3, Capacity: 10]"
	//when
	result := testEntry.String()
	//then
	if result != expected {
		t.Errorf("\nExpected:\n" + expected + "\nBut Got:\n" + result)
	}
}
