package tracing

import "testing"

func TestVersion(t *testing.T) {
	if Version() != "v0.1.0" {
		t.Fatal("expected version to be v0.1.0")
	}
}
