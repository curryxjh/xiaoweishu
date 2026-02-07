package logger

import "sync"

var gl LoggerV1
var LMutex sync.RWMutex

func SetGlobalLogger(l LoggerV1) {
	LMutex.Lock()
	defer LMutex.Unlock()
	gl = l
}

func L() LoggerV1 {
	LMutex.RLock()
	defer LMutex.RUnlock()
	return gl
}

var GL LoggerV1 = &NopLogger{}
