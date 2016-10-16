package main

import (
	"log"

	"fmt"
	"os/exec"

	"net"

	"github.com/qhsong/golang-vm-ci/common"
	"github.com/rgbkrk/libvirt-go"
	mgo "gopkg.in/mgo.v2"
)

func generatorTasks(connection libvirt.VirConnection, db *mgo.Database, task common.Task) error {
	args := fmt.Sprintf("--name %s -u %s", task.Name, task.UUID)
	cmd := exec.Command("virt-clone", args)
	err := cmd.Run()
	if err != nil {
		task.Status = common.TaskCreatFailed
	} else {
		task.Status = common.TaskPreparingStatus
	}
	db.C(common.C_TASK).UpdateId(task.ID, task)
	return err
}

func getMachineInternalIP(conn libvirt.VirConnection, name string) string {
	dom, _ := conn.LookupDomainByName(name)
	infs, _ := dom.ListAllInterfaceAddresses(0)
	for _, inf := range infs {
		for _, addr := range inf.Addrs {
			localIP := net.ParseIP(addr.Addr)
			_, IPNet, _ := net.ParseCIDR("192.168.122.0/24")
			if IPNet.Contains(localIP) {
				return addr.Addr
			}
		}
	}
	return ""
}

func taskWatchLoop(session *mgo.Session, tasksChan chan []common.Task, exitChan chan int) {
	cDB := session.Clone().DB("")
	conn, err := libvirt.NewVirConnection("qemu:///system")
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case tasks := <-tasksChan:
			for _, task := range tasks {
				go func(conn libvirt.VirConnection, db *mgo.Database, task common.Task) {
					err := generatorTasks(conn, db, task)
					if err != nil {
						log.Println("Generate Machine Failed:", err)
					} else {
						ip := getMachineInternalIP(conn, task.Name)
						if ip == "" {
							log.Fatal("Can not got internal IP address")
						}
						task.IP = ip
						cDB.C(common.C_TASK).UpdateId(task.ID, task)

					}
				}(conn, session.Clone().DB(""), task)
			}
		case <-exitChan:
			return
		}
	}
}
