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

func TestDedupTokensPreservesFirstSeenOrder(t *testing.T) {
	tokens, err := New().SplitTokens().DedupTokens().Run("alpha beta alpha gamma beta")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	want := []string{"alpha", "beta", "gamma"}
	if !reflect.DeepEqual(tokens, want) {
		t.Fatalf("Run() = %#v, want %#v", tokens, want)
	}
}

func TestDedupTokensIsCaseSensitive(t *testing.T) {
	tokens, err := New().SplitTokens().DedupTokens().Run("Red red RED Red")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	want := []string{"Red", "red", "RED"}
	if !reflect.DeepEqual(tokens, want) {
		t.Fatalf("Run() = %#v, want %#v", tokens, want)
	}
}

func TestDedupTokensWithFoldCaseUpstream(t *testing.T) {
	got, err := New().FoldCase().SplitTokens().DedupTokens().JoinTokens(" ").Run("Red red RED")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if got != "red" {
		t.Fatalf("Run() = %q, want %q", got, "red")
	}
}

func TestDedupTokensEmptyInput(t *testing.T) {
	tokens, err := New().SplitTokens().DedupTokens().Run("")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if len(tokens) != 0 {
		t.Fatalf("Run() = %#v, want []", tokens)
	}
}

func TestRemoveStopwordsDropsListedTokens(t *testing.T) {
	stop := map[string]struct{}{"the": {}, "a": {}, "of": {}}
	got, err := New().SplitTokens().RemoveStopwords(stop).JoinTokens(" ").Run("the cat sat on the mat of doom")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if got != "cat sat on mat doom" {
		t.Fatalf("Run() = %q, want %q", got, "cat sat on mat doom")
	}
}

func TestRemoveStopwordsNilSetIsNoOp(t *testing.T) {
	got, err := New().SplitTokens().RemoveStopwords(nil).JoinTokens(" ").Run("a b c")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if got != "a b c" {
		t.Fatalf("Run() = %q, want %q", got, "a b c")
	}
}

func TestRemoveStopwordsCaseSensitive(t *testing.T) {
	stop := map[string]struct{}{"the": {}}
	got, err := New().SplitTokens().RemoveStopwords(stop).JoinTokens(" ").Run("The the THE the")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if got != "The THE" {
		t.Fatalf("Run() = %q, want %q", got, "The THE")
	}
}

func TestRemoveStopwordsEmptyInput(t *testing.T) {
	stop := map[string]struct{}{"the": {}}
	tokens, err := New().SplitTokens().RemoveStopwords(stop).Run("")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if len(tokens) != 0 {
		t.Fatalf("Run() = %#v, want []", tokens)
	}
}

func TestDedupTokensSingleToken(t *testing.T) {
	tokens, err := New().SplitTokens().DedupTokens().Run("only")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	want := []string{"only"}
	if !reflect.DeepEqual(tokens, want) {
		t.Fatalf("Run() = %#v, want %#v", tokens, want)
	}
}

func TestDedupTokensAllDuplicates(t *testing.T) {
	tokens, err := New().SplitTokens().DedupTokens().Run("go go go")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	want := []string{"go"}
	if !reflect.DeepEqual(tokens, want) {
		t.Fatalf("Run() = %#v, want %#v", tokens, want)
	}
}

func TestRemoveStopwordsEmptySetIsNoOp(t *testing.T) {
	got, err := New().SplitTokens().RemoveStopwords(map[string]struct{}{}).JoinTokens(" ").Run("a b c")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if got != "a b c" {
		t.Fatalf("Run() = %q, want %q", got, "a b c")
	}
}
