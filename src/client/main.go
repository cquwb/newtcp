package main

import (
	"fmt"
	"protocol/cs"
	"time"
	"wabin"

	"github.com/golang/protobuf/proto"
)

type g_client_session struct {
	conn *wabin.Conn
}

func (this *g_client_session) Init(conn *wabin.Conn) {
	this.conn = conn
}

func (this *g_client_session) Close() {
	this.conn = nil
	//this.conn.Close()
}

type LoginMessage struct {
	str string
}

func main() {
	var config *wabin.Config
	var err error
	if config, err = wabin.LoadConifg("../../config/client_config.xml"); err != nil {
		fmt.Printf("begin tcp server  error %v \n", err)
		return
	}
	c := make(chan uint32, 1)
	fmt.Printf("begin tcp client \n")
	g_client_session := &g_client_session{}
	ok := wabin.TCPClient(g_client_session, config)
	if !ok {
		fmt.Printf("[client] begin error \n")
		return
	}
	ph := &wabin.PackHead{
		Len: 0,
		Cmd: uint32(25),
		Uid: 0,
		Sid: 0,
	}
	msg := &cs.C2S_Hello{}
	msg.Id = proto.Uint32(2)
	msg.Msg = proto.String("hello world")
	g_client_session.conn.Write(ph, msg)
	time.Sleep(time.Second * 5)
	g_client_session.Close()
	<-c
}
