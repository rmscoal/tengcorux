package skywalking

import "testing"

func TestStringToSpanID(t *testing.T) {
	tests := []struct {
		input    string
		expected int32
	}{
		{
			input:    "0",
			expected: 0,
		},
		{
			input:    "1",
			expected: 1,
		},
		{
			input:    "0.4",
			expected: 0,
		},
		{
			input:    "1s3hb",
			expected: 0,
		},
		{
			input:    "9321",
			expected: 9321,
		},
	}

	for _, test := range tests {
		id := stringToSpanID(test.input)
		if id != test.expected {
			t.Errorf("expected %d, got %d", test.expected, id)
		}
	}
}
