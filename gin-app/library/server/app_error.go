package server

var (
	ErrorNone = NewError(0, "ok")
)

type ErrorI interface {
	Code() int
	Message() string
}

type Error struct {
	code    int
	message string
}

func NewError(code int, message string) Error {
	return Error{code, message}
}

func (e Error) Code() int {
	return e.code
}

func (e Error) Message() string {
	return e.message
}

func (e Error) Error() string {
	return e.message
}
