package common

const(
	LOGINMSG = iota       //登录
	REGISTRATIONMSG       //注册
	NOTIFICATION_USERUP   //用户上线通知
	NOTIFICATION_USERDOWN //用户下线通知
	GROUPS_MESSAGES       //消息群发
	ONETOONE_MESSAGES     //一对一消息发送
	GROUPS_MSG_RESPONSE
)
//消息结构本体
type Msssage struct{
	MsgType int `json:"msg_type"`
	Data 	[]byte `json:"data"`
}

//登录消息
type LogMsg struct{
	Username string `json:"username"`
	Password string `json:"password"`
}
//登录响应消息
type LoginResMsg struct{
	Userinfo string `json:"userinfo"`
	Code int `json:"code"`
	OnlineUserList []byte `json:"online_user_list"`  //在线用户列表
	Err error `json:"err"`
}
//注册消息
type RegistrationMsg struct{
	UserId   int `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
//注册响应消息
type RegistResMsg struct{
	Code int `json:"code"`
	Err error `json:"err"`
}

//用户消息传递
type UserToUserMsg struct{
	Sender int    `json:"sender"` //发送者id
	Data   string `json:"data"`   //发送数据内容
	Name   string `json:"name"`   //发送对象 为空则是群发
}

type MsgResMsg struct{
	Code int `json:"code"`
	Err  string  `json:"err"`
}