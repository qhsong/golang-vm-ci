package main

import (
	"github.com/qhsong/golang-vm-ci/common"
	mgo "gopkg.in/mgo.v2"
)

func generatorTasks(tasks []common.Task, db *mgo.Session) error {

	return nil
}

func taskWatchLoop(db *mgo.Session, tasksChan chan []common.Task) {
	for {
		select {
		case tasks := <-tasksChan:
			cdb := db.Clone()

		}
	}
}
