package textnorm

import "testing"

func TestPreserveMeaningPunct(t *testing.T) {
	cases := []struct{ name, in, want string }{
		{"decimal kept", "size 4.5", "size 4.5"},
		{"decimal differs from int", "size 45", "size 45"},
		{"thousands comma kept", "1,000", "1,000"},
		{"separators not unified", "1.000", "1.000"},
		{"space thousands untouched", "1 000", "1 000"},
		{"trailing plus kept", "galaxy s22+", "galaxy s22+"},
		{"cpp keeps plus chain", "c++", "c++"},
		{"interior plus dropped", "a+b", "a b"},
		{"percent after digit kept", "100% cotton", "100% cotton"},
		{"percent standalone dropped", "% off", "  off"},
		{"dot not between digits dropped", "end. next", "end  next"},
		{"comma not between digits dropped", "red, blue", "red  blue"},
		{"other punct to space not deleted", "2.5-GHz", "2.5 GHz"},
		{"pipe to space", "a | b", "a   b"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := New().PreserveMeaningPunct().Run(c.in)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != c.want {
				t.Fatalf("PreserveMeaningPunct(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}
