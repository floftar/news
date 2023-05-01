package internal

import (
	"testing"
	"time"
)

func TestAge(t *testing.T) {
	tests := []struct {
		now      time.Time
		compare  time.Time
		expected string
	}{
		{
			time.Date(2009, 9, 6, 16, 20, 0, 0, time.UTC),
			time.Date(2009, 9, 6, 16, 20, 0, 0, time.UTC),
			"0 minutes",
		},

		{
			time.Date(2009, 9, 6, 16, 21, 0, 0, time.UTC),
			time.Date(2009, 9, 6, 16, 20, 0, 0, time.UTC),
			"1 minute",
		},

		{
			time.Date(2009, 9, 6, 16, 59, 0, 0, time.UTC),
			time.Date(2009, 9, 6, 16, 0, 0, 0, time.UTC),
			"59 minutes",
		},

		{
			time.Date(2009, 9, 6, 17, 0, 0, 0, time.UTC),
			time.Date(2009, 9, 6, 16, 0, 0, 0, time.UTC),
			"1 hour",
		},

		{
			time.Date(2009, 9, 6, 17, 1, 0, 0, time.UTC),
			time.Date(2009, 9, 6, 16, 0, 0, 0, time.UTC),
			"1 hour",
		},

		{
			time.Date(2009, 9, 6, 18, 0, 0, 0, time.UTC),
			time.Date(2009, 9, 6, 16, 0, 0, 0, time.UTC),
			"2 hours",
		},
	}

	for _, test := range tests {
		if actual := GetAge(test.now, test.compare); actual != test.expected {
			t.Errorf("GetAge(%q, %q) = %s, expected = %s",
				test.now, test.compare, actual, test.expected)
		}
	}
}
