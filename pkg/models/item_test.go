package models

import (
	"testing"
	"time"
)

func TestItem_String(t *testing.T) {
	//given
	testItem := Item{
		Id:           1,
		Name:         "Johannes",
		Capacity:     15,
		Availability: 15,
		Description:  "Test",
		Created:      time.Now(),
		Modified:     time.Now(),
	}
	expected := "Item[ID: 1, Name: Johannes, Capacity: 15, Available: 15, Description: Test]"
	//when
	result := testItem.String()
	//then
	if result != expected {
		t.Errorf("\nExpected:\n" + expected + "\nBut Got:\n" + result)
	}
}
