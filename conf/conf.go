package conf

import (
	"reflect"
)

var (
	LenStackBuf = 4096

	// log
	LogLevel  string
	LogPath   string
	LogFlag   int
	MsgLog    = false // 收发消息日志开关
	BlackList = make(map[reflect.Type]struct{})

	// console
	ConsolePort   int
	ConsolePrompt string = "Leaf# "
	ProfilePath   string

	// cluster
	ListenAddr      string
	ConnAddrs       []string
	PendingWriteNum int
)

func OpenMsgLog(o bool) {
	MsgLog = o
}

func AddBlackList(t reflect.Type) {
	BlackList[t] = struct{}{}
}

func InBlackList(t reflect.Type) bool {
	if _, ok := BlackList[t]; ok {
		return true
	}
	return false
}
