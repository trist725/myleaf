package gate

import (
	"github.com/trist725/myleaf/log"
	"github.com/trist725/myleaf/network"
	"net"
	"reflect"
)

type Agent interface {
	WriteMsg(msg interface{})
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	IsClosed() bool
	Destroy()
	UserData() interface{}
	SetUserData(data interface{})
}

type agent struct {
	conn     network.Conn
	gate     *Gate
	userData interface{}
}

func (a *agent) Run() {
	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			log.Debug("read message: %v", err)
			break
		}

		if a.gate.Processor != nil {
			msg, err := a.gate.Processor.Unmarshal(data)
			if err != nil {
				log.Debug("unmarshal message error: %v", err)
				break
			}
			err = a.gate.Processor.Route(msg, a)
			if err != nil {
				log.Debug("route message error: %v", err)
				break
			}
		}
	}
}

func (a *agent) OnClose() {
	if a.gate.AgentChanRPC != nil {
		var err error
		switch a.gate.Mode {
		case 1:
			err = a.gate.AgentChanRPC.Call0("CloseClient", a)
		case 2:
			err = a.gate.AgentChanRPC.Call0("CloseServer", a)
		default:
			err = a.gate.AgentChanRPC.Call0("CloseAgent", a)
		}
		if err != nil {
			log.Error("chanrpc error: %v", err)
		}
	}
}

func (a *agent) WriteMsg(msg interface{}) {
	if a.gate.Processor != nil {
		data, err := a.gate.Processor.Marshal(msg)
		if err != nil {
			log.Error("marshal message %v error: %v", reflect.TypeOf(msg), err)
			return
		}

		log.DebugMsg(reflect.TypeOf(msg), "send %T", msg)

		err = a.conn.WriteMsg(data...)
		if err != nil {
			log.Error("write message %v error: %v", reflect.TypeOf(msg), err)
		}
	}
}

func (a *agent) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *agent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

// Close 不能保证一定释放连接，但是它能保证一定尽可能的发送调用 Close 之前发送的消息（避免消息丢失）
func (a *agent) Close() {
	a.conn.Close()
}

// Destroy 方法不能保证调用 Destroy 之前发送的消息一定被发送，但是它能保证连接一定被释放
func (a *agent) Destroy() {
	a.conn.Destroy()
}

func (a *agent) IsClosed() bool {
	return a.conn.IsClosed()
}

func (a *agent) UserData() interface{} {
	return a.userData
}

func (a *agent) SetUserData(data interface{}) {
	a.userData = data
}
