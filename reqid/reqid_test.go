package reqid

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ////////////////////
// Tests
// ////////////////////
func TestGenerateRequestId(t *testing.T) {
	reqID := GenerateRequestId()
	assert.NotEmpty(t, reqID,
		"GenerateRequestId() should not return an empty string")
}

func TestGenerateUUID(t *testing.T) {
	uuid := GenerateUUID()
	assert.NotEmpty(t, uuid, "GenerateUUID() should not return an empty string")
}

func TestContext(t *testing.T) {
	ctx := context.Background()
	ctx = Inject(ctx)

	reqID := RetrieveFromContext(ctx)
	assert.NotEmpty(t, reqID,
		"RetrieveFromContext() should not return an empty string after Inject()")
	t.Log("Request ID in context:", reqID)
}

func TestRetrieveFromHttpHeader(t *testing.T) {
	header := http.Header{}
	expectedReqID := GenerateRequestId()
	header.Set(HeaderKey, expectedReqID)

	actualReqID := RetrieveFromHttpHeader(header)
	assert.NotEmpty(t, actualReqID,
		"RetrieveFromHttpHeader() should not return an empty string when the header is set")
	assert.Equal(t, expectedReqID, actualReqID,
		"Retrieved Request ID should match the expected value")
	t.Log("Request ID in header:", actualReqID)
}

func TestInjectValue(t *testing.T) {
	t.Run("InjectValue with empty context", func(t *testing.T) {
		// Create a new empty context
		ctx := context.Background()
		expectedValue := "test-request-id"

		// Inject the value into the context
		ctx = InjectValue(ctx, expectedValue)

		// Retrieve the value from the context
		actualValue := RetrieveFromContext(ctx)

		// Assert that the injected value is correctly stored
		assert.Equal(t, expectedValue, actualValue,
			"Expected and actual values should match")
	})

	t.Run("InjectValue with pre-existing context value", func(t *testing.T) {
		// Create a context with an existing value
		existingValue := "existing-request-id"
		ctx := context.WithValue(context.Background(), ContextKey,
			existingValue)

		// Attempt to inject a new value
		newValue := "new-request-id"
		ctx = InjectValue(ctx, newValue)

		// Retrieve the value from the context
		actualValue := RetrieveFromContext(ctx)

		// Assert that the existing value is not overwritten
		assert.Equal(t, existingValue, actualValue,
			"Existing value should not be overwritten")
	})

	t.Run("InjectValue with empty string value", func(t *testing.T) {
		// Create a new empty context
		ctx := context.Background()
		expectedValue := ""

		// Inject the empty value into the context
		ctx = InjectValue(ctx, expectedValue)

		// Retrieve the value from the context
		actualValue := RetrieveFromContext(ctx)

		// Assert that the empty value is correctly stored
		assert.Equal(t, expectedValue, actualValue,
			"Expected and actual values should match")
	})
}

//////////////////////
// Benchmarks
//////////////////////

func BenchmarkGenerateRequestId(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateRequestId()
	}
}

func BenchmarkGenerateUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateUUID()
	}
}
