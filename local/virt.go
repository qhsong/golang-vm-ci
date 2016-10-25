package main

import (
	"log"

	"os/exec"
	"time"

	"bytes"
	"errors"
	"io/ioutil"
	"net"
	"text/template"

	"github.com/qhsong/golang-vm-ci/common"
	"github.com/rgbkrk/libvirt-go"
	mgo "gopkg.in/mgo.v2"
)

func generatorTasks(connection *libvirt.VirConnection, db *mgo.Database, maTmpl *template.Template, task *common.Task) error {
	cmd := exec.Command("qemu-img", "create", "-b", "/root/rhel6.8.1.qcow2", "-f", "qcow2", task.DiskPath)
	//cmd.Args = strings.Split(args, " ")
	var out bytes.Buffer
	var errbuf bytes.Buffer
	cmd.Stderr = &errbuf
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		task.Status = common.TaskCreatFailed
		err = errors.New("\nCmd err: " + errbuf.String() + "\nstdout: " + out.String() + "\n")
	} else {
		var xmlBuffer bytes.Buffer
		err = maTmpl.Execute(&xmlBuffer, task)
		if err != nil {
			task.Status = common.TaskCreatFailed
		}
		var taskDomain libvirt.VirDomain
		taskDomain, err = connection.DomainCreateXML(xmlBuffer.String(), libvirt.VIR_DOMAIN_START_PAUSED)
		if err != nil {
			task.Status = common.TaskCreatFailed
			taskDomain.Free()
		} else {
			taskDomain.Resume()
			go getMachineInternalIP(&taskDomain, db, task)
		}

	}
	if task.Status == common.TaskCreatFailed {
		db.C(common.C_TASK).UpdateId(task.ID, task)
		db.Session.Close()
	}
	return err
}

func getMachineInternalIP(dom *libvirt.VirDomain, db *mgo.Database, task *common.Task) {
	timeDur, _ := time.ParseDuration("30s")
	for {
		infs, _ := dom.ListAllInterfaceAddresses(0)
		for _, inf := range infs {
			for _, addr := range inf.Addrs {
				localIP := net.ParseIP(addr.Addr)
				_, IPNet, _ := net.ParseCIDR("192.168.122.0/24")
				if IPNet.Contains(localIP) {
					task.IP = addr.Addr
					db.C(common.C_TASK).UpdateId(task.ID, task)
					db.Session.Close()
					return
				}
			}
		}
		time.Sleep(timeDur)
	}
}

func taskWatchLoop(session *mgo.Session, tasksChan chan []common.Task, exitChan chan int) {
	conn, err := libvirt.NewVirConnection("qemu:///system")
	if err != nil {
		log.Fatal(err)
	}
	tmplByte, err := ioutil.ReadFile("template.xml")
	if err != nil {
		log.Fatal("Error to read template.xml" + err.Error())
	}
	maTmpl, _ := template.New("machineXml").Parse(string(tmplByte))
	for {
		select {
		case tasks := <-tasksChan:
			for _, task := range tasks {
				go func(conn libvirt.VirConnection, db *mgo.Database, task common.Task) {
					err := generatorTasks(&conn, db, maTmpl, &task)
					if err != nil {
						log.Println("Generate Machine Failed:", err)
					}
				}(conn, session.Clone().DB(""), task)
			}
		case <-exitChan:
			return
		}
	}
}
