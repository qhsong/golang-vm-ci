package common

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	CommitNewStatus = iota + 1
	CommitPendingStatus
	CommitSuccessStatus
	CommitFailedStatus
)

//Push Info
type Push struct {
	Ref        string
	Before     string
	After      string
	Created    bool
	Delete     bool
	Forced     bool
	Commits    []*Commit
	HeadCommit *Commit `json:"head_commit"`
	Repo       Repository
}

//Commit struct
type Commit struct {
	ID        bson.ObjectId
	CommitID  string `json:"id"`
	Message   string
	Timestamp time.Time `json:"timestamp, string"`
	Added     []string
	Removed   []string
	Modified  []string
	URL       string
	Status    int
}

//Repository struct
type Repository struct {
	ID       string
	Name     string
	FullName string `json:"full_name"`
	URL      string `json:"url"`
	GitURL   string `json:"git_url"`
	SSHURL   string `json:"ssh_url"`
	CloneURL string `json:"clone_url"`
	Size     int
}
