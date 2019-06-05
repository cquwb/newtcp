package main

import (
	"fmt"
	"protocol/cs"

	//"time"
	"wabin"

	"github.com/golang/protobuf/proto"
)

type g_client_session struct {
	conn *wabin.Conn
}

func (this *g_client_session) Init(conn *wabin.Conn) {
	this.conn = conn
	this.conn.AddWaitGroup()
	go this.Run()
}

func (this *g_client_session) Close() {
	this.conn.DecWaitGroup()
	this.conn = nil
	//this.conn.Close()
}

func (this *g_client_session) Run() {
	for {
		select {
		case msg := <-this.conn.ReadChann:
			g_command.Dispatch(this, wabin.ByteToHead(msg[:20]), msg[20:])
		case <-this.conn.CloseChan:
			return
		}
	}
}

type LoginMessage struct {
	str string
}

var g_command = wabin.NewCommand()

func main() {
	Register()
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
		Cmd: uint32(cs.ID_ID_C2S_Hello),
		Uid: 0,
		Sid: 0,
	}
	msg := &cs.C2S_Hello{}
	msg.Id = proto.Uint32(23)
	msg.Msg = proto.String("hello world")
	g_client_session.conn.Write(ph, msg)
	//time.Sleep(time.Second * 5)
	//g_client_session.Close()
	<-c
}

func Register() {
	g_command.Register(uint32(cs.ID_ID_S2C_Hello), S2C_Hello)
}

func S2C_Hello(s wabin.Sessioner, ph *wabin.PackHead, data []byte) bool {
	sess, _ := s.(*g_client_session)
	receiveMsg := &cs.S2C_Hello{}
	if err := proto.Unmarshal(data, receiveMsg); err != nil {
		fmt.Printf("S2C_Hello proto unmarshal error %v \n", err)
		return false
	}
	fmt.Printf("respone msg is ret %d %s %v \n", receiveMsg.GetRet(), receiveMsg.GetMsg(), sess)
	return true
}
