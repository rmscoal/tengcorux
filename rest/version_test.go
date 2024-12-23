package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	assert.Equal(t, "v0.1.2", Version())
}
