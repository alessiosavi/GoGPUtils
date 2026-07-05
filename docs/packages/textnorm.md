---
layout: default
title: textnorm
parent: Packages
nav_order: 4
---

# Package `textnorm`

Deterministic text normalization pipelines for Go.

```go
import "github.com/alessiosavi/GoGPUtils/textnorm"
```

## Overview

The `textnorm` package provides a **fluent, deterministic pipeline API** for normalizing text. Instead of applying ad-hoc string transformations, you compose a `Pipeline` of discrete **stages** that transform input text in a predictable, repeatable order.

### The Pipeline Concept

A `Pipeline` is an ordered sequence of `Stage` functions, where each stage is a pure transformation:

```go
type Stage func(string) (string, error)
```

Pipelines are **immutable** — calling `Then()` returns a new pipeline with the stage appended, leaving the original unchanged. This makes pipelines safe to reuse, share, and extend without side effects.

```go
base := textnorm.New()                          // empty pipeline
decorated := base.Then(myStage)                 // base is unchanged
result, err := decorated.Run("some text")       // execute all stages
```

The package also supports **token pipelines** for operations that work on word-level token slices (deduplication, stopword removal, etc.). Token pipelines seamlessly bridge back to string pipelines via `JoinTokens()`.

---

## Types

### `Stage`

```go
type Stage func(string) (string, error)
```

A `Stage` transforms input text and may return an error. Stages are the building blocks of a `Pipeline`.

### `Pipeline`

```go
type Pipeline struct {
    stages []Stage
}
```

`Pipeline` holds an ordered list of normalization stages. The zero value is valid (an empty pipeline that returns input unchanged).

### `TokenStage`

```go
type TokenStage func([]string) ([]string, error)
```

A `TokenStage` transforms a slice of tokens and may return an error. Used within `TokenPipeline`.

### `TokenPipeline`

```go
type TokenPipeline struct {
    source Pipeline
    stages []TokenStage
}
```

`TokenPipeline` holds ordered token stages derived from a string pipeline. Created by calling `SplitTokens()` on a `Pipeline`.

### `PresetOption`

```go
type PresetOption func(*presetConfig)
```

`PresetOption` customizes preset builders. Currently supports `WithWidthFold()`.

---

## Pipeline Functions

### `New`

```go
func New() Pipeline
```

Returns a new empty pipeline.

### `Pipeline.Then`

```go
func (p Pipeline) Then(stage Stage) Pipeline
```

Returns a new pipeline with the given stage appended. If `stage` is `nil`, returns the original pipeline unchanged.

### `Pipeline.Run`

```go
func (p Pipeline) Run(input string) (string, error)
```

Executes all stages in declaration order, passing the output of each stage as input to the next. Returns the final result or the first error encountered.

---

## Preset Functions

Presets are pre-configured pipelines for common use cases. All presets accept optional `PresetOption` arguments.

### `SearchPreset`

```go
func SearchPreset(opts ...PresetOption) Pipeline
```

Builds a **search-key pipeline** optimized for indexing and search:

1. `SanitizeUTF8()` — remove invalid bytes and NUL characters
2. `NormalizeUnicode()` — decompose and strip accents
3. `FoldWidth()` _(optional, via `WithWidthFold()`)_ — fold full-width characters
4. `FoldCase()` — full Unicode case folding
5. `FilterRunes()` — keep only letters, numbers, and spaces
6. `TrimSpace()` — remove leading/trailing whitespace
7. `CollapseWhitespace()` — collapse consecutive whitespace to single spaces
8. `SplitTokens()` — split into tokens
9. `JoinTokens(" ")` — join tokens with single spaces

**Example:**

```go
result, err := textnorm.SearchPreset().Run("  Café,   go!  gophers ")
// result == "cafe go gophers"
```

To deduplicate tokens or strip stopwords, compose on top:

```go
import "github.com/alessiosavi/GoGPUtils/textnorm/stopwords"

result, err := textnorm.SearchPreset().
    SplitTokens().
    DedupTokens().
    RemoveStopwords(stopwords.English()).
    JoinTokens(" ").
    Run("the cat sat on the mat")
```

### `CanonicalPreset`

```go
func CanonicalPreset(opts ...PresetOption) Pipeline
```

Builds a **general-purpose canonicalization pipeline**:

1. `SanitizeUTF8()` — remove invalid bytes and NUL characters
2. `NormalizeUnicode()` — decompose and strip accents
3. `FoldWidth()` _(optional, via `WithWidthFold()`)_ — fold full-width characters
4. `FoldCase()` — full Unicode case folding
5. `TrimSpace()` — remove leading/trailing whitespace
6. `CollapseWhitespace()` — collapse consecutive whitespace to single spaces

**Example:**

```go
result, err := textnorm.CanonicalPreset().Run("  Hello,   World!  ")
// result == "hello, world!"
```

### `DBSafePreset`

```go
func DBSafePreset(opts ...PresetOption) Pipeline
```

Builds a **persistence-safe normalization pipeline** for database storage:

1. `SanitizeUTF8()` — remove invalid bytes and NUL characters
2. `NormalizeUnicode()` — decompose and strip accents
3. `FoldWidth()` _(optional, via `WithWidthFold()`)_ — fold full-width characters
4. `TrimSpace()` — remove leading/trailing whitespace
5. `CollapseWhitespace()` — collapse consecutive whitespace to single spaces

Unlike `SearchPreset` and `CanonicalPreset`, this preset **preserves case and punctuation** — it only sanitizes and normalizes whitespace.

**Example:**

```go
input := string([]byte{'g', 'o', 0x00, 0xff, '!'})
result, err := textnorm.DBSafePreset().Run(input)
// result == "go!" (NUL removed, invalid byte replaced with replacement char)
```

### `WithWidthFold`

```go
func WithWidthFold() PresetOption
```

Enables explicit width folding in a preset pipeline. Full-width characters (e.g., `Ｇｏ`) are folded to their half-width equivalents (e.g., `Go`).

**Example:**

```go
result, err := textnorm.SearchPreset(textnorm.WithWidthFold()).Run("Ｇｏ")
// result == "go"
```

---

## String Pipeline Stages

### `SanitizeUTF8`

```go
func (p Pipeline) SanitizeUTF8() Pipeline
```

Appends a UTF-8 and NUL sanitization stage. Invalid UTF-8 byte sequences are replaced with the Unicode replacement character (`U+FFFD`, `�`). NUL bytes (`\x00`) are removed entirely. Valid UTF-8 strings without NUL bytes pass through unchanged.

**Example:**

```go
input := string([]byte{'g', 'o', 0xff, 0x00, '!', 0xfe})
result, err := textnorm.New().SanitizeUTF8().Run(input)
// result == "go\uFFFD!\uFFFD"
```

### `NormalizeUnicode`

```go
func (p Pipeline) NormalizeUnicode() Pipeline
```

Appends a Unicode normalization stage. Decomposes characters using NFD, strips diacritical marks (Unicode category `Mn`), then recomposes using NFC. This effectively removes accents from Latin characters.

**Example:**

```go
result, err := textnorm.New().NormalizeUnicode().Run("café")
// result == "cafe"
```

### `RemoveAccents`

```go
func (p Pipeline) RemoveAccents() Pipeline
```

Appends a diacritic-removal stage. This is an alias for `NormalizeUnicode()` — both perform the same transformation.

**Example:**

```go
result, err := textnorm.New().RemoveAccents().Run("naïve")
// result == "naive"
```

### `FoldCase`

```go
func (p Pipeline) FoldCase() Pipeline
```

Appends full Unicode case folding. This is more aggressive than simple lowercasing — it handles special cases like the German eszett (`ß` → `ss`).

**Example:**

```go
folded, _ := textnorm.New().FoldCase().Run("Straße")
// folded == "strasse"

lowered, _ := textnorm.New().Lower().Run("Straße")
// lowered == "straße"
```

### `Lower`

```go
func (p Pipeline) Lower() Pipeline
```

Appends Unicode-aware lowercasing using the `golang.org/x/text/cases` package with the undetermined language tag.

### `FoldWidth`

```go
func (p Pipeline) FoldWidth() Pipeline
```

Appends an explicit width-folding stage. Full-width characters (common in East Asian typography) are folded to their half-width equivalents.

**Example:**

```go
result, err := textnorm.New().FoldWidth().Run("Ｇｏ")
// result == "Go"
```

### `TrimSpace`

```go
func (p Pipeline) TrimSpace() Pipeline
```

Appends a trimming stage that removes leading and trailing Unicode whitespace using `strings.TrimSpace`.

### `CollapseWhitespace`

```go
func (p Pipeline) CollapseWhitespace() Pipeline
```

Appends a whitespace-collapsing stage. All consecutive whitespace sequences (spaces, tabs, newlines) are collapsed to a single space character.

**Example:**

```go
result, err := textnorm.New().TrimSpace().CollapseWhitespace().Run("\t  hello   world \n")
// result == "hello world"
```

### `MapRunes`

```go
func (p Pipeline) MapRunes(fn func(rune) rune) Pipeline
```

Appends a rune-mapping stage. Each rune in the input is transformed by `fn`. If `fn` is `nil`, the pipeline is returned unchanged.

**Example:**

```go
result, err := textnorm.New().MapRunes(func(r rune) rune {
    if unicode.IsLetter(r) || unicode.IsDigit(r) {
        return r
    }
    return '-'
}).Run("Go! 2")
// result == "Go--2"
```

### `FilterRunes`

```go
func (p Pipeline) FilterRunes(keep runes.Set) Pipeline
```

Appends a rune-filtering stage. Only runes matching the `keep` set are retained; all others are removed. If `keep` is `nil`, the pipeline is returned unchanged.

**Example:**

```go
keep := runes.Predicate(func(r rune) bool {
    return unicode.IsLetter(r) || unicode.IsDigit(r)
})
result, err := textnorm.New().FilterRunes(keep).Run("Go 1.23!")
// result == "Go123"
```

---

## Token Pipeline Stages

### `SplitTokens`

```go
func (p Pipeline) SplitTokens() TokenPipeline
```

Turns the current string pipeline into a token pipeline. The source pipeline runs first, then the result is split into tokens using `strings.Fields` (whitespace-separated).

**Example:**

```go
tokens, err := textnorm.New().SplitTokens().Run("  go,   gophers are  fun  ")
// tokens == []string{"go,", "gophers", "are", "fun"}
```

### `TokenPipeline.Then`

```go
func (tp TokenPipeline) Then(stage TokenStage) TokenPipeline
```

Returns a new token pipeline with a stage appended. If `stage` is `nil`, returns the original pipeline unchanged.

### `TokenPipeline.Run`

```go
func (tp TokenPipeline) Run(input string) ([]string, error)
```

Executes the source pipeline, tokenizes the result, and applies all token stages in order. Returns the final token slice or the first error encountered.

### `TokenPipeline.MapTokens`

```go
func (tp TokenPipeline) MapTokens(fn func(string) string) TokenPipeline
```

Returns a new token pipeline that maps every token through `fn`. If `fn` is `nil`, the pipeline is returned unchanged.

**Example:**

```go
tokens, err := textnorm.New().SplitTokens().MapTokens(strings.ToUpper).Run("a b")
// tokens == []string{"A", "B"}
```

### `TokenPipeline.FilterTokens`

```go
func (tp TokenPipeline) FilterTokens(fn func(string) bool) TokenPipeline
```

Returns a new token pipeline that keeps only tokens for which `fn` returns `true`. If `fn` is `nil`, the pipeline is returned unchanged.

**Example:**

```go
result, err := textnorm.New().SplitTokens().
    FilterTokens(func(s string) bool { return s != "skip" }).
    JoinTokens("|").
    Run("go skip text")
// result == "go|text"
```

### `DedupTokens`

```go
func (tp TokenPipeline) DedupTokens() TokenPipeline
```

Returns a new token pipeline that drops duplicate tokens, preserving the first occurrence. Comparison is plain string equality (case-sensitive). Pair with `FoldCase()` upstream for case-insensitive deduplication.

**Example:**

```go
tokens, err := textnorm.New().SplitTokens().DedupTokens().Run("alpha beta alpha gamma beta")
// tokens == []string{"alpha", "beta", "gamma"}
```

Case-insensitive deduplication:

```go
result, err := textnorm.New().FoldCase().SplitTokens().DedupTokens().JoinTokens(" ").Run("Red red RED")
// result == "red"
```

### `RemoveStopwords`

```go
func (tp TokenPipeline) RemoveStopwords(set map[string]struct{}) TokenPipeline
```

Returns a new token pipeline that drops tokens present in `set`. A `nil` set is a no-op (the pipeline is returned unchanged), which lets callers wire a stopword set from configuration without branching. Comparison is plain string equality (case-sensitive). Pair with `FoldCase()` upstream and a lowercase set for case-insensitive filtering.

**Example:**

```go
stop := map[string]struct{}{"the": {}, "a": {}, "of": {}}
result, err := textnorm.New().SplitTokens().RemoveStopwords(stop).JoinTokens(" ").
    Run("the cat sat on the mat of doom")
// result == "cat sat on mat doom"
```

### `JoinTokens`

```go
func (tp TokenPipeline) JoinTokens(sep string) Pipeline
```

Joins token output back into a string pipeline. The token pipeline is executed, and tokens are joined with `sep` using `strings.Join`.

**Example:**

```go
result, err := textnorm.New().SplitTokens().MapTokens(strings.ToUpper).JoinTokens("|").Run("go skip text")
// result == "GO|SKIP|TEXT"
```

---

## Package `textnorm/stopwords`

```go
import "github.com/alessiosavi/GoGPUtils/textnorm/stopwords"
```

The `stopwords` subpackage provides language-specific stopword sets sourced from the [NLTK stopwords corpus](https://github.com/nltk/nltk_data). Built-in sets are exposed as accessor functions that lazily parse embedded text files on first call and cache the result.

### Functions

#### `English`

```go
func English() map[string]struct{}
```

Returns the embedded English stopword set. The first call parses the embedded file; subsequent calls return the same cached map. The returned map is shared and **must not be mutated**.

#### `French`

```go
func French() map[string]struct{}
```

Returns the embedded French stopword set. Same caching and immutability semantics as `English()`.

#### `Italian`

```go
func Italian() map[string]struct{}
```

Returns the embedded Italian stopword set. Same caching and immutability semantics as `English()`.

#### `LoadFromFile`

```go
func LoadFromFile(lang string, filepath string) (map[string]struct{}, error)
```

Reads a newline-delimited word list from `filepath` and returns an independent stopword set. Empty lines and lines beginning with `#` are skipped. The `lang` argument is informational only.

#### `LoadFromList`

```go
func LoadFromList(lang string, words []string) map[string]struct{}
```

Builds a stopword set from a slice of words. Empty strings (after trimming) are skipped. The `lang` argument is informational only.

#### `CleanAllStopwords`

```go
func CleanAllStopwords(languages []string) map[string]struct{}
```

Returns the union of the built-in stopword sets for the requested languages. Recognized values (case-insensitive, trimmed): `"english"`, `"french"`, `"italian"`. Unknown languages are silently ignored. With zero recognized inputs, returns a non-nil empty map.

#### `Union`

```go
func Union(sets ...map[string]struct{}) map[string]struct{}
```

Returns a new set containing every key found in any of the input sets. Does not mutate any input. With zero inputs, returns a non-nil empty map.

### Important Notes

- The returned maps from `English()`, `French()`, and `Italian()` are **shared across callers** and must not be mutated. Use `Union()` to combine sets safely.
- Unicode normalization precondition: the non-ASCII entries in French (e.g., `"à"`) and Italian (e.g., `"è"`, `"più"`) are stored in composed (NFC) form. Tokens fed to `RemoveStopwords` must be in the same normalization form. The textnorm presets already include `NormalizeUnicode` upstream.

---

## Usage Examples

### Basic Preset Usage

```go
package main

import (
    "fmt"
    "github.com/alessiosavi/GoGPUtils/textnorm"
)

func main() {
    // Search-optimized normalization
    search, _ := textnorm.SearchPreset().Run("  Café, go!  ")
    fmt.Println(search) // "cafe go"

    // Canonical normalization
    canonical, _ := textnorm.CanonicalPreset().Run("  Hello, World!  ")
    fmt.Println(canonical) // "hello, world!"

    // Database-safe normalization
    dbsafe, _ := textnorm.DBSafePreset().Run("  Go\x00  ")
    fmt.Println(dbsafe) // "Go"
}
```

### Custom Pipeline Composition

```go
package main

import (
    "fmt"
    "unicode"

    "github.com/alessiosavi/GoGPUtils/textnorm"
    "golang.org/x/text/runes"
)

func main() {
    // Build a custom pipeline: sanitize, normalize, lowercase, keep only alphanumeric
    keep := runes.Predicate(func(r rune) bool {
        return unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r)
    })

    pipe := textnorm.New().
        SanitizeUTF8().
        NormalizeUnicode().
        Lower().
        FilterRunes(keep).
        TrimSpace().
        CollapseWhitespace()

    result, _ := pipe.Run("  Hello, 世界!  ")
    fmt.Println(result) // "hello 世界"
}
```

### Token Pipeline with Stopwords

```go
package main

import (
    "fmt"
    "github.com/alessiosavi/GoGPUtils/textnorm"
    "github.com/alessiosavi/GoGPUtils/textnorm/stopwords"
)

func main() {
    // Remove English stopwords and deduplicate
    result, _ := textnorm.SearchPreset().
        SplitTokens().
        DedupTokens().
        RemoveStopwords(stopwords.English()).
        JoinTokens(" ").
        Run("The the quick brown fox jumps over the lazy dog")

    fmt.Println(result) // "quick brown fox jumps over lazy dog"
}
```

### Width Folding for East Asian Text

```go
package main

import (
    "fmt"
    "github.com/alessiosavi/GoGPUtils/textnorm"
)

func main() {
    // Without width folding
    defaultOut, _ := textnorm.SearchPreset().Run("Ｇｏ")
    fmt.Println(defaultOut) // "Ｇｏ" (unchanged)

    // With width folding
    foldedOut, _ := textnorm.SearchPreset(textnorm.WithWidthFold()).Run("Ｇｏ")
    fmt.Println(foldedOut) // "go"
}
```

---

## Design Notes

- **Immutability**: Pipelines are immutable — `Then()` always returns a new pipeline. This makes them safe to reuse and share.
- **Nil safety**: All stage methods and token pipeline methods handle `nil` functions/sets gracefully by returning the original pipeline unchanged.
- **Error propagation**: `Run()` stops at the first error and returns it. Individual stages in this package never return errors, but custom stages may.
- **Zero dependencies for core**: The `textnorm` package depends only on `golang.org/x/text` for Unicode-aware operations. The stopwords subpackage has zero external dependencies.
- **Streaming deferred**: Streaming adapters are intentionally deferred until real usage proves they are worth the extra surface area.
