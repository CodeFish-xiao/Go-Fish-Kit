package logs

type nopLogger struct{}

func (n *nopLogger) Print(kvpair ...interface{}) {}
