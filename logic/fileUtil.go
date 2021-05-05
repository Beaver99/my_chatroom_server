package logic

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/GGBooy/message"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"strconv"
)

//func fileExist(path string) (bool, error) {
//	_, err := os.Stat(path)
//	if err == nil {
//		return true, nil
//	}
//	if os.IsNotExist(err) {
//		return false, nil
//	}
//	return false, err
//}
type fileID struct {
	File_id bson.ObjectId `bson:"file_id"`
}

// 发送文件
func SendFile(ctx context.Context, ReceiverID UUID, c *websocket.Conn, msg map[string]interface{}) {
	filename := msg["Filename"].(string)
	sendername := msg["Sendername"].(string)
	offsetStr := msg["Offset"].(string)
	offsetInt, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}

	session, err := mgo.Dial("mongodb://127.0.0.1:27017")
	defer session.Close()
	// get triple
	var fileIDs fileID
	//TODO: send a dup file
	query := bson.M{"sender_id": UUID(sendername), "receiver_id": ReceiverID, "filename": filename}
	//SingleChatFileStoreCollection.Find(query).Select(bson.M{"file_id": 1, "_id": 0}).All(&fileIDs)

	SingleChatFileStoreCollection.Find(query).Select(bson.M{"file_id": 1, "_id": 0}).One(&fileIDs)

	//log.Println(filename,"Here are the fileIDs that are found:", fileIDs,reflect.TypeOf(fileIDs))
	//for i := range fileIDs {
	//	buff, _ := bson.Marshal(fileIDs)
	//	var fileIdMap fileID
	//	bson.Unmarshal(buff, &fileIdMap)
	//	fmt.Println(fileIdMap,reflect.TypeOf(fileIdMap))
	f, err := session.DB("gridfs").GridFS("fs").OpenId(fileIDs.File_id)
	if err != nil {
		fmt.Println("can't open this file")
		log.Println(err)
		return
	}
	defer f.Close()
	SendSeg(ctx, c, f, sendername, offsetInt)
	//}
	//f, err := session.DB("gridfs").GridFS("fs").FindID(fileIDs[i]).Open(filename)
	//if err != nil {
	//	fmt.Println("can't open this file")
	//	log.Println(err)
	//	return
	//}
	//defer f.Close()
	//SendSeg(ctx, c, f, sendername, offsetInt)
}

func SendSeg(ctx context.Context, c *websocket.Conn, f *mgo.GridFile, sendername string, offsetInt int64) {
	// 从指定的Offset处发送之后的数据
	filename := f.Name()
	fileLenInt := f.Size()
	fileLenStr := strconv.FormatInt(fileLenInt, 10)
	buffer := make([]byte, 4096)

	for {
		// 偏移至offset位置后发送一段数据
		_, _ = f.Seek(offsetInt, io.SeekStart)
		num, err := f.Read(buffer)
		if err == io.EOF {
			// 如果已经到达文件结尾，停止发送
			fmt.Println("send over")
			return
		} else if err != nil {
			log.Println(err)
			return
		}

		offsetStr := strconv.FormatInt(offsetInt, 10)
		fileseg := message.DataMessage{
			MessageType: "6",
			Sendername:  sendername,
			Filename:    filename,
			Length:      fileLenStr,
			Offset:      offsetStr,
			Data:        buffer[:num],
		}
		err = wsjson.Write(ctx, c, fileseg)
		if err != nil {
			return
		}
		offsetInt = offsetInt + int64(num)
	}
}

func RecvSeg(ctx context.Context, ReceiverID UUID, c *websocket.Conn, msg map[string]interface{}) {
	sendername := msg["Sendername"].(string)
	filename := msg["Filename"].(string)
	session, err := mgo.Dial("mongodb://127.0.0.1:27017")
	if err != nil {
		log.Println(err)
		return
	}
	defer session.Close()
	//f, err := session.DB("gridfs").GridFS("fs").Open(filename)
	//if err == mgo.ErrNotFound {
	//	f, _ = session.DB("gridfs").GridFS("fs").Create(filename)
	//} else if err != nil {
	//	log.Println(err)
	//	return
	//}
	// ##
	f, _ := session.DB("gridfs").GridFS("fs").Create(filename)
	defer f.Close()

	for {
		bufferStr := msg["Data"].(string)
		buffer, err := base64.StdEncoding.DecodeString(bufferStr)
		if err != nil {
			log.Println(err)
		}

		fileLenStr := msg["Length"].(string)
		fileLenInt, err := strconv.ParseInt(fileLenStr, 10, 64)
		if err != nil {
			log.Println(err)
			return
		}

		offsetStr := msg["Offset"].(string)
		offsetInt, err := strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			log.Println(err)
			return
		}

		// 在Offset处(即末尾)写入数据
		_, _ = f.Seek(offsetInt, io.SeekStart)
		_, err = f.Write(buffer)
		if err != nil {
			log.Println(err)
		}

		// 获取写入数据后末尾偏移，判断是否接受完成
		posInt, err := f.Seek(0, io.SeekEnd)
		if err != nil {
			log.Println(err)
		}
		if posInt == fileLenInt {
			fmt.Println("receive over")
			SendNotify(ctx, c, sendername, filename, fileLenStr)
			SingleChatFileTripleInsert(SingleChatFileTriple{
				SenderID:   UUID(sendername),
				ReceiverID: ReceiverID,
				Filename:   filename,
				FileID:     f.Id(),
			})
			return
		}
		msg = <-chFile
	}
}

func SendNotify(ctx context.Context, c *websocket.Conn, sendername string, filename string, lengthStr string) {
	notifyMsg := message.FileNotifyMessage{
		MessageType: "3",
		Sendername:  sendername,
		Filename:    filename,
		Length:      lengthStr,
	}
	wsjson.Write(ctx, c, notifyMsg)
}

func WriteSeg(ctx context.Context, c *websocket.Conn, buffer []byte, offsetInt int64) {

}
