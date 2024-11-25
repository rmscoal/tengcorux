package tracer

import (
	"context"
	"sync"
	"time"
)

var (
	mu           sync.RWMutex
	globalTracer Tracer = &NoopTracer{}
)

// SetGlobalTracer replaces the current tracer to the provided.
func SetGlobalTracer(tracer Tracer) {
	mu.Lock()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Nanosecond)
	defer cancel()
	err := globalTracer.Shutdown(ctx)
	if err != nil {
		return
	}
	defer mu.Unlock()
	globalTracer = tracer
}

// GetGlobalTracer retrieves the globalTracer.
func GetGlobalTracer() Tracer {
	mu.RLock()
	defer mu.RUnlock()
	return globalTracer
}
