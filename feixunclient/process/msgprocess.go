package process

import (
	"encoding/json"
	"feixun/common"
	"feixun/tools"
	"fmt"
)

//点对点、群发消息处理

type msgprocess struct{
	id   int
	name string
	ut   *tools.Utils
}

//从用户获取消息类型和信息内容
func(this msgprocess)AccessToInformation(){
	var a int
	var content string
	println("1、给其他用户发送消息")
	println("2、消息群发")

	fmt.Scanln(&a)
	switch a {
		case 1:
		case 2:
			println("请输入消息内容")
			fmt.Scanln(&content)
			this.sendMessagesInGroups(content)
	}
}

//群发消息
func(this msgprocess)sendMessagesInGroups(content string){
	var newMsg common.Msssage
	newMsg.MsgType = common.GROUPS_MESSAGES

	var newGroupsMsg common.UserToUserMsg
	newGroupsMsg.Data = content
	newGroupsMsg.Name = this.name
	newGroupsMsg.Sender = this.id

	groupsJson,err := json.Marshal(newGroupsMsg)
	if err != nil{
		println("json失败")
	}
	newMsg.Data = groupsJson
	msgJson,err := json.Marshal(newMsg)
	if err != nil {
		println("json失败")
	}

	this.ut.RequestMsg(msgJson)
}