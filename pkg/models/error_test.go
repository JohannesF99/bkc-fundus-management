package models

import (
	"testing"
	"time"
)

func TestError_Error(t *testing.T) {
	//given
	testError := Error{
		Details: "Test1",
		Path:    "Test2",
		Object:  "Test3",
		Time:    time.Now(),
	}
	expected := "Error[Details: Test1, Path: Test2, Object: Test3]"
	//when
	result := testError.Error()
	//then
	if result != expected {
		t.Errorf("\nExpected:\n" + expected + "\nBut Got:\n" + result)
	}
}
