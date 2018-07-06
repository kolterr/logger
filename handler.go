package logger

type Handler interface {
	Write([]byte) (int, error)
	Close() error
}
