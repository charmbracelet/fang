package fang

import (
	"errors"
	"testing"
)

func TestIsUsageError(t *testing.T) {
	for _, err := range []string{
		"flag needs an argument: 'foo'",
		"unknown flag: 'bar'",
		"unknown shorthand flag: 'b'",
		"unknown command 'baz'",
		"invalid argument 'qux'",
	} {
		t.Run(err, func(t *testing.T) {
			if !isUsageError(errors.New(err)) {
				t.Errorf("expected %q to be a usage error", err)
			}
		})
	}

	t.Run("something else", func(t *testing.T) {
		if isUsageError(errors.New("this is not a usage error")) {
			t.Errorf("expected 'this is not a usage error' to not be a usage error")
		}
	})
}
