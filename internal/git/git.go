package git

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Runner interface {
	Run(args ...string) (string, error)
}

type DefaultRunner struct{}

func (r *DefaultRunner) Run(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

type Commit struct {
	Hash      string
	Subject   string
	Author    string
	Timestamp time.Time
	Prefix    string
}

var validPrefixes = map[string]bool{
	"feat":     true,
	"fix":      true,
	"docs":     true,
	"chore":    true,
	"refactor": true,
	"test":     true,
	"style":    true,
	"perf":     true,
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

func ExtractPrefix(subject string) string {
	colonIndex := strings.Index(subject, ":")
	if colonIndex == -1 {
		return "other"
	}

	prefix := subject[:colonIndex]

	parenIndex := strings.Index(prefix, "(")
	if parenIndex != -1 {
		prefix = prefix[:parenIndex]
	}

	if validPrefixes[prefix] {
		return prefix
	}

	return "other"
}
