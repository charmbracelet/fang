package fang_test

import (
	"bytes"
	"strings"
	"testing"

	"charm.land/fang/v2"
	"github.com/spf13/cobra"
)

func TestWithLongRenderer(t *testing.T) {
	root := &cobra.Command{
		Use:  "app",
		Long: "the-raw-long",
		Run:  func(*cobra.Command, []string) {},
	}

	var got string
	var stdout bytes.Buffer
	root.SetOut(&stdout)
	root.SetArgs([]string{"--help"})

	err := fang.Execute(t.Context(), root, fang.WithLongRenderer(func(s string) string {
		got = s
		return "renderedmarker"
	}))
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	if got != "the-raw-long" {
		t.Fatalf("renderer received %q, want the raw long", got)
	}
	if !strings.Contains(stdout.String(), "renderedmarker") {
		t.Fatalf("rendered text missing from help output:\n%s", stdout.String())
	}
}
