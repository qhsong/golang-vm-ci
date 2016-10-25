package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"text/template"

	"github.com/qhsong/golang-vm-ci/common"
	"github.com/rgbkrk/libvirt-go"
	"github.com/satori/go.uuid"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func TestGeneratorTasks(t *testing.T) {
	conn, err := libvirt.NewVirConnection(`qemu:///system`)
	if err != nil {
		t.Fatal(err)
	}
	mongo, err := mgo.Dial(MongoDBURL)
	if err != nil {
		t.Fatal("Can not dial to mongoDB")
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
	task.DiskPath = fmt.Sprintf("/root/%s.qcow2", task.Name)

	err = db.C(common.C_TASK).Insert(task)
	if err != nil {
		t.Fatal("insert task error")
	}
	tmplByte, err := ioutil.ReadFile("template.xml")
	if err != nil {
		log.Fatal("Error to read template.xml" + err.Error())
	}
	maTmpl, _ := template.New("machineXml").Parse(string(tmplByte))

	err = generatorTasks(&conn, db.Session.Clone().DB("test"), maTmpl, &task)
	if err != nil {
		t.Error("Failed to run generatorTasks:" + err.Error())
		t.Fail()
	}
	dom, err := conn.LookupByUUIDString(task.UUID)

	if err != nil {
		t.Error("Can not find the Machine")
	}

	getMachineInternalIP(&dom, db.Session.Clone().DB("test"), &task)
	var newTask common.Task
	err = db.C(common.C_TASK).FindId(task.ID).One(&newTask)
	if err != nil {
		t.Error("Can not find IP", err)
	}
	if len(newTask.IP) <= 0 {
		t.Error("Can not got machine IP.")
	}
	log.Println()

	defer os.Remove(task.DiskPath)
	defer dom.Free()
	defer dom.Undefine()
	defer dom.Destroy()
	err = db.C(common.C_TASK).FindId(task.ID).One(&newTask)
	if err != nil {
		t.Fatal("Can not find task in DB")
	}

}
