package contract

type ILogger interface {
	Info(msg string)
	Warn(msg string)
	Debug(msg string)
	Error(msg string)
	Secure(tag string, msg string)
}
