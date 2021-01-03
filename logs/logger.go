package logs

// Logger is a logger interface.
type Logger interface {
	Print(kvpair ...interface{})
}
