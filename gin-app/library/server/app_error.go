package server

type AppErrorI interface {
	Code() int
	Message() string
}

type AppError struct {
	code    int
	message string
}

func NewAppError(code int, message string) AppError {
	return AppError{code, message}
}

func (e AppError) Code() int {
	return e.code
}

func (e AppError) Message() string {
	return e.message
}

func (e AppError) Error() string {
	return e.message
}

var (
	AppErrorNone = NewAppError(0, "ok")
)
