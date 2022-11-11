package main

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

var ErrInsecure = errors.New("config has insecure permissions (need 0?00)")

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
	info, err := os.Stat(location)
	if err != nil {
		return err
	}

	if info.Mode().Perm()&0077 > 0 {
		return ErrInsecure
	}

	f, err := os.Open(location)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(c)
}
