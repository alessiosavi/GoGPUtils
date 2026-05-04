// Package stopwords provides language-specific stopword sets sourced from
// the NLTK stopwords corpus (https://github.com/nltk/nltk_data).
//
// Built-in sets are exposed as accessor functions (English, French,
// Italian) that lazily parse embedded text files on first call and cache
// the result. The returned map[string]struct{} is shared across callers
// and MUST NOT be mutated. Use Union to combine sets safely; it always
// returns a new independent map.
//
// Unicode normalization precondition: the non-ASCII entries in French
// (e.g. "à") and Italian (e.g. "è", "più") are stored in composed (NFC)
// form. Tokens fed to RemoveStopwords must be in the same normalization
// form. The textnorm presets (SearchPreset, CanonicalPreset, DBSafePreset)
// already include NormalizeUnicode upstream.
//
// To refresh the bundled data files from upstream:
//
//	wget https://github.com/nltk/nltk_data/raw/refs/heads/gh-pages/packages/corpora/stopwords.zip
//	unzip -o stopwords.zip
//	cp stopwords/{english,french,italian} textnorm/stopwords/data/
package stopwords

import (
	"bufio"
	"embed"
	"io"
	"os"
	"strings"
	"sync"
)

//go:embed data/english data/french data/italian
var dataFS embed.FS

var (
	englishOnce sync.Once
	englishSet  map[string]struct{}

	frenchOnce sync.Once
	frenchSet  map[string]struct{}

	italianOnce sync.Once
	italianSet  map[string]struct{}
)

// English returns the embedded English stopword set (NLTK corpus). The
// first call parses the embedded file; subsequent calls return the same
// cached map. The returned map is shared and MUST NOT be mutated.
func English() map[string]struct{} {
	englishOnce.Do(func() {
		englishSet = mustLoadEmbedded("data/english")
	})
	return englishSet
}

// French returns the embedded French stopword set (NLTK corpus). The first
// call parses the embedded file; subsequent calls return the same cached
// map. The returned map is shared and MUST NOT be mutated.
func French() map[string]struct{} {
	frenchOnce.Do(func() {
		frenchSet = mustLoadEmbedded("data/french")
	})
	return frenchSet
}

// Italian returns the embedded Italian stopword set (NLTK corpus). The
// first call parses the embedded file; subsequent calls return the same
// cached map. The returned map is shared and MUST NOT be mutated.
func Italian() map[string]struct{} {
	italianOnce.Do(func() {
		italianSet = mustLoadEmbedded("data/italian")
	})
	return italianSet
}

// LoadFromFile reads a newline-delimited word list from filepath and
// returns a stopword set independent of the built-in singletons. Empty
// lines and lines beginning with '#' are skipped; remaining lines are
// trimmed of surrounding whitespace.
//
// The lang argument is informational and is not used to dispatch into a
// built-in set; callers wanting a built-in should call English/French/
// Italian directly. Use this helper to ship a custom override.
func LoadFromFile(lang string, filepath string) (map[string]struct{}, error) {
	_ = lang
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return parseWords(f)
}

// LoadFromList builds a stopword set from a slice of words. Empty strings
// (after trimming) are skipped; duplicates collapse to a single entry.
// The returned map is independent of the built-in singletons.
//
// The lang argument is informational and is not used to dispatch into a
// built-in set; callers wanting a built-in should call English/French/
// Italian directly.
func LoadFromList(lang string, words []string) map[string]struct{} {
	_ = lang
	out := make(map[string]struct{}, len(words))
	for _, w := range words {
		w = strings.TrimSpace(w)
		if w == "" {
			continue
		}
		out[w] = struct{}{}
	}
	return out
}

// CleanAllStopwords returns the union of the built-in stopword sets for
// the requested languages. Recognized values (case-insensitive, trimmed):
// "english", "french", "italian". Unknown languages are silently ignored.
// With zero recognized inputs, returns a non-nil empty map.
func CleanAllStopwords(languages []string) map[string]struct{} {
	sets := make([]map[string]struct{}, 0, len(languages))
	for _, lang := range languages {
		switch strings.ToLower(strings.TrimSpace(lang)) {
		case "english":
			sets = append(sets, English())
		case "french":
			sets = append(sets, French())
		case "italian":
			sets = append(sets, Italian())
		}
	}
	return Union(sets...)
}

// Union returns a new set containing every key found in any of the input
// sets. It does not mutate any input. With zero inputs it returns a
// non-nil empty map.
func Union(sets ...map[string]struct{}) map[string]struct{} {
	total := 0
	for _, s := range sets {
		total += len(s)
	}
	out := make(map[string]struct{}, total)
	for _, s := range sets {
		for k := range s {
			out[k] = struct{}{}
		}
	}
	return out
}

func mustLoadEmbedded(name string) map[string]struct{} {
	f, err := dataFS.Open(name)
	if err != nil {
		panic("stopwords: open embedded " + name + ": " + err.Error())
	}
	defer f.Close()
	out, err := parseWords(f)
	if err != nil {
		panic("stopwords: parse embedded " + name + ": " + err.Error())
	}
	return out
}

func parseWords(r io.Reader) (map[string]struct{}, error) {
	out := make(map[string]struct{})
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		out[line] = struct{}{}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
