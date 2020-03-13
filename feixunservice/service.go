package main

import (
	"feixun/common"
	"feixunservice/process"
	"fmt"
	"io"
	"net"
	"feixun/tools"
	"feixunservice/redis"
)

//创建一个客户连接
func clientHandle(conn net.Conn){
	tool := tools.Utils{Conn:conn}

	mainHandles := process.Handles{conn}
	var message common.Msssage
	for {
		err := tool.ResponseMsg(&message)
		if err == nil{
			mainHandles.MainHandle(message)
		}else if err == io.EOF{
			fmt.Printf("链接已断开%s",conn.RemoteAddr().String())
			break
		}else {
			println(err)
			break
		}
	}

}


func main(){
	ln,err := net.Listen("tcp","127.0.0.1:8888")

	 //初始化redisPool
	redis.RedisPool = redis.NewRedisPool()
	defer redis.RedisPool.Close()
	if err != nil{
		println(err.Error())
		panic("服务器开始失败")
	}
	println("服务器运行中……")

	for {
		conn,err := ln.Accept()
		if err == nil{
			go clientHandle(conn)
		}


	}
	defer ln.Close()
}
