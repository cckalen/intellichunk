package example_test

import (
	"testing"

	"github.com/hlindberg/testutils"
)

func Test_examples_of_testutils(t *testing.T) {
	// Create a tester which supports the testing check functions as methods and thus avoids having to
	// pass *testing.T in every call. This makes the test source less cluttered.
	tt := testutils.NewTester(t)
	tt.CheckEqual("A", "A")

	// CheckNumericXXX handles all numeric data types (which is otherwise a pain)
	// Here, the tested value should be less then the expected (first arg) value
	tt.CheckNumericLess(10, 0)

	// Looping test over results
	values := []int{1, 2, 3}
	for i := 0; i < len(values); i++ {
		// At(i) makes the tester reports error/fail for that index
		// CheckEqual handles all types of values
		tt.At(i).CheckEqual(i, values[i]-1)
	}
}
