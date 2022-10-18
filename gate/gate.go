package gate

import (
	"time"

	"github.com/trist725/myleaf/chanrpc"
	"github.com/trist725/myleaf/network"
)

type Gate struct {
	MaxConnNum      int
	PendingWriteNum int
	MaxMsgLen       uint32
	Processor       network.Processor
	AgentChanRPC    *chanrpc.Server

	// websocket
	wsServer    *network.WSServer
	WSAddr      string
	HTTPTimeout time.Duration
	CertFile    string
	KeyFile     string

	// tcp
	tcpServer    *network.TCPServer
	TCPAddr      string
	LenMsgLen    int
	LittleEndian bool

	Mode int
}

func (gate *Gate) Run(closeSig chan bool) {
	if gate.WSAddr != "" {
		gate.wsServer = new(network.WSServer)
		gate.wsServer.Addr = gate.WSAddr
		gate.wsServer.MaxConnNum = gate.MaxConnNum
		gate.wsServer.PendingWriteNum = gate.PendingWriteNum
		gate.wsServer.MaxMsgLen = gate.MaxMsgLen
		gate.wsServer.HTTPTimeout = gate.HTTPTimeout
		gate.wsServer.CertFile = gate.CertFile
		gate.wsServer.KeyFile = gate.KeyFile
		gate.wsServer.NewAgent = func(conn *network.WSConn) network.Agent {
			a := &agent{conn: conn, gate: gate}
			if gate.AgentChanRPC != nil {
				switch gate.Mode {
				case 1:
					gate.AgentChanRPC.Go("NewClient", a)
				case 2:
					gate.AgentChanRPC.Go("NewServer", a)
				default:
					gate.AgentChanRPC.Go("NewAgent", a)
				}
			}
			return a
		}
	}

	if gate.TCPAddr != "" {
		gate.tcpServer = new(network.TCPServer)
		gate.tcpServer.Addr = gate.TCPAddr
		gate.tcpServer.MaxConnNum = gate.MaxConnNum
		gate.tcpServer.PendingWriteNum = gate.PendingWriteNum
		gate.tcpServer.LenMsgLen = gate.LenMsgLen
		gate.tcpServer.MaxMsgLen = gate.MaxMsgLen
		gate.tcpServer.LittleEndian = gate.LittleEndian
		gate.tcpServer.NewAgent = func(conn *network.TCPConn) network.Agent {
			a := &agent{conn: conn, gate: gate}
			if gate.AgentChanRPC != nil {
				switch gate.Mode {
				case 1:
					gate.AgentChanRPC.Go("NewClient", a)
				case 2:
					gate.AgentChanRPC.Go("NewServer", a)
				default:
					gate.AgentChanRPC.Go("NewAgent", a)
				}
			}
			return a
		}
	}

	if gate.wsServer != nil {
		gate.wsServer.Start()
	}
	if gate.tcpServer != nil {
		gate.tcpServer.Start()
	}
	<-closeSig
	if gate.wsServer != nil {
		gate.wsServer.Close()
	}
	if gate.tcpServer != nil {
		gate.tcpServer.Close()
	}
}

func (gate *Gate) OnDestroy() {}

func (gate *Gate) TCPConnsCount() int {
	return gate.tcpServer.ConnCount()
}

func (gate *Gate) WSConnsCount() int {
	return gate.wsServer.ConnCount()
}
