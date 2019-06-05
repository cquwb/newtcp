package wabin

import (
	"fmt"
	"net"
)

func TCPServer(server Server, config *Config) {
	//address := "127.0.0.1:8080"
	//netType := "tcp"
	fmt.Printf("config is %s \n", config.Addr)
	l, e := net.Listen("tcp", config.Addr)
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
		c := NewConn(rw)
		s := server.NewSession()
		c.Server(s)
		//handlerConn(c) //1111构建面向对象的方式来处理
	}
}
