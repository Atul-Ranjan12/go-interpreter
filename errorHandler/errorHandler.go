package errorHandler

type ErrorHandler interface {
	Error(line int, message string)
}
