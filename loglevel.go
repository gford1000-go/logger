package logger

// LogLevel represents the level applied to the statement being logged
type LogLevel int

const (
	None LogLevel = iota
	Error
	Warn
	Info
	Debug
	All
)

func (l LogLevel) String() string {
	switch l {
	case None:
		return "NONE"
	case Error:
		return "ERROR"
	case Warn:
		return "WARN"
	case Info:
		return "INFO"
	case Debug:
		return "DEBUG"
	}
	return "Unknown"
}
