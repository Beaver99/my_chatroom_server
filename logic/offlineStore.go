package logic

import (
	"gopkg.in/mgo.v2"
)

type OfflineMsg struct{
	UserID UUID
	Msg map[string]interface{}
}

func InitOfflineStore() *mgo.Session{
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)

	//OfflineMessage := session.DB("OfflineMessage").C("to be sent")
	return session
}
