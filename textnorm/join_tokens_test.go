package textnorm

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"
)

func TestJoinTokensMatchesTokenRun(t *testing.T) {
	appendWord := New().Then(func(s string) (string, error) { return s + " x", nil })
	for _, tc := range []struct {
		name  string
		pipe  TokenPipeline
		input string
	}{
		{"zero", TokenPipeline{}, "a b"},
		{"empty-input", New().SplitTokens(), ""},
		{"whitespace", New().SplitTokens(), " \t\n\u2003 "},
		{"single-token", New().SplitTokens(), " a "},
		{"append-word", appendWord.SplitTokens(), "a"},
		{"html-entity", New().DecodeHTMLEntities().SplitTokens(), "&amp;amp;"},
		{"map-filter", appendWord.SplitTokens().MapTokens(strings.ToUpper).FilterTokens(func(s string) bool { return s != "B" }), "a b"},
		{"dedup-stopwords", appendWord.SplitTokens().DedupTokens().RemoveStopwords(map[string]struct{}{"a": {}}), "a b b"},
		{"all-filtered", appendWord.SplitTokens().FilterTokens(func(string) bool { return false }), "a"},
		{"nil-tokens", appendWord.SplitTokens().Then(func([]string) ([]string, error) { return nil, nil }), "a"},
		{"empty-tokens", appendWord.SplitTokens().Then(func([]string) ([]string, error) { return []string{}, nil }), "a"},
		{"empty-token-elements", New().SplitTokens().Then(func([]string) ([]string, error) { return []string{"", "a", ""}, nil }), ""},
		{"token-bytes", New().SplitTokens().Then(func([]string) ([]string, error) { return []string{"a b", "\xff\x00", "界"}, nil }), ""},
		{"nested", appendWord.SplitTokens().JoinTokens(" ").SplitTokens(), "a"},
	} {
		for _, sep := range []string{"", " ", "|", "::", "界", "\x00\xff"} {
			t.Run(fmt.Sprintf("%s/sep=%q", tc.name, sep), func(t *testing.T) {
				tokens, err := tc.pipe.Run(tc.input)
				if err != nil {
					t.Fatal(err)
				}
				want := strings.Join(tokens, sep)
				got, err := tc.pipe.JoinTokens(sep).Run(tc.input)
				if err != nil || got != want {
					t.Fatalf("Run(%q) = %q, %v; token Run + Join = %q", tc.input, got, err, want)
				}
			})
		}
	}
}

func TestJoinTokensRunsCallbacksOncePerRun(t *testing.T) {
	var sourceCalls, mapCalls, filterCalls, downstreamCalls int
	var trace []string
	pipe := New().Then(func(s string) (string, error) {
		sourceCalls++
		trace = append(trace, "source:"+s)
		return s + " x", nil
	}).SplitTokens().MapTokens(func(s string) string {
		mapCalls++
		trace = append(trace, "map:"+s)
		return fmt.Sprintf("%s%d", s, mapCalls)
	}).FilterTokens(func(s string) bool {
		filterCalls++
		trace = append(trace, "filter:"+s)
		return filterCalls%2 == 1
	}).JoinTokens("|").Then(func(s string) (string, error) {
		downstreamCalls++
		trace = append(trace, "downstream:"+s)
		return "[" + s + "]", nil
	})
	if sourceCalls != 0 || mapCalls != 0 || filterCalls != 0 || downstreamCalls != 0 {
		t.Fatal("building a pipeline ran callbacks")
	}
	for run := 1; run <= 2; run++ {
		trace = nil
		got, err := pipe.Run("a")
		a, x := fmt.Sprintf("a%d", 2*run-1), fmt.Sprintf("x%d", 2*run)
		if err != nil || got != "["+a+"]" {
			t.Errorf("run %d = %q, %v; want [%s], nil", run, got, err, a)
		}
		if sourceCalls != run || mapCalls != 2*run || filterCalls != 2*run || downstreamCalls != run {
			t.Errorf("run %d: calls source/map/filter/downstream = %d/%d/%d/%d; want %d/%d/%d/%d", run, sourceCalls, mapCalls, filterCalls, downstreamCalls, run, 2*run, 2*run, run)
		}
		wantTrace := []string{"source:a", "map:a", "map:x", "filter:" + a, "filter:" + x, "downstream:" + a}
		if !slices.Equal(trace, wantTrace) {
			t.Errorf("run %d: trace = %q, want %q", run, trace, wantTrace)
		}
	}
}

func TestJoinTokensNestedRunsSourceOnce(t *testing.T) {
	for depth := 1; depth <= 3; depth++ {
		t.Run(fmt.Sprint(depth), func(t *testing.T) {
			sourceCalls := 0
			mapCalls := make([]int, depth)
			pipe := New().Then(func(s string) (string, error) {
				sourceCalls++
				return s + " x", nil
			})
			for level := range depth {
				pipe = pipe.SplitTokens().MapTokens(func(s string) string {
					mapCalls[level]++
					return s
				}).JoinTokens(" ")
			}
			for run := 1; run <= 2; run++ {
				got, err := pipe.Run("a")
				if err != nil || got != "a x" || sourceCalls != run {
					t.Errorf("run %d = %q, %v; source calls = %d; want a x, nil, %d", run, got, err, sourceCalls, run)
				}
				for level, calls := range mapCalls {
					if calls != 2*run {
						t.Errorf("run %d level %d: mapper calls = %d, want %d", run, level, calls, 2*run)
					}
				}
			}
		})
	}
}

func TestJoinTokensPropagatesFailuresInOrder(t *testing.T) {
	stages := []string{"source-first", "source-last", "token-first", "map", "filter", "token-last", "downstream-first", "downstream-last"}
	for index, failAt := range stages {
		for _, mode := range []string{"error", "panic"} {
			if mode == "error" && (failAt == "map" || failAt == "filter") {
				continue // Map/filter callbacks do not return errors.
			}
			t.Run(failAt+"/"+mode, func(t *testing.T) {
				// Use a wrapped sentinel to distinguish exact identity from errors.Is.
				failure := fmt.Errorf("callback: %w", errors.New("failure"))
				var trace []string
				record := func(name string) error {
					trace = append(trace, name)
					if name == failAt {
						if mode == "panic" {
							panic(failure)
						}
						return failure
					}
					return nil
				}
				pipe := New().Then(func(s string) (string, error) {
					return s, record("source-first")
				}).Then(func(s string) (string, error) {
					return s, record("source-last")
				}).SplitTokens().Then(func(tokens []string) ([]string, error) {
					return tokens, record("token-first")
				}).MapTokens(func(s string) string {
					_ = record("map")
					return s
				}).FilterTokens(func(string) bool {
					_ = record("filter")
					return true
				}).Then(func(tokens []string) ([]string, error) {
					return tokens, record("token-last")
				}).JoinTokens("|").Then(func(s string) (string, error) {
					return s, record("downstream-first")
				}).Then(func(s string) (string, error) {
					return s, record("downstream-last")
				})
				if mode == "panic" {
					func() {
						defer func() {
							if got := recover(); got != failure { //nolint:errorlint // A panic must preserve the exact value, not just errors.Is equivalence.
								t.Errorf("panic = %v, want original failure %v", got, failure)
							}
						}()
						_, _ = pipe.Run("a")
					}()
				} else {
					got, err := pipe.Run("a")
					if got != "" || err != failure { //nolint:errorlint // Preserve exact first-pass error identity, including wrapping.
						t.Errorf("Run = %q, %v; want empty output and original failure %v", got, err, failure)
					}
				}
				if want := stages[:index+1]; !slices.Equal(trace, want) {
					t.Errorf("trace = %q, want %q", trace, want)
				}
			})
		}
	}
}

func TestJoinTokensDoesNotTriggerSecondSourceFailure(t *testing.T) {
	for _, mode := range []string{"error", "panic"} {
		t.Run(mode, func(t *testing.T) {
			calls := 0
			pipe := New().Then(func(s string) (string, error) {
				calls++
				if calls == 2 {
					if mode == "panic" {
						panic("second source execution")
					}
					return "", errors.New("second source execution")
				}
				return s + " x", nil
			}).SplitTokens().JoinTokens("|")
			defer func() {
				if p := recover(); p != nil {
					t.Errorf("unexpected panic: %v", p)
				}
			}()
			if got, err := pipe.Run("a"); err != nil || got != "a|x" || calls != 1 {
				t.Errorf("Run = %q, %v; source calls = %d; want a|x, nil, 1", got, err, calls)
			}
		})
	}
}

func TestJoinTokensBranchesAreIndependent(t *testing.T) {
	source := New().Then(func(s string) (string, error) { return s + " x", nil })
	tokens := source.SplitTokens().MapTokens(strings.ToUpper)
	joined := tokens.JoinTokens("|")
	left := joined.Then(func(s string) (string, error) { return s + "!", nil })
	right := joined.Then(func(s string) (string, error) { return "[" + s + "]", nil })
	filtered := tokens.FilterTokens(func(s string) bool { return s != "X" }).JoinTokens("/")
	mapped := tokens.MapTokens(func(s string) string { return s + "?" }).JoinTokens("/")
	sourceBranch := source.Then(func(s string) (string, error) { return s + " y", nil })
	for _, tc := range []struct {
		name string
		pipe Pipeline
		want string
	}{
		{"joined", joined, "A|X"},
		{"left", left, "A|X!"},
		{"right", right, "[A|X]"},
		{"filtered", filtered, "A"},
		{"mapped", mapped, "A?/X?"},
		{"source-branch", sourceBranch, "a x y"},
		{"source", source, "a x"},
		{"other-separator", tokens.JoinTokens("::"), "A::X"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Concurrent runs also exercise sharing the captured token pipeline.
			t.Parallel()
			for range 2 {
				if got, err := tc.pipe.Run("a"); err != nil || got != tc.want {
					t.Fatalf("Run = %q, %v; want %q, nil", got, err, tc.want)
				}
			}
		})
	}
	if got, err := tokens.Run("a"); err != nil || !slices.Equal(got, []string{"A", "X"}) {
		t.Fatalf("token Run = %q, %v; want [A X], nil", got, err)
	}
}
