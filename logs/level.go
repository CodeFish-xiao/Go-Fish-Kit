package logs

import (
	"fmt"
)

type Level int8

const (
	// TraceLevel level. 指定比调试更细粒度的信息事件。
	TraceLevel Level = iota - 2
	// DebugLevel level.通常只在调试时启用。非常详细的日志记录。
	DebugLevel
	// InfoLevel InfoLevel是默认的日志优先级。
	// 关于应用程序内部运行情况的一般操作条目。
	InfoLevel
	// WarnLevel level. 值得关注的非关键条目。
	WarnLevel
	// ErrorLevel level. Logs. 日志用于应该明确指出的错误。
	ErrorLevel
	// FatalLevel level.日志，然后调用`logger.Exit(1)`。严重程度最高。
	FatalLevel
)

func (l Level) String() string {
	switch l {
	case TraceLevel:
		return "TRACE"
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return ""
	}

}

// Enabled 如果给定级别位于或高于此级别，则启用返回true。
func (l Level) Enabled(lvl Level) bool {
	return lvl >= l
}

// GetLevel GetLevel将一个级别字符串转换为一个日志记录器的级别值。
// 如果输入字符串与已知值不匹配，则返回错误。
func GetLevel(levelStr string) (Level, error) {
	switch levelStr {
	case TraceLevel.String():
		return TraceLevel, nil
	case DebugLevel.String():
		return DebugLevel, nil
	case InfoLevel.String():
		return InfoLevel, nil
	case WarnLevel.String():
		return WarnLevel, nil
	case ErrorLevel.String():
		return ErrorLevel, nil
	case FatalLevel.String():
		return FatalLevel, nil
	}
	return InfoLevel, fmt.Errorf("Unknown Level String: '%s', defaulting to InfoLevel", levelStr)
}
