package process

import (
	"encoding/binary"
	"encoding/json"
	"feixun/common"
	"feixun/tools"
	"feixunservice/models"
	"net"
)

//维护动态用户信息列表
//每个key都是user的id，每个value代表一个user和相对应的Conn方便即时通讯
//时刻更新用户Userlist
//添加或删除一个key
type UserInfoLister interface{
	Add(user *models.User,conn net.Conn) //添加用户到在线用户列表
	Delete(id int)						 //从在线用户列表删除一个用户
	ReturnMap()map[int]*UserAndConn		 //返回用户在线列表
	GetJsonUserList()[]byte				 //返回[]byte形式的在线用户列表
	SendUserUpSignal(user *models.User)  //向客户端发送用户上线信号
	SendUserDownSignal(i int)			 //向客户端发送X用户下线信号
}

type UserInfoList struct{
	UserList map[int]*UserAndConn
	until	tools.Utils
}

type UserAndConn struct{
	User OnlineUser
	Conn net.Conn
}

type OnlineUser struct{
	UserId int `json:"user_id"`
	UserName string `json:"user_name"`
}

func NewUserInfoList ()UserInfoLister{
	userlist := make(map[int]*UserAndConn,1000)
	var NewUserInfoList UserInfoList
	NewUserInfoList.UserList = userlist
	NewUserInfoList.until = tools.NewUtils(nil)
	return &NewUserInfoList
}

func(p *UserInfoList)Add(user *models.User,conn net.Conn){
	var newUserAndConn UserAndConn
	var newOnlineUser OnlineUser
	newOnlineUser.UserId = (*user).UserId
	newOnlineUser.UserName = (*user).UserName
	newUserAndConn.User = newOnlineUser
	newUserAndConn.Conn = conn

	p.UserList[(*user).UserId] = &newUserAndConn
}
func(p *UserInfoList)Delete(id int){
	delete(p.UserList,id)
}
func(p *UserInfoList)ReturnMap()map[int]*UserAndConn{
	return p.UserList
}
//向客户端发送XX用户上线信号
func(p *UserInfoList)SendUserUpSignal(user *models.User){
	var newOnlineUser OnlineUser
	newOnlineUser.UserName = user.UserName
	newOnlineUser.UserId = user.UserId

	var NewMsg common.Msssage
	NewMsg.MsgType = common.NOTIFICATION_USERUP
	jsondata,err := json.Marshal(newOnlineUser)
	if err != nil{
		println("json出错")
	}
	NewMsg.Data = jsondata
	msgJson,err := json.Marshal(NewMsg)
	println(string(msgJson))
	if err != nil{
		println("json失败")
	}
	p.RequestSignalToClient(msgJson)

}

func(p  UserInfoList)GetJsonUserList()[]byte{
	userlistMap := make(map[int]OnlineUser,1000)
	onlineList := GlobalOnlineUserList.ReturnMap()
	for k,v := range onlineList {
		userlistMap[k] = v.User
	}
	jsonUserListMap,err := json.Marshal(userlistMap)
	if err != nil{
		println("json出错")
	}
	return jsonUserListMap
}

func(p *UserInfoList)RequestSignalToClient(data []byte){
	for _,c := range p.UserList {
		conn := c.Conn
		p.until.Conn = conn

		p.until.RequestMsg(data)
		}

}

func(p *UserInfoList)SendUserDownSignal(i int){
	//删除要下线的客户端
	delete(p.UserList,i)
	//通知其他客户端删除要下线的客户端id
	buf := make([]byte,4)
	binary.BigEndian.PutUint32(buf,uint32(i))

	var newMsg common.Msssage
	newMsg.MsgType = common.NOTIFICATION_USERDOWN
	newMsg.Data = buf

	msgJson,err := json.Marshal(newMsg)
	if err != nil {
		println("json失败")
	}

	for _,v := range p.ReturnMap(){
		p.until.Conn = v.Conn
		p.until.RequestMsg(msgJson)
	}


}