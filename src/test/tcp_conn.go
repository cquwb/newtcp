package wabin

import (
	"fmt"
	"net"
)

func TCPServer(address string) {
	//address := "127.0.0.1:8080"
	//netType := "tcp"
	l, e := net.Listen("tcp", address)
	if e != nil {
		fmt.Printf("connect net error %v \n", e)
		panic(e.Error()) //111 panic
	}
	defer l.Close() //111111 defer close
	for {
		rw, e := l.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				continue
			}
			fmt.Printf("[TCPServer] accept error %v", e)
			return
		}
		handlerConn(rw) //1111构建面向对象的方式来处理
	}
}

func handlerConn(rwc net.Conn) {
	conn := NewConn(rwc)
	go conn.readLoop()
	go conn.readLoop()
}



type conn struct {
	rwc net.Conn
}

func NewConn(rwc net.Conn) *conn {
	return &conn {
		rwc:rwc
	}
}

func (this *conn) readLoop() {
	l := make([]byte, 4, 4)
	for {
		n, err := this.rwc.Read(l)
		if err != nil {
			fmt.Printf("[tcp_conn] readLoop error %v \n", err)
			break
		}
		if n != len(l) {
			fmt.Printf("[tcp_conn] readLoop len error %d", n)
			break
		}
		
		
	}
	
}

func byteToUint32(b []byte) uint32 {
	
}

func (this *conn) writeLoop() {

}

func TCPClient(address string)  {
	rw, e := net.Dial("tcp", address)
	if e != nil {
		fmt.Printf("[TCPClient] error %v \n", e)
		panic(e.Error())
	}
	conn := NewConn(rwc)
	rw.Write([]byte(30))
	rw.Write([]byte("Hello World"))
}