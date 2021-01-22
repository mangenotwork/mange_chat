package constname

const (
	//一对一聊天存入Redis List, obo:<房间token>
	OnebyOneMsgKey = "obo:%s"

	//用哈希 记录来自谁的未读消息  unmsg:<用户>
	UnreadMsgKey = "unmsg:%s"

	//用户信息  u:<用户>
	UserInfoKey = "u:%s"

	//用户列表 一个大的zset,  alluser
	AllUserKey = "alluser"
)
