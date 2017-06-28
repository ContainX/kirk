package db

import (
	"gopkg.in/mgo.v2"
	//"fmt"
	//"gopkg.in/mgo.v2/bson"
	"fmt"
)

var MongoSession *mgo.Session
var mongoErr error
var Instance *mgo.Database

func Init() (*mgo.Database, *mgo.Session) {

	MongoSession, mongoErr = mgo.Dial("mongodb://mongo:27017")
	if mongoErr != nil {
		fmt.Println("Error connecting to mongo", mongoErr)
		panic(mongoErr)
	}
	MongoSession.SetMode(mgo.Monotonic, true)
	Instance = MongoSession.DB("kirk")

	return Instance, MongoSession
}
