// Package example contains examples - the entire package can be deleted.
package example

import "strings"

// Greet returns "hello world" with space and name appended if given.
func Greet(s string) string {
	out := "hello world"
	if s != "" {
		return out + " " + s
	}
	return out
}

// GreetWonderful returns "hello wonderful world" with space and name appended if given.
func GreetWonderful(s string) string {
	out := "hello wonderful world"
	if s != "" {
		return out + " " + s
	}
	return out
}

// GreetUpper returns upper cased version of [Greet].
func GreetUpper(s string) string {
	return strings.ToUpper(Greet(s))
}

// GreetWonderfulUpper returns upper cased version of [GreetWonderful].
func GreetWonderfulUpper(s string) string {
	return strings.ToUpper(GreetWonderful(s))
}
