package response

import "net/http"

// ResponseWriter records the status code and body size of a response.
type ResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
	wroteHeader  bool
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	// 1xx responses other than 101 are informational and precede the final status.
	isInformational := statusCode < http.StatusOK && statusCode != http.StatusSwitchingProtocols

	if !rw.wroteHeader && !isInformational {
		rw.statusCode = statusCode
		rw.wroteHeader = true
	}

	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	rw.wroteHeader = true

	n, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten += n

	return n, err
}

func (rw *ResponseWriter) StatusCode() int {
	return rw.statusCode
}

func (rw *ResponseWriter) BytesWritten() int {
	return rw.bytesWritten
}

func (rw *ResponseWriter) WroteHeader() bool {
	return rw.wroteHeader
}

func (rw *ResponseWriter) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}
