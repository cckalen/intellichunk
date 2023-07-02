package example_test

import (
	"testing"

	"github.com/hlindberg/testutils"
	"github.com/wyrth-io/goapp-template/internal/example"
)

func Test_Greet_returns_hello_without_name(t *testing.T) {
	testutils.CheckEqual("hello world", example.Greet(""), t)
}

func Test_Greet_returns_hello_with_name_when_given(t *testing.T) {
	testutils.CheckEqual("hello world Albert", example.Greet("Albert"), t)
}

func Test_GreetUpper_returns_HELLO_without_name(t *testing.T) {
	testutils.CheckEqual("HELLO WORLD", example.GreetUpper(""), t)
}

func Test_GreetUpper_returns_HELLO_with_name_upcased_when_given(t *testing.T) {
	testutils.CheckEqual("HELLO WORLD ALBERT", example.GreetUpper("Albert"), t)
}

func Test_GreetWonderful_returns_hello_world_without_name(t *testing.T) {
	testutils.CheckEqual("hello wonderful world", example.GreetWonderful(""), t)
}

func Test_GreetWonderful_returns_hello_with_name_when_given(t *testing.T) {
	testutils.CheckEqual("hello wonderful world Albert", example.GreetWonderful("Albert"), t)
}

func Test_GreetWonderfulUpper_returns_HELLO_WORLD_without_name(t *testing.T) {
	testutils.CheckEqual("HELLO WONDERFUL WORLD", example.GreetWonderfulUpper(""), t)
}

func Test_GreetWorldUpper_returns_hello_with_name_when_given(t *testing.T) {
	testutils.CheckEqual("HELLO WONDERFUL WORLD ALBERT", example.GreetWonderfulUpper("Albert"), t)
}
