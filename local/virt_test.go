package main

import (
	"testing"

	"github.com/qhsong/golang-vm-ci/common"
	"github.com/rgbkrk/libvirt-go"
	"github.com/satori/go.uuid"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func TestGeneratorTasks(t *testing.T) {
	conn, _ := libvirt.NewVirConnection(`qemu:///system`)
	mongo, err := mgo.Dial(MongoDBURL)
	if err != nil {
		t.Error("Can not dial to mongoDB")
		t.Fail()
	}
	db := mongo.DB("test")
	task := common.Task{
		ID:       bson.NewObjectId(),
		CommitID: "testCommitID",
		Status:   common.TaskNewStatus,
		Folder:   "testFolder",
	}
	task.Name = "testName"
	task.UUID = uuid.NewV4().String()

	err = generatorTasks(conn, db, task)
	if err != nil {
		t.Error("Failed to run generatorTasks", err)
		t.Fail()
	}
	dom, err := conn.LookupByUUIDString(task.UUID)
	defer func(err error) {
		if err == nil {
			dom.Undefine()
		}
	}(err)
	defer dom.Free()

	if err != nil {
		t.Error("Can not find the Machine")
	}
	var newTask common.Task
	db.C(common.C_TASK).FindId(task.ID).One(&newTask)
	if newTask.Status != common.TaskPreparingStatus {
		t.Error("Failed to change task status")
		t.Fail()
	}
}
