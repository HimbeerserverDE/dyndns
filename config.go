package main

import (
	"encoding/json"
	"os"
	"time"
)

type config struct {
	User      string
	Passwd    string
	Link4     string
	Link6     string
	Interval  time.Duration
	PrefixLen int
	Records4  []int
	Records6  []int
}

func (c *config) parse(location string) error {
	f, err := os.Open(location)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(c)
}
