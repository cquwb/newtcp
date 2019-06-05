package main

import (
	"fmt"
	"time"
	"wabin"
)

type g_client_session struct {
	conn *wabin.Conn
}

func (this *g_client_session) Init(conn *wabin.Conn) {
	this.conn = conn
}

func (this *g_client_session) Close() {
	this.conn.Close()
}

func main() {
	for i := 0; i < 1000; i++ {
		go testServer(i)
	}
	time.Sleep(time.Second * 5)

}

func testServer(i int) {
	fmt.Printf("test %d \n", i)
	g_client_session := &g_client_session{}
	ok := wabin.TCPClient(g_client_session, "127.0.0.1:8080")
	if !ok {
		fmt.Printf("[client] begin error \n")
		return
	}
	ph := &wabin.PackHead{
		Len: 0,
		Cmd: uint32(1),
		Uid: 0,
		Sid: 0,
	}
	//msg := &LoginMessage{"hello world"}
	g_client_session.conn.Write(ph, []byte("hello world"))

}
