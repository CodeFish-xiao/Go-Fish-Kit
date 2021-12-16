package log

import "testing"

func TestDebug(t *testing.T) {
	Debug("cao")
}
func TestError(t *testing.T) {
	Error("cao")
}
func TestWarning(t *testing.T) {
	Warning("cao")
}
func TestWarn(t *testing.T) {
	Warn("cao")
}
func TestNotice(t *testing.T) {
	Notice("cao")
}
func TestInfo(t *testing.T) {
	Info("cao")
}
func TestInformational(t *testing.T) {
	Informational("cao")
}
