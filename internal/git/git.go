package git

import "time"

type Commit struct {
	Hash      string
	Subject   string
	Author    string
	Timestamp time.Time
	Prefix    string
}
