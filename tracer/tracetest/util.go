package tracetest

import (
	"math/rand"
	"time"
)

var globalRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func newRandomIntegerID() uint64 {
	return uint64(globalRand.Uint32())<<32 + uint64(globalRand.Uint32())
}
