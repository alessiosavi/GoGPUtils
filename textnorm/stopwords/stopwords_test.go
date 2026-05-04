package stopwords

import "testing"

func TestEnglishContainsAnchors(t *testing.T) {
	for _, w := range []string{"the", "a", "an", "is", "are", "of", "with", "and", "or", "to"} {
		if _, ok := English[w]; !ok {
			t.Errorf("English missing anchor word %q", w)
		}
	}
}

func TestFrenchContainsAnchors(t *testing.T) {
	for _, w := range []string{"le", "la", "les", "et", "de"} {
		if _, ok := French[w]; !ok {
			t.Errorf("French missing anchor word %q", w)
		}
	}
}

func TestItalianContainsAnchors(t *testing.T) {
	for _, w := range []string{"il", "la", "i", "le", "e", "di"} {
		if _, ok := Italian[w]; !ok {
			t.Errorf("Italian missing anchor word %q", w)
		}
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
	// Mutating the output must not bleed into the inputs.
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
