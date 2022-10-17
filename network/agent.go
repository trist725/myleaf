package network

type Agent interface {
	Run()
	OnClose()
	IsClose() bool
}
