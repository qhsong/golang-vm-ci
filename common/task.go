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
	TaskFailedStatus
)

type Task struct {
	ID          bson.ObjectId
	CommitID    string
	Status      int
	MachineName string
	MachineIP   string
	Folder      string //Generally they are same as folder
	Name        string
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
