package db

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"os"
)

var MongoSession *mgo.Session
var mongoErr error
var Instance *mgo.Database

func Init() (*mgo.Database, *mgo.Session) {

	mongoUrl := os.Getenv("MONGO_URL")
	fmt.Println("Mongo URL", mongoUrl)
	MongoSession, mongoErr = mgo.Dial(mongoUrl)
	if mongoErr != nil {
		fmt.Println("Error connecting to mongo", mongoErr)
		panic(mongoErr)
	}
	MongoSession.SetMode(mgo.Monotonic, true)
	Instance = MongoSession.DB("kirk")

	return Instance, MongoSession
}
