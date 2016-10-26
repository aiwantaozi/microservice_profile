package configs

import (
	"gopkg.in/mgo.v2"
)

var GlobalMgoSession *mgo.Session

func InitMongo() *mgo.Session {
	var err error
	GlobalMgoSession, err = mgo.Dial(AppConfig.MONGO_PRIMARY_URL)
	if err != nil {
		panic(err)
	}
	return GlobalMgoSession
}

func GetDatabase() *mgo.Database {
	session := GlobalMgoSession.Clone()
	MongoMainDatabase := session.DB(AppConfig.DEFAULT_DATABASE_NAME)
	return MongoMainDatabase
}
