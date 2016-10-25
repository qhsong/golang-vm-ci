package common

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	TaskNewStatus = iota + 1
	TaskPreparingStatus
	TaskWaitRunningStatus
	TaskRunningStatus
	TaskSuccessStatus
	TaskCreatFailed
)

type Task struct {
	ID       bson.ObjectId `bson:"_id"`
	CommitID string
	Status   int
	Folder   string //Generally they are same as folder
	VMachine
}

type VMachine struct {
	Name     string
	UUID     string
	IP       string
	DiskPath string
}
type Log struct {
	ID       bson.ObjectId
	TaskID   bson.ObjectId
	FileName string
	URL      string
}

type Console struct {
	ID      bson.ObjectId
	Console string
}
