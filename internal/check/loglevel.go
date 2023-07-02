// Package check contains value domain checkers.
package check

// Loglevel returns true if the given string level is an acceptable loglevel name
// subset of logrus levels are supported as it is expected that panics as exceptions and fatal
// should really just be the final exit reason. This distinction is not meaningful to
// a user. Errors are faults that the runtime survives, Fatal is the final exit.
func Loglevel(level string) bool {
	switch level {
	case "debug":
	case "info":
	case "warn", "warning":
	case "error":
	case "fatal":

	default:
		return false
	}
	return true
}
