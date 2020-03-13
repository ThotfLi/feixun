package process

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"feixun/common"
	"feixun/tools"
	"feixunservice/models"
	"feixunservice/redis"
	"fmt"
	"net"
	"sync"
)
//用户处理对象
var UserDao models.UserDao
var once sync.Once

//全局在线用户列表
var GlobalOnlineUserList UserInfoLister

//总控器 ，注册、登录、等功能
type Handles struct{
	Conn net.Conn
}
func (this Handles)MainHandle(message common.Msssage){
	//初始化信息处理对象
	RsQp := tools.Utils{Conn:this.Conn}

	//初始化用户处理对象
	once.Do(func() {
			UserDao = models.UserDao{Pool:redis.RedisPool}
			GlobalOnlineUserList = NewUserInfoList()
		})
	switch message.MsgType {
		//登录验证
		case common.LOGINMSG :
			//反序列化登录信息对象 拿到账号密码
			var logmsge common.LogMsg
			json.Unmarshal(message.Data,&logmsge)
			//初始化登录响应对象
			var rseMsg common.LoginResMsg
			//用户信息正确返回json后的User对象给客户端
			//验证登录请求
			user,err := UserDao.VerificationUser(logmsge)
			if err == nil{
				byteUser,_ := json.Marshal(*user)
				rseMsg.Userinfo = string(byteUser)
				//设置要返回的用户在线列表 其中不包括自己

				rseMsg.OnlineUserList = GlobalOnlineUserList.GetJsonUserList()

				//通知其他用户有新用户上线
				GlobalOnlineUserList.SendUserUpSignal(user)

				//添加新用户信息到在线用户列表
				GlobalOnlineUserList.Add(user,this.Conn)
				rseMsg.Code = 200
				fmt.Printf("用户：%s,已登录FX系统\n",(*user).UserName)
			}else if err == models.ERROR_USER_EXISTS || err== models.ERROR_USER_NOTEXISTS{
				rseMsg.Code = 403
				rseMsg.Err = errors.New("用户名或密码错误")
			}else {
				rseMsg.Code = 500
				rseMsg.Err = err
			}

			//发送登录成功msg
			rspDataByte,_ := json.Marshal(rseMsg)
			RsQp.RequestMsg(rspDataByte)
		//注册验证
		case common.REGISTRATIONMSG:
			//反序列化拿到msg
			var registMsg common.RegistrationMsg
			json.Unmarshal(message.Data,&registMsg)
			var rsgMsg common.RegistResMsg
			err := UserDao.Registration(registMsg)
			if err == nil{
				rsgMsg.Code = 200
			}else {
				rsgMsg.Code = 403
				rsgMsg.Err = models.ERROR_USER_EXISTS
			}
			rspJson,_ := json.Marshal(rsgMsg)
			RsQp.RequestMsg(rspJson)
		//X用户下线通知
		case common.NOTIFICATION_USERDOWN:
			var i uint32
			i = binary.BigEndian.Uint32(message.Data)
			println("用户下线了id为",int(i))
			GlobalOnlineUserList.SendUserDownSignal(int(i))
		//消息群发信号
		case common.GROUPS_MESSAGES:
			var groupMsg common.UserToUserMsg
			json.Unmarshal(message.Data,&groupMsg)

			var megRes common.MsgResMsg
			newT := tools.NewUtils(nil)
			//用户数量大于1代表有用户在线
			if len(GlobalOnlineUserList.ReturnMap()) > 1{
				var newMsg common.Msssage
				newMsg.Data = message.Data
				newMsg.MsgType = common.GROUPS_MESSAGES
				msgJson,err := json.Marshal(newMsg)
				if err != nil {
					println("json失败")
				}
				megRes.Code = 200
				for i,v := range GlobalOnlineUserList.ReturnMap(){
					if i == groupMsg.Sender {
						continue
					}
					newT.Conn = v.Conn
					newT.RequestMsg(msgJson)
				}
			}else {
				megRes.Code = 403
				megRes.Err = "没有用户在线"
			}
			rspJson,err := json.Marshal(megRes)
			if err != nil{
				println("json出错")
			}
			var newMsg common.Msssage
			newMsg.Data = rspJson
			newMsg.MsgType = common.GROUPS_MSG_RESPONSE
			msgJson,err := json.Marshal(newMsg)
			if err != nil{
				println("json失败")
			}
			RsQp.RequestMsg(msgJson)
	}

}
