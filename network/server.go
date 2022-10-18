package network

type IServer interface {
	Start()
	Close()
	ConnCount() int
}
