package main

import (
	"encoding/json"
	"feixun/process"
	"feixunservice/models"
	"fmt"
	"os"
)
var username string
var userpwd  string
var userid   int
func main(){
	var loop bool = true
	var number int
	for loop{
		fmt.Println("------欢迎登录飞讯聊天系统------")
		fmt.Println("1.登录")
		fmt.Println("2.注册")
		fmt.Println("3.退出")
		fmt.Scanf("%d\n",&number)

		//登录和注册处理对象
		pup := process.NewUserProcess()

		defer pup.Close()

		switch number {
		case 1:
			println("请输入用户名")
			fmt.Scanln(&username)
			println("请输入密码")
			fmt.Scanln(&userpwd)

			logRsp := pup.Login(username,userpwd)

			var userStut models.User
			json.Unmarshal([]byte(logRsp.Userinfo),&userStut)

			if logRsp.Code == 200 {
				//初始化在线用户列表
				json.Unmarshal(logRsp.OnlineUserList,&process.OnlineUsers)

				ser := process.NewServer(userStut,pup.ReturnClint())
				ser.TwoLoginIn()
			}
		case 2:
			println("请输入用户id")
			fmt.Scanln(&userid)
			println("请输入用户名")
			fmt.Scanln(&username)
			println("请输入密码")
			fmt.Scanln(&userpwd)

			registRsp := pup.Registration(userid,username,userpwd)
			if registRsp.Code == 200{
				println("注册成功")
			}else {
				println("注册失败")
				println(registRsp.Code)
				fmt.Println(registRsp.Err)
			}
		case 3:
			fmt.Println("退出系统")
			os.Exit(0)
		}

	}

}