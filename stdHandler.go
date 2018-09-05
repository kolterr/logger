package logger

import "io"

type StdHandler struct {
	w io.Writer
}

func NewStdHandler(w io.Writer) *StdHandler {
	return &StdHandler{w: w}
}

func (s *StdHandler) Write(d []byte) (error) {
	_, err := s.w.Write(d)
	return err
}

func (s *StdHandler) Close() error {
	return nil
}
