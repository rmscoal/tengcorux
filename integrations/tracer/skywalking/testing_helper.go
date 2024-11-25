package skywalking

import "testing"

// recoverPanic recover from panics and set the testing as fail
func recoverPanic(t *testing.T) func() {
	return func() {
		if r := recover(); r != nil {
			t.Error("should not panic")
			t.FailNow()
		}
	}
}
