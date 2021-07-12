package contract

import (
	"io"
	"net/http"
)

type MessageInterface interface {
	GetBody() io.ReadCloser
	GetHeader() http.Header
}

type ResponseContract interface {
	MessageInterface
	GetStatusCode() int
	//WithStatus() object.HashMap
	//GetReasonPhrase() object.HashMap
}

type PromiseInterface interface {
}
