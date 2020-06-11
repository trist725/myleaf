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
	var clientID int32
	if p.littleEndian {
		clientID = int32(binary.LittleEndian.Uint32(msgByte[:4]))
	} else {
		clientID = int32(binary.BigEndian.Uint32(msgByte[:4]))
	}

	if p.defaultRouter != nil {
		p.defaultRouter.Go("ServerForward", msgByte[:4], userData, clientID)
	}

	return nil
}

// goroutine safe
func (p *ServerProcessor) Unmarshal(data []byte) (interface{}, error) {
	if len(data) < 2 {
		return nil, errors.New("protobuf data too short")
	}
	return data, nil
}

// goroutine safe
func (p *ServerProcessor) Marshal(msg interface{}) ([][]byte, error) {
	msgByte := msg.([]byte)
	if len(msgByte) < 2 {
		return nil, errors.New("protobuf data too short")
	}
	return [][]byte{msgByte[:4], msgByte[4 : 4+2], msgByte[4+2:]}, nil
}
