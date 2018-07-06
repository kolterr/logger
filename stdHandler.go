package logger

import "io"

type StdHandler struct {
	w io.Writer
}

func NewStdHandler(w io.Writer) *StdHandler {
	return &StdHandler{w: w}
}

func (s *StdHandler) Write(d []byte) (int, error) {
	return s.w.Write(d)
}

func (s *StdHandler) Close() error {
	return nil
}
