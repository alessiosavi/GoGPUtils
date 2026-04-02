package textnorm

import "testing"

func TestMixedNormalizationPipelineRegression(t *testing.T) {
	got, err := New().NormalizeUnicode().FoldCase().TrimSpace().CollapseWhitespace().SplitTokens().MapTokens(func(s string) string {
		return s
	}).JoinTokens(" ").Run("  Stra\u00dfe   Caf\u00e9  ")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if got != "strasse cafe" {
		t.Fatalf("Run() = %q, want %q", got, "strasse cafe")
	}
}
