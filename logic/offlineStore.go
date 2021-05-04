package logic

import (
	"gopkg.in/mgo.v2"
)

type OfflineMsg struct{
	UserID UUID
	Msg map[string]interface{}
}
type SingleChatFileTriple struct{
	SenderID UUID `bson:"sender_id"`
	ReceiverID UUID `bson:"receiver_id"`
	Filename string `bson:"filename"`
	FileID interface{} `bson:"file_id"`
}

func InitOfflineMsgStore() *mgo.Session{
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)

	//OfflineMessage := session.DB("OfflineMessage").C("to be sent")
	return session
}

func SingleChatFileTripleInsert(record SingleChatFileTriple){
	SingleChatFileStoreCollection.Insert(record)
}
