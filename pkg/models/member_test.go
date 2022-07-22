package models

import (
	"testing"
	"time"
)

func TestMember_String(t *testing.T) {
	//given
	testMember := Member{
		Id:                1,
		Name:              "Johannes",
		BorrowedItemCount: 15,
		Comment:           "Test",
		Active:            true,
		Created:           time.Now(),
		Modified:          time.Now(),
	}
	expected := "Member[ID: 1, Name: Johannes, Borrowed Items: 15, Comment: Test, Active: true]"
	//when
	result := testMember.String()
	//then
	if result != expected {
		t.Errorf("\nExpected:\n" + expected + "\nBut Got:\n" + result)
	}
}
