package main

import (
	"fmt"
	"net"
	"protocol/cs"
	"time"
	"wabin"

	"github.com/golang/protobuf/proto"

	"net/http"
	_ "net/http/pprof"
)

type g_server struct {
	s    map[uint32]*ReceiveSession
	msgs chan *ServerMsg
}

func (this *g_server) NewSession() wabin.Sessioner {
	return &ReceiveSession{}
}

type ServerMsg struct {
	s    *ReceiveSession
	data []byte
}

type ReceiveSession struct {
	Conn *wabin.Conn
}

func (this *ReceiveSession) Init(c *wabin.Conn) {
	this.Conn = c
	this.Conn.AddWaitGroup()
	go this.Run()
}

func (this *ReceiveSession) Close() {
	this.Conn = nil
}

func (this *ReceiveSession) Run() {
	defer this.Conn.DecWaitGroup()
	for {
		select {
		case msg := <-this.Conn.ReadChann:
			smsg := &ServerMsg{this, msg}
			server.msgs <- smsg
		case <-this.Conn.CloseChan:
			return

		}
	}

}

var server = &g_server{}

var g_command = wabin.NewCommand()

func main() {
	Register()
	var config *wabin.Config
	var err error
	if config, err = wabin.LoadConifg("../../config/server_config.xml"); err != nil {
		fmt.Printf("begin tcp server  error %v \n", err)
		return
	}
	go func() {
		http.ListenAndServe("127.0.0.1:6060", nil)
	}()
	fmt.Printf("begin tcp server \n")
	server.msgs = make(chan *ServerMsg, 1024)
	go hanlderMsg()
	wabin.TCPServer(server, config)
	/*
		l, e := net.Listen("tcp", "127.0.0.1:8080")
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
	*/
}

func hanlderMsg() {
	for {
		select {
		case msg := <-server.msgs:
			g_command.Dispatch(msg.s, wabin.ByteToHead(msg.data[:20]), msg.data[20:])

			/*
				fmt.Printf("server msgs %v \n", msg)
				data := &cs.C2S_Hello{}
				ph := wabin.ByteToHead(msg[:20])
				fmt.Printf("msg cmd is %d \n", ph.Cmd)
				if err := proto.Unmarshal(msg[20:], data); err == nil {
					fmt.Printf("proto msg is %d %s \n", data.GetId(), data.GetMsg())
				} else {
					fmt.Printf("proto msg error  %v \n", err)
				}*/
		}
	}
}

func handlerConn(rw net.Conn) {

	i := 0
	for {
		i++
		fmt.Printf("[tcp_conn] readLoop times %d \n", i)
		rw.SetReadDeadline(time.Now().Add(time.Second * 2))
		l := make([]byte, 5, 5)
		n, err := rw.Read(l)
		if err != nil {
			fmt.Printf("[tcp_conn] readLoop error %v \n", err)
			break
		}
		fmt.Printf("hanlderConn read data %d %v \n", n, l)
		//time.Sleep(time.Second * 1)
	}
}

func Register() {
	g_command.Register(uint32(cs.ID_ID_C2S_Hello), C2S_Hello)
}

func C2S_Hello(s wabin.Sessioner, ph *wabin.PackHead, data []byte) bool {
	sess, _ := s.(*ReceiveSession)
	receiveMsg := &cs.C2S_Hello{}
	if err := proto.Unmarshal(data, receiveMsg); err != nil {
		fmt.Printf("C2S_Hello proto unmarshal error %v \n", err)
		return false
	}
	fmt.Printf("get msg is %d %s %v \n", receiveMsg.GetId(), receiveMsg.GetMsg(), sess)
	msg := &cs.S2C_Hello{}
	msg.Ret = proto.Uint32(uint32(1))
	msg.Msg = proto.String("ok")
	newHead := &wabin.PackHead{
		Cmd: uint32(cs.ID_ID_S2C_Hello),
		Uid: ph.Uid,
		Sid: ph.Sid,
	}
	sess.Conn.Write(newHead, msg)
	return true
}
