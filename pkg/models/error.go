package models

import "time"

type Error struct {
	Details string
	Path    string
	Object  string
	Time    time.Time
}

func (e Error) Error() string {
	return "Error{" +
		"Details:" + e.Details +
		", Path:" + e.Path +
		", Object:" + e.Object +
		"}"
}
