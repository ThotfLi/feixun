package process

import (
	"encoding/json"
	"feixun/client"
	"feixun/common"
	"feixun/tools"
	"fmt"
	"net"
)

//处理用户登录、注册等
type UserProcesser interface {
	//处理用户登录请求
	Login(username string,userpwd string) common.LoginResMsg
	//处理用户注册请求
	Registration(userid int,username string,userpwd string)common.RegistResMsg
	Close()
	ReturnClint() net.Conn
}

type UserProcess struct{
	//
	conn net.Conn
	ut tools.Utils
}

//登录 输入用户名密码
//给服务器传送消息
//接受服务器响应消息
func NewUserProcess()UserProcesser{
	NewUProcess := UserProcess{}
	conn := client.NewClient()
	NewUProcess.conn = conn
	NewUProcess.ut = tools.Utils{Conn:conn}
	return &NewUProcess
}

func (this *UserProcess)Login(username string,userpwd string)common.LoginResMsg{
		//封装登录消息
		var loginmessage common.LogMsg
		loginmessage.Username = username
		loginmessage.Password = userpwd

		var messages common.Msssage
		messages.MsgType = common.LOGINMSG
		messages.Data,_ = json.Marshal(loginmessage)

		msgJson,_:=json.Marshal(messages)

		//接受登录是否成功消息
		this.ut.RequestMsg(msgJson)
		var loginRepMsg common.LoginResMsg
		this.ut.ResponseMsg(&loginRepMsg)
		//返回 登录信息 包括登录是否成功和错误
		return loginRepMsg

}

func (this *UserProcess)Registration(userid int,username string,userpwd string)common.RegistResMsg{
	var registinfo common.Msssage
	var registMsg common.RegistrationMsg

	registMsg.UserId = userid
	registMsg.Username = username
	registMsg.Password = userpwd

	data,err := json.Marshal(registMsg)
	if err != nil{
		fmt.Println("注册序列化信息错误")
	}
	registinfo.Data = data
	registinfo.MsgType = common.REGISTRATIONMSG

	msgJson,err := json.Marshal(registinfo)
	if err != nil{
		fmt.Println("注册Msg 序列化出错")
	}
	this.ut.RequestMsg(msgJson)
	var registRsgMsg common.RegistResMsg
	this.ut.ResponseMsg(&registRsgMsg)
	return registRsgMsg
}

func (this *UserProcess)Close(){
	this.conn.Close()
}

func (this *UserProcess)ReturnClint()net.Conn{
	return this.conn
}