package client
import(
	"net"
)

func NewClient()net.Conn{
	conn,err := net.Dial("tcp","127.0.0.1:8888")
	if err != nil{
		panic("链接服务器失败")
	}
	return conn
}
