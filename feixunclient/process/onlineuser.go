package process

import (
	"encoding/binary"
	"encoding/json"
	"feixunservice/process"
	"fmt"
)

type onlineUsers map[int]*process.OnlineUser

var OnlineUsers onlineUsers

func(p *onlineUsers)UpSignalProcess(data []byte){
	var newUserInfo process.OnlineUser
	json.Unmarshal(data,&newUserInfo)

	OnlineUsers[newUserInfo.UserId] = &newUserInfo
}

func(p *onlineUsers)DownSignalProcess(data []byte){
	var i uint32
	i = binary.BigEndian.Uint32(data)
	delete(OnlineUsers,int(i))
}

func(p onlineUsers)PrintUserInfo(){
	for i,v := range p{
		fmt.Println(i,v)
	}
}
