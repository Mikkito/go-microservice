package utils

import "log"

type AsyncError struct {
	Source string
	Err    error
}

type ErrorHandler struct {
	ch chan AsyncError
}

func NewErrorHandler(buffer int) *ErrorHandler {
	return &ErrorHandler{
		ch: make(chan AsyncError, buffer),
	}
}

func (h *ErrorHandler) Start() {
	go func() {
		for e := range h.ch {
			log.Printf("[ERROR] source=%s err=%v\n", e.Source, e.Err)
		}
	}()
}

func (h *ErrorHandler) Report(e AsyncError) {
	select {
	case h.ch <- e:
	default:
	}
}
