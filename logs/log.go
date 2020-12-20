package logs

// level is a log level
type Level int32

const (
	LevelFatal Level = iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
)
