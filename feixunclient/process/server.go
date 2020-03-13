package process

import (
	"encoding/binary"
	"encoding/json"
	"feixun/common"
	"feixun/tools"
	"feixunservice/models"
	"fmt"
	"net"
	"os"
	"sync"
)
type UserToOtherUser interface{
	//持续刷新用户在线列表
	goRunContinuousInfo()
	//处理登录后的交互
	TwoLoginIn()
}

type Server struct{
	*sync.Mutex  	     //确保信息收发一致性
	until tools.Utils    //信息发送工具
	UserInfo models.User //当前用户模型
	msg  msgprocess      //消息处理
}

func NewServer(user models.User,conn net.Conn)UserToOtherUser{
	newTools := tools.NewUtils(conn)
	s := Server{
		until:newTools,
		UserInfo:user,
	}
	//初始化消息处理模块
	ms := msgprocess{ut:&s.until,
			id:s.UserInfo.UserId,
				name:s.UserInfo.UserName}
	s.msg = ms
	return s
}

func (this Server)goRunContinuousInfo(){
	//持续刷新用户在线列表等信息
	//监听服务器发送的时时在线用户信息
	for{
		var newMsg common.Msssage
		err := this.until.ResponseMsg(&newMsg)
		if err != nil {
			break
		}
		//newMsg.Data  jsonToInt
		switch newMsg.MsgType {

		//处理用户上线信号
		case common.NOTIFICATION_USERUP:
			OnlineUsers.UpSignalProcess(newMsg.Data)

		//处理用户下线信号
		case common.NOTIFICATION_USERDOWN:
			println("有人下限拉")
			OnlineUsers.DownSignalProcess(newMsg.Data)

		//群发消息接受
		case common.GROUPS_MESSAGES:
			var groupMsg common.UserToUserMsg
			json.Unmarshal(newMsg.Data,&groupMsg)

			fmt.Printf("%s对所有人说：%s\n",groupMsg.Name,groupMsg.Data)
		//发送群消息后的响应结果 是否发送成功
		case common.GROUPS_MSG_RESPONSE:
			var rspMsg common.MsgResMsg
			json.Unmarshal(newMsg.Data,&rspMsg)
			if rspMsg.Code == 200 {
				println("发送成功")
			}else {
				println("发送失败")
			}
		}
	}
}

func (this Server)TwoLoginIn() {
	go this.goRunContinuousInfo()
	for {
		fmt.Println("1.显示用户在线列表")
		fmt.Println("2.发送信息")
		fmt.Println("3.信息列表")
		fmt.Println("4.退出系统")

	var key int
	fmt.Scanln(&key)
	switch key {
		case 1:
			OnlineUsers.PrintUserInfo()
		case 2:
			println("发消息")
			this.msg.AccessToInformation()
		case 3:
			println("信息")
		case 4:
			this.DownSignalRequest()
			os.Exit(0)
	}
	}
}

func (this Server)DownSignalRequest(){
	var i uint32
	i = uint32(this.UserInfo.UserId)
	buf := make([]byte,4)
	binary.BigEndian.PutUint32(buf,i)

	var newMsg common.Msssage
	newMsg.MsgType = common.NOTIFICATION_USERDOWN
	newMsg.Data = buf
	msgJson,err := json.Marshal(newMsg)
	if err != nil{
		println("json出错")
	}
	ut := tools.NewClientUtils()
	ut.RequestMsg(msgJson)


}