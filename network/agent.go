package network

type Agent interface {
	Run()
	OnClose()
	IsClosed() bool
}
