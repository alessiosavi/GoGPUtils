package stringutil

import (
	"math"
	"slices"
	"testing"
	"time"
)

func TestCosineSimilarity_UsesAccurateMagnitude(t *testing.T) {
	got := CosineSimilarity("ab", "a", 1)
	want := 1 / math.Sqrt(2)
	if got != want {
		t.Errorf("CosineSimilarity(ab, a, 1) = %.17g, want %.17g", got, want)
	}
}

func TestBetween_PreservesEmptyDelimiterResults(t *testing.T) {
	tests := []struct {
		name, start, end, want string
	}{
		{"empty start", "", "c", "ab"},
		{"empty end", "a", "", ""},
		{"both empty", "", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, ok := Between("abc", tt.start, tt.end); got != tt.want || !ok {
				t.Errorf("Between(%q, %q, %q) = %q, %v; want %q, true", "abc", tt.start, tt.end, got, ok, tt.want)
			}
		})
	}
}

func TestBetweenAll_EmptyDelimitersAndProgress(t *testing.T) {
	tests := []struct {
		name, s, start, end string
		want                []string
	}{
		{"empty start", "a,b,c", "", ",", []string{"a", "b"}},
		{"empty end", "a,b,c", ",", "", []string{"", ""}},
		{"empty start with bracket", "[value]", "", "]", []string{"[value"}},
		{"empty end with bracket", "[value]", "[", "", []string{""}},
		{"empty start consumes last byte", ",", "", ",", []string{""}},
		{"empty end consumes last byte", ",", ",", "", []string{""}},
		{"empty input and start", "", "", ",", nil},
		{"empty input and end", "", ",", "", nil},
		{"both empty", "value", "", "", nil},
		{"empty input and delimiters", "", "", "", nil},
		{"adjacent", "[][one][two]", "[", "]", []string{"", "one", "two"}},
		{"identical delimiters", "|a||b|", "|", "|", []string{"a", "b"}},
		{"missing end", "[one][two", "[", "]", []string{"one"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			done := make(chan []string, 1)
			go func() {
				done <- BetweenAll(tt.s, tt.start, tt.end)
			}()

			select {
			case got := <-done:
				if !slices.Equal(got, tt.want) || (tt.want == nil && got != nil) {
					t.Errorf("BetweenAll() = %#v, want %#v", got, tt.want)
				}
			case <-time.After(time.Second):
				t.Fatal("BetweenAll did not terminate within one second")
			}
		})
	}
}
