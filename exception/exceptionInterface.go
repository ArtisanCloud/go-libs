package exception

type Throwable interface {
	GetMessage() string
	GetCode() int
	GetFile() string
	GetLine() int
	GetTrace() []string
	GetTraceAsString() string
	GetPrevious() *Throwable
	__toString() string
}
