package response

import "net/http"

var (
	StatusCodeUninitialized = -1
)

type Writer struct {
	http.ResponseWriter
	statusCode int
}

func NewWriter(w http.ResponseWriter) *Writer {
	return &Writer{w, StatusCodeUninitialized}
}

func (w *Writer) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func (w *Writer) GetStatusCodeOrPanic() int {
	if w.statusCode == StatusCodeUninitialized {
		panic("status code not initialized")
	}

	return w.statusCode
}
