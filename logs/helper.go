package logs

var nop Logger = new(nopLogger)

// Helper is a logger helper.
type Helper struct {
	Logger

	opts Options
}
