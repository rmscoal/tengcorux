package tracetest

import "testing"

func TestNewRandomInteger(t *testing.T) {
	firstId := newRandomIntegerID()
	if firstId == 0 {
		t.Fatal("id should not be zero")
	}

	secondId := newRandomIntegerID()
	if secondId == firstId {
		t.Fatal("the second id should not be equal to the first id")
	}
}
