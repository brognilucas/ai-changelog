package git

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type Commit struct {
	Hash      string
	Subject   string
	Author    string
	Timestamp time.Time
	Prefix    string
}

func ParseCommitLine(line string) (Commit, error) {
	parts := strings.Split(line, "|")
	if len(parts) != 4 {
		return Commit{}, errors.New("invalid commit line format: expected 4 pipe-separated fields")
	}

	timestamp, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		return Commit{}, errors.New("invalid timestamp: must be a unix timestamp")
	}

	return Commit{
		Hash:      parts[0],
		Subject:   parts[1],
		Author:    parts[2],
		Timestamp: time.Unix(timestamp, 0),
	}, nil
}
