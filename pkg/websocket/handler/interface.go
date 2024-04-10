package handler

type Interface interface {
	RegisterId(id int64)
	Id() int64
	RegisterRemoveChan(ch chan<- int64)
	RegisterConnWriteChan(ch chan<- []byte)
	RegisterConnClose(do func())
	RegisterConnPing(do func())
	Ping()
	Handler(data []byte) (res []byte, err error)
	Run()
	Close()
}

type MessageHandler interface {
	Read() <-chan []byte
	Write(in []byte) error
}
