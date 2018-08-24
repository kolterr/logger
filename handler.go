package logger

type Handler interface {
	Write([]byte) (error)
	Close() error
}
