package wabin

import (
	"fmt"
	"net"
)

func TCPClient(s Sessioner, config *Config) bool {
	rw, e := net.Dial("tcp", config.Addr)
	if e != nil {
		fmt.Printf("[TCPClient] error %v \n", e)
		return false
	}
	conn := NewConn(rw)
	conn.Server(s)
	return true
}
