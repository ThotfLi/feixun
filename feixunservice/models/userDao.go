package models

import (
	"encoding/json"
	"feixun/common"
	"feixun/tools"
	"github.com/gomodule/redigo/redis"
)
//定义UserDao 结构体
//和redis连接池交互
//向redis中添加用户信息

type UserDao struct{
	Pool *redis.Pool
}
//根据用户名 返回用户User + err
func(this *UserDao) getUserByName(conn redis.Conn,name string)(*User,error){
	//通过id去redis中查询用户
	userString,err := redis.String(conn.Do("HGET","FXUSERS",name))
	if err != nil{
		//实际是未找到用户
		return nil,err
	}

	//反序列化从数据库拿到的用户信息
	var userObj User
	tools.JsonToStruct([]byte(userString),&userObj)
	return &userObj,nil
}

//输入客户端发来的登录请求用户对象和数据库用户对象做比对
func(this *UserDao)VerificationUser(user common.LogMsg)(*User,error){
	dataBaseUser,err := this.getUserByName(this.Pool.Get(),user.Username)
	//没找到用户
	if err!= nil{
		return nil,err
	}
	if dataBaseUser == nil{
		return nil, ERROR_USER_NOTEXISTS
	}

	//验证密码是否正确
	if (*dataBaseUser).UserPwd == user.Password{
		return dataBaseUser,nil
	}else {
		//密码不正确
		return nil, ERROR_USER_EXISTS
	}
}

//注册
//返回 err == nil 则注册成功
func(this *UserDao)Registration(info common.RegistrationMsg)error{
	user,_ := this.getUserByName(this.Pool.Get(),info.Username)
	//err ！= nil 证明不存在此用户可以注册
	if user == nil {
		//注册用户
		err := this.add(info)
		//这个err == nil为注册成功
		return err

	}else{
		//err == nil
		return ERROR_USER_EXISTS
	}
}

func(this *UserDao)add(info common.RegistrationMsg)error{
	var newUser User
	newUser.UserName = info.Username
	newUser.UserPwd = info.Password
	newUser.UserId = info.UserId

	strJson,_ := json.Marshal(newUser)
	conn := this.Pool.Get()
	_ ,err := conn.Do("HSET","FXUSERS",newUser.UserName,string(strJson))
	return err
}