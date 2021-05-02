package logic

// store those global variables who are used by many other modules.
var userAccountDB = initUserInfoRedis()

var userConnRegister userConnMap

var (
	OfflineMsgStoreSession = InitOfflineStore()
	OfflineMsgStoreCollection = OfflineMsgStoreSession.DB("OfflineMessage").C("to_be_sent")
)
