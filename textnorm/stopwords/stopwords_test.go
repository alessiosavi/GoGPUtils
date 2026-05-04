package stopwords

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestEnglishContainsAnchors(t *testing.T) {
	set := English()
	for _, w := range []string{"the", "a", "an", "is", "are", "of", "with", "and", "or", "to"} {
		if _, ok := set[w]; !ok {
			t.Errorf("English missing anchor word %q", w)
		}
	}
}

func TestFrenchContainsAnchors(t *testing.T) {
	set := French()
	for _, w := range []string{"le", "la", "les", "et", "de"} {
		if _, ok := set[w]; !ok {
			t.Errorf("French missing anchor word %q", w)
		}
	}
}

func TestItalianContainsAnchors(t *testing.T) {
	set := Italian()
	for _, w := range []string{"il", "la", "i", "le", "e", "di"} {
		if _, ok := set[w]; !ok {
			t.Errorf("Italian missing anchor word %q", w)
		}
	}
}

// TestSingletonsCacheReturnSameMap verifies the lazy-loaded singletons
// return the same underlying map across calls. Identity is checked via
// reflect.Value.Pointer(), which exposes the map header.
func TestSingletonsCacheReturnSameMap(t *testing.T) {
	cases := []struct {
		name string
		fn   func() map[string]struct{}
	}{
		{"English", English},
		{"French", French},
		{"Italian", Italian},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a := tc.fn()
			b := tc.fn()
			pa := reflect.ValueOf(a).Pointer()
			pb := reflect.ValueOf(b).Pointer()
			if pa != pb {
				t.Fatalf("%s() returned different maps on successive calls (%x vs %x)", tc.name, pa, pb)
			}
		})
	}
}

func TestLoadFromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "custom")
	contents := "" +
		"# comment line\n" +
		"hello\n" +
		"\n" +
		"  world  \n" +
		"hello\n" + // duplicate, collapses
		"# trailing comment\n"
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := LoadFromFile("custom", path)
	if err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}
	want := map[string]struct{}{"hello": {}, "world": {}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("LoadFromFile() = %v, want %v", got, want)
	}
}

func TestLoadFromFileMissingPath(t *testing.T) {
	_, err := LoadFromFile("custom", filepath.Join(t.TempDir(), "does-not-exist"))
	if err == nil {
		t.Fatal("LoadFromFile() with missing path: want error, got nil")
	}
}

func TestLoadFromList(t *testing.T) {
	got := LoadFromList("custom", []string{"alpha", "  beta  ", "alpha", "", "gamma"})
	want := map[string]struct{}{"alpha": {}, "beta": {}, "gamma": {}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("LoadFromList() = %v, want %v", got, want)
	}
}

func TestLoadFromListEmpty(t *testing.T) {
	got := LoadFromList("custom", nil)
	if got == nil {
		t.Fatal("LoadFromList(nil) returned nil, want non-nil empty map")
	}
	if len(got) != 0 {
		t.Fatalf("LoadFromList(nil) size = %d, want 0", len(got))
	}
}

func TestCleanAllStopwordsUnionsBuiltins(t *testing.T) {
	got := CleanAllStopwords([]string{"english", "french"})
	if len(got) <= len(English()) {
		t.Fatalf("union size %d should exceed English alone (%d)", len(got), len(English()))
	}
	// Anchor words from both languages must appear in the union.
	for _, w := range []string{"the", "and", "le", "et"} {
		if _, ok := got[w]; !ok {
			t.Errorf("CleanAllStopwords union missing %q", w)
		}
	}
}

func TestCleanAllStopwordsIsCaseInsensitiveAndTrimmed(t *testing.T) {
	got := CleanAllStopwords([]string{"  English  ", "ITALIAN"})
	if len(got) <= len(English()) {
		t.Fatalf("union should include Italian additions; size %d, English %d", len(got), len(English()))
	}
}

func TestCleanAllStopwordsIgnoresUnknownLanguages(t *testing.T) {
	got := CleanAllStopwords([]string{"klingon", "english"})
	want := English()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("CleanAllStopwords([klingon, english]) should equal English() alone")
	}
}

func TestCleanAllStopwordsEmptyInput(t *testing.T) {
	got := CleanAllStopwords(nil)
	if got == nil {
		t.Fatal("CleanAllStopwords(nil) returned nil, want non-nil empty map")
	}
	if len(got) != 0 {
		t.Fatalf("CleanAllStopwords(nil) size = %d, want 0", len(got))
	}
}

func TestUnionMergesSetsWithoutDuplicates(t *testing.T) {
	a := map[string]struct{}{"x": {}, "y": {}}
	b := map[string]struct{}{"y": {}, "z": {}}
	got := Union(a, b)
	if len(got) != 3 {
		t.Fatalf("Union size = %d, want 3", len(got))
	}
	for _, w := range []string{"x", "y", "z"} {
		if _, ok := got[w]; !ok {
			t.Errorf("Union missing %q", w)
		}
	}
}

func TestUnionWithZeroInputsReturnsEmptyMap(t *testing.T) {
	got := Union()
	if got == nil {
		t.Fatal("Union() returned nil, want non-nil empty map")
	}
	if len(got) != 0 {
		t.Fatalf("Union() size = %d, want 0", len(got))
	}
}

func TestUnionDoesNotMutateInputs(t *testing.T) {
	a := map[string]struct{}{"x": {}}
	b := map[string]struct{}{"y": {}}
	out := Union(a, b)
	out["z"] = struct{}{}
	if len(a) != 1 || len(b) != 1 {
		t.Fatalf("Union mutated inputs after output mutation: a=%v b=%v", a, b)
	}
	if _, ok := a["z"]; ok {
		t.Fatal("output mutation leaked into input a")
	}
	if _, ok := b["z"]; ok {
		t.Fatal("output mutation leaked into input b")
	}
}
