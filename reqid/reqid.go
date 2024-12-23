package reqid

import (
	"context"
	cryptorand "crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	mathrand "math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

func init() {
	// Assert that a cryptographically secure PRNG is available.
	// Panic otherwise.
	buf := make([]byte, 1)

	_, err := io.ReadFull(cryptorand.Reader, buf)
	if err != nil {
		panic(fmt.Sprintf("crypto/rand is unavailable: Read() failed with %#v", err))
	}
}

// contextKey will be the type of key used in context.Context
type contextKey struct{}

var ContextKey contextKey

const HeaderKey = "X-Request-Id"

// GenerateRequestId creates a new cryptographically random string. If failed,
// calls GenerateUUID instead.
// Adapted from https://elithrar.github.io/article/generating-secure-random-numbers-crypto-rand/
func GenerateRequestId() string {
	seed := time.Now().UnixNano()
	lengthRandomizer := mathrand.New(mathrand.NewSource(seed))
	reqIdLength := lengthRandomizer.Intn(16) + 16 // Max: 32, min: 16

	b := make([]byte, reqIdLength)
	_, err := cryptorand.Read(b)
	if err != nil {
		return GenerateUUID()
	}

	encoded := base64.URLEncoding.EncodeToString(b)
	return strings.ReplaceAll(encoded, "=", "")
}

// GenerateUUID generates uuid V7 with V4 as fallback.
func GenerateUUID() string {
	v7, err := uuid.NewV7()
	if err != nil {
		// Fallback use v4
		return uuid.NewString()
	}
	return v7.String()
}

// RetrieveFromContext extracts the value from context by looking
// for the context key.
func RetrieveFromContext(ctx context.Context) string {
	val := ctx.Value(ContextKey)
	if val == nil {
		return ""
	}
	return val.(string)
}

// RetrieveFromHttpHeader gets the HeaderKey value for header.
func RetrieveFromHttpHeader(header http.Header) string {
	return header.Get(HeaderKey)
}

// Inject automatically generates and inserts into the
// context only if the previous value has been set.
func Inject(ctx context.Context) context.Context {
	if RetrieveFromContext(ctx) == "" {
		ctx = context.WithValue(ctx, ContextKey, GenerateRequestId())
	}

	return ctx
}

func InjectValue(ctx context.Context, value string) context.Context {
	if RetrieveFromContext(ctx) == "" {
		ctx = context.WithValue(ctx, ContextKey, value)
	}

	return ctx
}
