package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type BadRequestError string

type DateTime struct {
	time time.Time
}

func (dt *DateTime) format() string {
	return time.RFC3339
}
func (dt *DateTime) String() string {
	return dt.time.Format(dt.format())
}
func (dt *DateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(dt.String())
}
func (dt *DateTime) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		panic(BadRequestError("not a valid string"))
	}
	if s == "" {
		panic(BadRequestError("must not be empty"))
	}
	t, err := time.Parse(dt.format(), s)
	if err != nil {
		panic(BadRequestError("format must be YYYY-MM-DDTHH:mm:ssZ"))
	}
	dt.time = t
	return nil
}

type Task struct {
	Id          uint64
	Content     string
	DateCreated DateTime
	Completed   bool
}

func (t Task) String() string {
	done := "done"
	if !t.Completed {
		done = "not done"
	}
	return fmt.Sprintf("%d - %s - %s - %s", t.Id, t.Content, t.DateCreated.String(), done)
}
