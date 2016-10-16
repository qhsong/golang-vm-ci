package main

import (
	"log"
	"net/http"

	mgo "gopkg.in/mgo.v2"

	"github.com/gin-gonic/gin"
	"github.com/qhsong/golang-vm-ci/common"
)

func postTriggerHandler(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)
	method := c.Request.Header["X-Github-Event"][0]
	tasksChan := c.MustGet("tasksChan").(chan []common.Task)
	switch method {
	case "ping":
		var tmp interface{}
		if c.Bind(&tmp) == nil {
			err := db.C(common.C_GITHUB_INFO).Insert(tmp)
			if err != nil {
				log.Fatal("Insert ping info error!" + err.Error())
			}
			c.String(http.StatusOK, "Config Success")

		} else {
			log.Fatal("Error to decoder ping result!")
		}
	case "push":
		var PushInfo common.Push
		if c.BindJSON(&PushInfo) == nil {
			err := triggerPush(PushInfo, db, tasksChan)
			if err != nil {
				c.String(http.StatusInternalServerError, "Intern Server Error!")
			}
		} else {
			c.String(http.StatusNotFound, "Can not generator JSON String")
		}
	}
}
