package tools

import (
	"fmt"
	"strings"
)

const (
	networkSplit = "@"
)

func ParseNetwork(str string) (network, addr string, err error) {
	//返回子串networkSplit在字符串str中第一次出现的位置
	if idx := strings.Index(str, networkSplit); idx == -1 {
		err = fmt.Errorf("addr: \"%s\" error, must be network@tcp:port or network@unixsocket", str)
		return
	} else {
		network = str[:idx] //连接方式(tcp)
		addr = str[idx+1:]  //连接地址
		return
	}
}
