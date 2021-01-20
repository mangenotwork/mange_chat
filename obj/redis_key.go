package obj

const (
	//一对一聊天存入Redis List, obo:<房间token>
	OnebyOneMsgKey = "obo:%s"

	//用哈希 记录来自谁的未读消息  unmsg:<用户>
	UnreadMsgKey = "unmsg:%s"
)
