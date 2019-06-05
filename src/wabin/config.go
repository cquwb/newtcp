package wabin

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Addr         string `xml:"addr"`
	ReadTimeOut  uint32
	WriteTimeOut uint32
	MaxLen       uint32
}

func LoadConifg(addr string) (*Config, error) {
	f, err := os.Open(addr)
	if err != nil {
		fmt.Printf("read file error %s %v \n", addr, err)
		return nil, err
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Printf("read data error %v \n", err)
		return nil, err
	}
	c := Config{}
	err = xml.Unmarshal(data, &c)
	if err != nil {
		fmt.Printf("unmarshal data error %v \n", err)
		return nil, err
	}
	fmt.Printf("config is %+v \n", c)
	return &c, nil
}
