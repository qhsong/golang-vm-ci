package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/qhsong/golang-vm-ci/common"
	mgo "gopkg.in/mgo.v2"
)

//DB struct
type DB struct {
	Session *mgo.Session
	Mongo   *mgo.DialInfo
}

var dbinfo DB

//MongoDBURL Default MongoDB Connect string
const MongoDBURL = "mongodb://127.0.0.1:27017/vm"

func connectDB() {
	uri := os.Getenv("MONGODB_URL")
	if len(uri) == 0 {
		uri = MongoDBURL
	}
	mongo, err := mgo.ParseURL(uri)
	if err != nil {
		log.Fatal("Can not parse DB" + err.Error())
	}
	s, err := mgo.Dial(uri)
	if err != nil {
		log.Fatal("Can not connect to DB." + err.Error())
	}
	log.Println("Connect to DB successful!")
	dbinfo.Mongo = mongo
	dbinfo.Session = s
}

//For Connect DB Data middileware
func DBMiddileWare(c *gin.Context) {
	s := dbinfo.Session.Clone()
	defer s.Close()

	c.Set("db", s.DB(""))
	c.Next()
}

func ChanMiddileWare(tasksChan chan []common.Task) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("tasksChan", tasksChan)
		c.Next()
	}
}
func main() {
	r := gin.Default()
	connectDB()
	r.Use(DBMiddileWare)

	tasksChan := make(chan []common.Task, 1000)
	r.Use(ChanMiddileWare(tasksChan))
	r.POST("/trigger/", postTriggerHandler)
	r.Run(":8080")
}
