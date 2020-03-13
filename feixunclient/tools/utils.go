package tools

import (
	"encoding/binary"
	"encoding/json"
	"feixun/client"
	"net"
)

type Utils struct{
	buf []byte  //储存rspMsg
	Conn net.Conn
}
func NewUtils(conn net.Conn)Utils{
	return Utils{
		Conn:conn,
	}
}
func NewClientUtils()Utils{
	return Utils{
		Conn: client.NewClient(),
	}
}
func (this *Utils)RequestMsg(msgJson []byte){
	//第一次发送msg长度
	dataLen := uint32(len(msgJson))
	lenBuf := make([]byte,4)
	binary.BigEndian.PutUint32(lenBuf,dataLen)
	this.Conn.Write(lenBuf)
	//发送msg
	this.Conn.Write(msgJson)
}


// 传入msg 是要从服务器得到的消息结构体的零值
//返回要得到消息结构体赋值后的值

func (this *Utils)ResponseMsg(msg interface{})error{
	var datalen uint32
	lenBuf := make([]byte,4)
	n,err := this.Conn.Read(lenBuf)
	if err != nil{
		return err
	}
	datalen = binary.BigEndian.Uint32(lenBuf)
	this.buf = make([]byte,datalen)

	n,err = this.Conn.Read(this.buf)
	if err != nil{
		return err
	}
	if uint32(n) != datalen{
		return err
	}

	json.Unmarshal(this.buf,msg)
	return nil
}