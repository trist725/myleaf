package protobuf

import (
	"encoding/binary"
	"errors"

	"github.com/trist725/myleaf/chanrpc"
)

// -------------------------
// | clientID | id | protobuf message |
// -------------------------
type ServerProcessor struct {
	littleEndian bool
	//默认做转发的router
	defaultRouter *chanrpc.Server
	//暂存客户端ID
	clientID int32
}

func NewServerProcessor() *ServerProcessor {
	p := new(ServerProcessor)
	p.littleEndian = true
	return p
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *ServerProcessor) SetByteOrder(littleEndian bool) {
	p.littleEndian = littleEndian
}

func (p *ServerProcessor) SetDefaultRouter(msgRouter *chanrpc.Server) {
	p.defaultRouter = msgRouter
}

// goroutine safe
func (p *ServerProcessor) Route(msg interface{}, userData interface{}) error {
	msgByte := msg.([]byte)
	if p.defaultRouter != nil {
		p.defaultRouter.Go("ServerForward", msgByte[:4], userData, p.clientID)
	}

	return nil
}

// goroutine safe
func (p *ServerProcessor) Unmarshal(data []byte) (interface{}, error) {
	if len(data) < 2 {
		return nil, errors.New("protobuf data too short")
	}
	if p.littleEndian {
		p.clientID = int32(binary.LittleEndian.Uint32(data[:4]))
	} else {
		p.clientID = int32(binary.BigEndian.Uint32(data[:4]))
	}
	return data, nil
}

// goroutine safe
func (p *ServerProcessor) Marshal(msg interface{}) ([][]byte, error) {
	msgByte := msg.([][]byte)
	if len(msgByte) < 2 {
		return nil, errors.New("protobuf data too short")
	}
	return msgByte, nil
}
