package logic

// store those global variables who are used by many other modules.
var userAccountDB = initUserInfoRedis()

var groupInfoDB = initGroupInfoRedis()

var userConnRegister userConnMap

// FIXME: session set wrong
// FIXME: global session is left close
var MonogoDBSession = InitOfflineMsgStore()

var (
	OfflineMsgStoreCollection = MonogoDBSession.DB("OfflineStore").C("Message")
)

var (
	SingleChatFileStoreCollection = MonogoDBSession.DB("OfflineStore").C("SingleChatFile")
	//GroupChatFileStoreCollection = MonogoDBSession.DB("OfflineStore").C("GroupChatFile")

)

var chFile = make(chan map[string]interface{}, 64)
