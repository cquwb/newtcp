package wabin

import (
	"fmt"
	"net"

	//"encoding/binary"
	"bufio"
	"errors"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
)

var (
	ErrorMsgType  = errors.New("tcp_conn: error msg type")
	WriteOverflow = errors.New("tcp_conn: wirte buffer over flow")
)

const HEADLEN = 24
const MAX_SIZE = 512
const READ_TIME_OUT = 2

const (
	StateInit = iota
	StateDisconnected
	StateConnected
)

type Conn struct {
	rwc        net.Conn
	sess       Sessioner
	ReadChann  chan []byte
	WriteChann chan *Message

	CloseChan chan struct{} //保证出错了能够退出程序

	wg sync.WaitGroup //保证所有goroutine都退出了在进行清理操作

	state int32 //保证只close CloseChan 一次
}

func NewConn(rwc net.Conn) *Conn {
	return &Conn{
		rwc: rwc,
	}
}

//启动链接
func (this *Conn) Server(s Sessioner) {
	if !atomic.CompareAndSwapInt32(&this.state, StateInit, StateConnected) {
		return
	}
	this.ReadChann = make(chan []byte, 1024)
	this.WriteChann = make(chan *Message, 1024)
	this.CloseChan = make(chan struct{})
	this.sess = s
	this.wg.Add(1)
	go this.readLoop()
	this.wg.Add(1)
	go this.writeLoop()
	s.Init(this)
}

func (this *Conn) Close() {
	this.rwc.Close()
}

func (this *Conn) Stop() {
	fmt.Printf("conn begin stop \n")
	if !atomic.CompareAndSwapInt32(&this.state, StateConnected, StateDisconnected) {
		return
	}
	fmt.Printf("stop sucess 1 \n")
	close(this.CloseChan)
	this.Close()
	go func() {
		this.wg.Wait()
		//todo 一些clear
		this.sess.Close()
		this.sess = nil
		fmt.Printf("conn stop sucess \n")
		atomic.StoreInt32(&this.state, StateInit)
	}()
}

func (this *Conn) readLoop() {
	i := 0
	rbuf := bufio.NewReader(this.rwc)
	for {
		i++
		this.rwc.SetReadDeadline(time.Now().Add(READ_TIME_OUT * time.Second))
		fmt.Printf("[tcp_conn] readLoop times %d \n", i)
		l := make([]byte, MESSAGE_LEN, MESSAGE_LEN)
		_, err := io.ReadFull(rbuf, l)
		if err != nil {
			fmt.Printf("[tcp_conn] readLoop error %v \n", err)
			goto exit
		}

		nl := byteToUint32(l)
		if nl < MESSAGE_LEN || nl > MAX_SIZE {
			fmt.Printf("[tcp_conn] readLoop new len error %d \n", nl)
			goto exit
		}
		fmt.Printf("[tcp_conn] readLoop new len %d \n", nl)
		nb := make([]byte, nl-MESSAGE_LEN, nl-MESSAGE_LEN)
		_, e := io.ReadFull(rbuf, nb)
		if e != nil {
			fmt.Printf("[tcp_conn] readLoop read data  error %v \n", e)
			goto exit
		}
		fmt.Printf("read data 1 %v \n", nb)
		this.ReadChann <- nb
	}
exit:
	this.Stop()
	this.wg.Done()
}

func (this *Conn) writeLoop() {
	write_data := make([]byte, MAX_SIZE)
	head_buff := make([]byte, HEADLEN)
	data_buff := make([]byte, MAX_SIZE-HEADLEN)
	for {
		index := 0
		select {
		case msg := <-this.WriteChann:
			length, data, err := Marshal(msg.PH, msg.Info, MAX_SIZE, head_buff, data_buff)
			if err == nil {
				copy(write_data, head_buff)
				copy(write_data[HEADLEN:], data)
				index += length
				for more := true; more; {
					select {
					case msg := <-this.WriteChann:
						length, data, err := Marshal(msg.PH, msg.Info, MAX_SIZE, head_buff, data_buff)
						if err == nil {
							if length+index > MAX_SIZE {
								this.rwc.Write(data[:index])
								copy(write_data, head_buff)
								copy(write_data[HEADLEN:], data)
								index = length
							} else {
								copy(write_data[index:], head_buff)
								copy(write_data[index+HEADLEN:], data)
								index += length
							}
						}
					case <-this.CloseChan:
						goto exit
					default:
						more = false
					}
				}
				fmt.Printf("write data is %v \n", write_data[:index])
				this.rwc.Write(write_data[:index])
			}

		case <-this.CloseChan:
			goto exit
		}
	}
exit:
	this.Stop()
	this.wg.Done()
}

func (this *Conn) Write(PH *PackHead, msg interface{}) bool {
	select {
	case this.WriteChann <- &Message{PH, msg}:
	case <-this.CloseChan:
		return false
	}
	return true
}

func (this *Conn) AddWaitGroup() {
	this.wg.Add(1)
}

func (this *Conn) DecWaitGroup() {
	this.wg.Done()
}

func Marshal(ph *PackHead, msg interface{}, max_size uint32, head_buff, data_buff []byte) (int, []byte, error) {
	var data []byte
	var err error
	switch v := msg.(type) {
	case []byte:
		data = v
	case proto.Message:
		data, err = proto.Marshal(v)
		if err != nil {
			fmt.Printf("[Marshal] error %v \n", err)
			return 0, nil, ErrorMsgType
		}
	default:
		fmt.Printf("[Marshal] error %v \n", ph)
		return 0, nil, ErrorMsgType
	}
	length := len(data) + HEADLEN
	if length > MAX_SIZE {
		fmt.Printf("[Marshal] over flow %v \n", ph)
		return 0, nil, WriteOverflow
	}
	ph.Len = uint32(length)
	HeadToByte(ph, head_buff)
	return length, data, nil

}
