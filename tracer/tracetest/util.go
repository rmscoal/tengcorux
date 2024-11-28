package tracetest

import "math/rand"

// newRandomIntegerID generate a random integer ID. Please note, that it might
// not be the most cryptic random integers generated.
func newRandomIntegerID() uint64 {
	return uint64(rand.Uint32())<<32 + uint64(rand.Uint32())
}
