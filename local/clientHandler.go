package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qhsong/golang-vm-ci/common"
	mgo "gopkg.in/mgo.v2"
)

func getClientInfoHandler(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)
	queryIP := c.Request.RemoteAddr
	var qTask common.Task
	qTask.IP = queryIP
	var resultTask common.Task
	err := db.C(common.Task).Find(task).One(&resultTask)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Can not find task"})
	}
	c.JSON(http.StatusOK, resultTask)
}
