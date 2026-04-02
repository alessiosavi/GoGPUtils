package textnorm

import (
	"reflect"
	"strings"
	"testing"
)

func TestSplitTokensUsesWhitespaceBoundaries(t *testing.T) {
	tokens, err := New().SplitTokens().Run("  go,   gophers are  fun  ")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	want := []string{"go,", "gophers", "are", "fun"}
	if !reflect.DeepEqual(tokens, want) {
		t.Fatalf("Run() = %#v, want %#v", tokens, want)
	}
}

func TestTokenPipelineMapsFiltersAndJoins(t *testing.T) {
	got, err := New().SplitTokens().MapTokens(strings.ToUpper).FilterTokens(func(s string) bool {
		return s != "SKIP"
	}).JoinTokens("|").Run("go skip text")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if got != "GO|TEXT" {
		t.Fatalf("Run() = %q, want %q", got, "GO|TEXT")
	}
}

func TestTokenPipelineReturnsTokens(t *testing.T) {
	tokens, err := New().SplitTokens().MapTokens(func(s string) string { return s + "!" }).Run("a b")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	want := []string{"a!", "b!"}
	if !reflect.DeepEqual(tokens, want) {
		t.Fatalf("Run() = %#v, want %#v", tokens, want)
	}
}
