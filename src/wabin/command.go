package wabin

import (
	"fmt"
)

type Service func(Sessioner, *PackHead, []byte) bool

type Command struct {
	Services map[uint32]Service
}

func NewCommand() *Command {
	return &Command{
		Services: make(map[uint32]Service),
	}
}
func (this *Command) Register(id uint32, s Service) {
	if _, exists := this.Services[id]; !exists {
		this.Services[id] = s
	}
}

func (this *Command) Dispatch(sess Sessioner, ph *PackHead, data []byte) bool {
	if s, exists := this.Services[ph.Cmd]; exists {
		return s(sess, ph, data)
	} else {
		fmt.Printf("command id %d not exitss \n", ph.Cmd)
		return false
	}
}
