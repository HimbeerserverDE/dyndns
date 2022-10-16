package main

import (
	"encoding/json"
	"os"
)

type config struct {
	User      string
	Passwd    string
	Link      string
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
