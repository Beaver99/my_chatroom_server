package logic

// store those global variables who are used by many other modules.
var userAccountDB = initUserInfoRedis()

var userConnRegister userConnMap

// FIXME: session set wrong
// FIXME: global session is left close
var

var (
	OfflineMsgStoreSession = InitOfflineStore()
	OfflineMsgStoreCollection = OfflineMsgStoreSession.DB("OfflineMessage").C("to_be_sent")
)
