package connect

import (
	"net"

	"chat/proto"

	"github.com/gorilla/websocket"
)

type Channel struct {
	Room      *Room    //指向room的指针
	Next      *Channel //channel间互连形成双向链表
	Prev      *Channel
	broadcast chan *proto.Msg
	userId    int
	conn      *websocket.Conn //websocket的连接
	connTcp   *net.TCPConn    //tcp的连接
}

func NewChannel(size int) (c *Channel) {
	c = new(Channel)
	c.broadcast = make(chan *proto.Msg, size)
	c.Next = nil
	c.Prev = nil
	return
}

func (ch *Channel) Push(msg *proto.Msg) (err error) {
	select {
	case ch.broadcast <- msg:
	default:
	}
	return
}
