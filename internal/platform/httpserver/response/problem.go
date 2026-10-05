package response

import "net/http"

const problemContentType = "application/problem+json"

// Problem is an RFC 9457 problem details object.
type Problem struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail,omitempty"`
}

func NewProblem(statusCode int, detail string) Problem {
	return Problem{
		Type:   "about:blank",
		Title:  http.StatusText(statusCode),
		Status: statusCode,
		Detail: detail,
	}
}
