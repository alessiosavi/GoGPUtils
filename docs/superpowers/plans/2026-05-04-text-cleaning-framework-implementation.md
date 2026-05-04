# Text Cleaning — GoGPUtils + framework-golang Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add token-level dedup/stopword primitives + a stopwords subpackage to GoGPUtils, then expose a `pkg/text/` package on framework-golang (plus an opt-in `GetQueryParamsSearchKey` and CLAUDE.md guidance) that wraps GoGPUtils so AWS Lambda services have one stable text-cleaning entry point.

**Architecture:** GoGPUtils is the single source of truth for cleaning logic — additive changes to `textnorm/` and a new `textnorm/stopwords/` subpackage. framework-golang depends on GoGPUtils and thinly re-exports the presets under framework-versioned names. No existing code is modified beyond godoc updates and one additive sibling function (`GetQueryParamsSearchKey`).

**Tech Stack:** Go 1.25, `golang.org/x/text/transform`, standard `testing`. Existing patterns: table-driven tests with `t.Run(name, ...)`, fluent pipelines, returning `(value, error)`.

**Spec:** `docs/superpowers/specs/2026-05-04-text-cleaning-framework-design.md`

---

## File Map

**GoGPUtils (`/opt/SP/Workspace/Go/GoGPUtils/`):**

| File | Action | Responsibility |
|---|---|---|
| `textnorm/tokens.go` | Modify | Add `DedupTokens`, `RemoveStopwords` token stages. |
| `textnorm/tokens_test.go` | Modify | New table-driven tests for the two stages. |
| `textnorm/stopwords/stopwords.go` | Create | `English`, `French`, `Italian` sets + `Union` helper. |
| `textnorm/stopwords/stopwords_test.go` | Create | Sanity tests for the sets and `Union`. |
| `textnorm/presets.go` | Modify | Godoc-only: point at `DedupTokens`/`RemoveStopwords`. |

**framework-golang (`/opt/SP/Workspace/Go/framework-golang/`):**

| File | Action | Responsibility |
|---|---|---|
| `go.mod` / `go.sum` | Modify | `require github.com/alessiosavi/GoGPUtils …` |
| `pkg/text/doc.go` | Create | Package overview + decision table. |
| `pkg/text/text.go` | Create | `Clean`, `SearchKey`, `DBSafe`, `Slug`, `Pipeline`, `Stopwords`. |
| `pkg/text/text_test.go` | Create | Wrapper smoke tests (one per public symbol). |
| `pkg/request/params.go` | Modify | Add `GetQueryParamsSearchKey`. |
| `pkg/request/params_test.go` | Modify | Add `TestGetQueryParamsSearchKey`. |
| `Usage.md` | Modify | Append "Package `text`" section. |
| `CLAUDE.md` | Create | Top-level pointer to `pkg/text` for Claude. |

---

## Phase 1 — GoGPUtils

All commands in this phase run from `/opt/SP/Workspace/Go/GoGPUtils`.

### Task 1: `textnorm.DedupTokens` — failing tests

**Files:**
- Modify: `textnorm/tokens_test.go`

- [ ] **Step 1.1: Append failing tests for `DedupTokens`**

Append this block to the bottom of `textnorm/tokens_test.go`:

```go
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
```

- [ ] **Step 1.2: Confirm tests fail**

Run: `go test ./textnorm/ -run TestDedupTokens -v`
Expected: build error — `tp.DedupTokens undefined`.

### Task 2: `textnorm.DedupTokens` — implementation

**Files:**
- Modify: `textnorm/tokens.go`

- [ ] **Step 2.1: Append `DedupTokens` to `textnorm/tokens.go`**

Append at the bottom of `textnorm/tokens.go`:

```go
// DedupTokens returns a new token pipeline that drops duplicate tokens,
// preserving the first occurrence. Comparison is plain string equality
// (case-sensitive). Pair with FoldCase upstream for case-insensitive dedup.
func (tp TokenPipeline) DedupTokens() TokenPipeline {
	return tp.Then(func(tokens []string) ([]string, error) {
		seen := make(map[string]struct{}, len(tokens))
		out := make([]string, 0, len(tokens))
		for _, token := range tokens {
			if _, ok := seen[token]; ok {
				continue
			}
			seen[token] = struct{}{}
			out = append(out, token)
		}
		return out, nil
	})
}
```

- [ ] **Step 2.2: Confirm tests pass**

Run: `go test ./textnorm/ -run TestDedupTokens -v`
Expected: PASS for all four `TestDedupTokens*` tests.

- [ ] **Step 2.3: Commit**

```bash
git add textnorm/tokens.go textnorm/tokens_test.go
git commit -m "feat(textnorm): add DedupTokens token stage"
```

### Task 3: `textnorm.RemoveStopwords` — failing tests

**Files:**
- Modify: `textnorm/tokens_test.go`

- [ ] **Step 3.1: Append failing tests**

Append to `textnorm/tokens_test.go`:

```go
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
```

- [ ] **Step 3.2: Confirm tests fail**

Run: `go test ./textnorm/ -run TestRemoveStopwords -v`
Expected: build error — `tp.RemoveStopwords undefined`.

### Task 4: `textnorm.RemoveStopwords` — implementation

**Files:**
- Modify: `textnorm/tokens.go`

- [ ] **Step 4.1: Append `RemoveStopwords` to `textnorm/tokens.go`**

Append at the bottom of `textnorm/tokens.go`:

```go
// RemoveStopwords returns a new token pipeline that drops tokens present
// in set. A nil set is a no-op (the pipeline is returned unchanged), which
// lets callers wire a stopword set from configuration without branching.
// Comparison is plain string equality (case-sensitive). Pair with FoldCase
// upstream and a lowercase set for case-insensitive filtering.
func (tp TokenPipeline) RemoveStopwords(set map[string]struct{}) TokenPipeline {
	if set == nil {
		return tp
	}
	return tp.Then(func(tokens []string) ([]string, error) {
		out := make([]string, 0, len(tokens))
		for _, token := range tokens {
			if _, drop := set[token]; drop {
				continue
			}
			out = append(out, token)
		}
		return out, nil
	})
}
```

- [ ] **Step 4.2: Confirm tests pass**

Run: `go test ./textnorm/ -run TestRemoveStopwords -v`
Expected: PASS for all four `TestRemoveStopwords*` tests.

- [ ] **Step 4.3: Run full textnorm package to confirm no regressions**

Run: `go test ./textnorm/...`
Expected: all tests pass.

- [ ] **Step 4.4: Commit**

```bash
git add textnorm/tokens.go textnorm/tokens_test.go
git commit -m "feat(textnorm): add RemoveStopwords token stage"
```

### Task 5: `textnorm/stopwords` — failing tests

**Files:**
- Create: `textnorm/stopwords/stopwords_test.go`

- [ ] **Step 5.1: Create `textnorm/stopwords/stopwords_test.go`**

Write this entire file:

```go
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
	_ = Union(a, b)
	if len(a) != 1 || len(b) != 1 {
		t.Fatalf("Union mutated inputs: a=%v b=%v", a, b)
	}
}
```

- [ ] **Step 5.2: Confirm tests fail**

Run: `go test ./textnorm/stopwords/...`
Expected: build error — package directory has no Go file.

### Task 6: `textnorm/stopwords` — implementation

**Files:**
- Create: `textnorm/stopwords/stopwords.go`

- [ ] **Step 6.1: Create `textnorm/stopwords/stopwords.go`**

Write this entire file:

```go
// Package stopwords provides language-specific stopword sets for use with
// textnorm.RemoveStopwords. The sets are exported as map[string]struct{} so
// callers get O(1) membership checks and can pass them straight to the
// RemoveStopwords stage.
//
// The variables are documented as read-only — Go has no const map, but
// callers MUST NOT mutate them. Use Union to combine sets safely.
package stopwords

// English is a default English stopword set covering ~150 common function
// words. Suitable for search-key filtering, embeddings hygiene, and slug
// generation.
var English = map[string]struct{}{
	"a": {}, "about": {}, "above": {}, "after": {}, "again": {}, "against": {},
	"all": {}, "am": {}, "an": {}, "and": {}, "any": {}, "are": {}, "as": {},
	"at": {}, "be": {}, "because": {}, "been": {}, "before": {}, "being": {},
	"below": {}, "between": {}, "both": {}, "but": {}, "by": {}, "can": {},
	"could": {}, "did": {}, "do": {}, "does": {}, "doing": {}, "don": {},
	"down": {}, "during": {}, "each": {}, "few": {}, "for": {}, "from": {},
	"further": {}, "had": {}, "has": {}, "have": {}, "having": {}, "he": {},
	"her": {}, "here": {}, "hers": {}, "herself": {}, "him": {}, "himself": {},
	"his": {}, "how": {}, "i": {}, "if": {}, "in": {}, "into": {}, "is": {},
	"it": {}, "its": {}, "itself": {}, "just": {}, "me": {}, "more": {},
	"most": {}, "my": {}, "myself": {}, "no": {}, "nor": {}, "not": {},
	"now": {}, "of": {}, "off": {}, "on": {}, "once": {}, "only": {}, "or": {},
	"other": {}, "our": {}, "ours": {}, "ourselves": {}, "out": {}, "over": {},
	"own": {}, "same": {}, "she": {}, "should": {}, "so": {}, "some": {},
	"such": {}, "than": {}, "that": {}, "the": {}, "their": {}, "theirs": {},
	"them": {}, "themselves": {}, "then": {}, "there": {}, "these": {},
	"they": {}, "this": {}, "those": {}, "through": {}, "to": {}, "too": {},
	"under": {}, "until": {}, "up": {}, "very": {}, "was": {}, "we": {},
	"were": {}, "what": {}, "when": {}, "where": {}, "which": {}, "while": {},
	"who": {}, "whom": {}, "why": {}, "will": {}, "with": {}, "would": {},
	"you": {}, "your": {}, "yours": {}, "yourself": {}, "yourselves": {},
}

// French is a starter French stopword set. Expected to grow; treat the API
// (variable name + type) as the contract, the contents as data.
var French = map[string]struct{}{
	"à": {}, "au": {}, "aux": {}, "avec": {}, "ce": {}, "ces": {}, "cette": {},
	"d": {}, "dans": {}, "de": {}, "des": {}, "du": {}, "elle": {}, "en": {},
	"est": {}, "et": {}, "eu": {}, "il": {}, "ils": {}, "j": {}, "je": {},
	"l": {}, "la": {}, "le": {}, "les": {}, "leur": {}, "lui": {}, "ma": {},
	"mais": {}, "me": {}, "mes": {}, "moi": {}, "mon": {}, "n": {}, "ne": {},
	"nos": {}, "notre": {}, "nous": {}, "on": {}, "ou": {}, "par": {},
	"pas": {}, "pour": {}, "qu": {}, "que": {}, "qui": {}, "s": {}, "sa": {},
	"se": {}, "ses": {}, "son": {}, "sur": {}, "ta": {}, "te": {}, "tes": {},
	"toi": {}, "ton": {}, "tu": {}, "un": {}, "une": {}, "vos": {},
	"votre": {}, "vous": {}, "y": {},
}

// Italian is a starter Italian stopword set. Expected to grow; treat the API
// (variable name + type) as the contract, the contents as data.
var Italian = map[string]struct{}{
	"a": {}, "ad": {}, "al": {}, "alla": {}, "alle": {}, "agli": {}, "ai": {},
	"allo": {}, "anche": {}, "che": {}, "chi": {}, "ci": {}, "come": {},
	"con": {}, "da": {}, "dal": {}, "dalla": {}, "del": {}, "della": {},
	"delle": {}, "dello": {}, "degli": {}, "dei": {}, "di": {}, "e": {},
	"ed": {}, "è": {}, "gli": {}, "i": {}, "il": {}, "in": {}, "io": {},
	"l": {}, "la": {}, "le": {}, "lei": {}, "li": {}, "lo": {}, "loro": {},
	"lui": {}, "ma": {}, "me": {}, "mi": {}, "mio": {}, "ne": {}, "nei": {},
	"nel": {}, "nella": {}, "nelle": {}, "no": {}, "noi": {}, "non": {},
	"o": {}, "per": {}, "più": {}, "questa": {}, "queste": {}, "questi": {},
	"questo": {}, "se": {}, "sei": {}, "si": {}, "sia": {}, "sono": {},
	"su": {}, "sul": {}, "sulla": {}, "te": {}, "ti": {}, "tra": {}, "tu": {},
	"tuo": {}, "un": {}, "una": {}, "uno": {}, "voi": {},
}

// Union returns a new set containing every key found in any of the input
// sets. It does not mutate any input. With zero inputs it returns a non-nil
// empty map.
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
```

- [ ] **Step 6.2: Confirm tests pass**

Run: `go test ./textnorm/stopwords/... -v`
Expected: all six tests pass.

- [ ] **Step 6.3: Commit**

```bash
git add textnorm/stopwords/stopwords.go textnorm/stopwords/stopwords_test.go
git commit -m "feat(textnorm/stopwords): add English/French/Italian sets and Union helper"
```

### Task 7: `textnorm` presets godoc

**Files:**
- Modify: `textnorm/presets.go`

- [ ] **Step 7.1: Update godoc on `SearchPreset`**

In `textnorm/presets.go`, replace this comment block:

```go
// SearchPreset builds a search-key pipeline.
func SearchPreset(opts ...PresetOption) Pipeline {
```

with:

```go
// SearchPreset builds a search-key pipeline (UTF-8 sanitize, Unicode
// normalize, fold case, keep letters/numbers/spaces, trim, collapse).
//
// To dedup repeated tokens or strip stopwords, compose on top:
//
//	textnorm.SearchPreset().
//	    SplitTokens().
//	    DedupTokens().
//	    RemoveStopwords(stopwords.English).
//	    JoinTokens(" ")
func SearchPreset(opts ...PresetOption) Pipeline {
```

- [ ] **Step 7.2: Confirm package still builds and tests pass**

Run: `go test ./textnorm/...`
Expected: all tests pass.

- [ ] **Step 7.3: Commit**

```bash
git add textnorm/presets.go
git commit -m "docs(textnorm): document DedupTokens/RemoveStopwords composition on SearchPreset"
```

### Task 8: GoGPUtils integration check

- [ ] **Step 8.1: Run full GoGPUtils test suite**

Run: `go test ./...`
Expected: all tests pass (existing + new).

- [ ] **Step 8.2: Run go vet and go fmt**

Run:
```bash
go vet ./...
gofmt -l textnorm/
```
Expected: empty output from both.

- [ ] **Step 8.3: Tidy modules**

Run: `go mod tidy`
Expected: no changes (we added no new imports).

- [ ] **Step 8.4: Note for human handoff**

After Phase 1 lands on a release tag, capture the tag (e.g. `v0.0.106` or whatever the maintainer chooses) and pass it to Phase 2, Step 9.1. Tagging is a manual maintainer action and is not part of this plan.

---

## Phase 2 — framework-golang

All commands in this phase run from `/opt/SP/Workspace/Go/framework-golang`.

### Task 9: Add GoGPUtils dependency

**Files:**
- Modify: `go.mod`, `go.sum`

- [ ] **Step 9.1: Add the GoGPUtils dependency**

Run, substituting the tag captured at the end of Phase 1 (e.g. `v0.0.106`):

```bash
cd /opt/SP/Workspace/Go/framework-golang
go get github.com/alessiosavi/GoGPUtils@<tag>
go mod tidy
```

Expected: `go.mod` gains `github.com/alessiosavi/GoGPUtils <tag>` in the `require` block; `go.sum` gains the matching hashes.

- [ ] **Step 9.2: Verify build still works**

Run: `go build ./...`
Expected: clean build.

- [ ] **Step 9.3: Commit**

```bash
git add go.mod go.sum
git commit -m "chore(deps): add github.com/alessiosavi/GoGPUtils for pkg/text"
```

### Task 10: `pkg/text/text.go` — `Clean` (TDD)

**Files:**
- Create: `pkg/text/text_test.go`
- Create: `pkg/text/text.go`

- [ ] **Step 10.1: Write the failing test for `Clean`**

Create `pkg/text/text_test.go`:

```go
package text

import "testing"

func TestCleanLowercaseTrimsAndCollapses(t *testing.T) {
	got, err := Clean("  Héllo   World  ")
	if err != nil {
		t.Fatalf("Clean() error = %v", err)
	}
	if got != "hello world" {
		t.Fatalf("Clean() = %q, want %q", got, "hello world")
	}
}

func TestCleanEmptyInput(t *testing.T) {
	got, err := Clean("")
	if err != nil {
		t.Fatalf("Clean() error = %v", err)
	}
	if got != "" {
		t.Fatalf("Clean() = %q, want %q", got, "")
	}
}
```

- [ ] **Step 10.2: Confirm test fails**

Run: `go test ./pkg/text/...`
Expected: build error — package `text` has no Go file.

- [ ] **Step 10.3: Create `pkg/text/text.go` with `Clean`**

Write this entire file:

```go
package text

import "github.com/alessiosavi/GoGPUtils/textnorm"

// Clean applies a single fixed-order canonical clean: UTF-8 sanitize →
// Unicode normalize (diacritics stripped) → fold case → trim → collapse
// whitespace. Use for free-text fields stored or compared as-is.
// Wraps textnorm.CanonicalPreset.
func Clean(s string) (string, error) {
	return textnorm.CanonicalPreset().Run(s)
}
```

- [ ] **Step 10.4: Confirm test passes**

Run: `go test ./pkg/text/... -v`
Expected: PASS for `TestCleanLowercaseTrimsAndCollapses` and `TestCleanEmptyInput`.

- [ ] **Step 10.5: Commit**

```bash
git add pkg/text/text.go pkg/text/text_test.go
git commit -m "feat(text): add Clean wrapper over textnorm.CanonicalPreset"
```

### Task 11: `SearchKey`

**Files:**
- Modify: `pkg/text/text_test.go`, `pkg/text/text.go`

- [ ] **Step 11.1: Append failing test**

Append to `pkg/text/text_test.go`:

```go
func TestSearchKeyLowercasesStripsAccentsAndPunctuation(t *testing.T) {
	got, err := SearchKey(" Café, Hermès! ")
	if err != nil {
		t.Fatalf("SearchKey() error = %v", err)
	}
	if got != "cafe hermes" {
		t.Fatalf("SearchKey() = %q, want %q", got, "cafe hermes")
	}
}
```

- [ ] **Step 11.2: Confirm fails**

Run: `go test ./pkg/text/ -run TestSearchKey -v`
Expected: build error — `SearchKey undefined`.

- [ ] **Step 11.3: Append `SearchKey` to `pkg/text/text.go`**

Append at the bottom of `pkg/text/text.go`:

```go
// SearchKey produces a search-friendly form: Clean + diacritic strip +
// non-letter/digit/space removal + token split + rejoin. Use for search
// indexing, query matching, deterministic comparison.
// Wraps textnorm.SearchPreset.
func SearchKey(s string) (string, error) {
	return textnorm.SearchPreset().Run(s)
}
```

- [ ] **Step 11.4: Confirm passes**

Run: `go test ./pkg/text/ -run TestSearchKey -v`
Expected: PASS.

- [ ] **Step 11.5: Commit**

```bash
git add pkg/text/text.go pkg/text/text_test.go
git commit -m "feat(text): add SearchKey wrapper over textnorm.SearchPreset"
```

### Task 12: `DBSafe`

**Files:**
- Modify: `pkg/text/text_test.go`, `pkg/text/text.go`

- [ ] **Step 12.1: Append failing test**

Append to `pkg/text/text_test.go`:

```go
func TestDBSafeStripsDiacriticsAndCollapsesButKeepsCase(t *testing.T) {
	got, err := DBSafe("  Héllo   Wörld  ")
	if err != nil {
		t.Fatalf("DBSafe() error = %v", err)
	}
	if got != "Hello World" {
		t.Fatalf("DBSafe() = %q, want %q", got, "Hello World")
	}
}
```

- [ ] **Step 12.2: Confirm fails**

Run: `go test ./pkg/text/ -run TestDBSafe -v`
Expected: build error — `DBSafe undefined`.

- [ ] **Step 12.3: Append `DBSafe` to `pkg/text/text.go`**

Append at the bottom of `pkg/text/text.go`:

```go
// DBSafe produces persistence-safe text: UTF-8 sanitize, Unicode normalize
// (NFD + drop combining marks + NFC — diacritics are stripped), trim,
// collapse whitespace. Does NOT lowercase. Use for storing scraper or
// normalizer output where you want stable byte-for-byte storage but want
// to keep original case. Wraps textnorm.DBSafePreset.
func DBSafe(s string) (string, error) {
	return textnorm.DBSafePreset().Run(s)
}
```

- [ ] **Step 12.4: Confirm passes**

Run: `go test ./pkg/text/ -run TestDBSafe -v`
Expected: PASS.

- [ ] **Step 12.5: Commit**

```bash
git add pkg/text/text.go pkg/text/text_test.go
git commit -m "feat(text): add DBSafe wrapper over textnorm.DBSafePreset"
```

### Task 13: `Slug`

**Files:**
- Modify: `pkg/text/text_test.go`, `pkg/text/text.go`

- [ ] **Step 13.1: Append failing test**

Append to `pkg/text/text_test.go`:

```go
func TestSlugProducesURLFriendlyForm(t *testing.T) {
	got, err := Slug("Hello, World! Café")
	if err != nil {
		t.Fatalf("Slug() error = %v", err)
	}
	if got != "hello-world-cafe" {
		t.Fatalf("Slug() = %q, want %q", got, "hello-world-cafe")
	}
}
```

- [ ] **Step 13.2: Confirm fails**

Run: `go test ./pkg/text/ -run TestSlug -v`
Expected: build error — `Slug undefined`.

- [ ] **Step 13.3: Add `Slug` import + function**

In `pkg/text/text.go`, change the import block to:

```go
import (
	"github.com/alessiosavi/GoGPUtils/stringutil"
	"github.com/alessiosavi/GoGPUtils/textnorm"
)
```

Then append at the bottom of `pkg/text/text.go`:

```go
// Slug produces a URL-friendly slug (lowercased, diacritics stripped,
// non-alphanumerics replaced with single hyphens, leading/trailing hyphens
// trimmed). Use for URLs, filenames, IDs.
// Wraps stringutil.Slugify.
func Slug(s string) (string, error) {
	return stringutil.Slugify(s)
}
```

- [ ] **Step 13.4: Confirm passes**

Run: `go test ./pkg/text/ -run TestSlug -v`
Expected: PASS.

- [ ] **Step 13.5: Commit**

```bash
git add pkg/text/text.go pkg/text/text_test.go
git commit -m "feat(text): add Slug wrapper over stringutil.Slugify"
```

### Task 14: `Pipeline` + `Stopwords`

**Files:**
- Modify: `pkg/text/text_test.go`, `pkg/text/text.go`

- [ ] **Step 14.1: Append failing tests**

Append to `pkg/text/text_test.go`:

```go
func TestPipelineComposesCustomFlow(t *testing.T) {
	got, err := Pipeline().
		SanitizeUTF8().
		NormalizeUnicode().
		FoldCase().
		TrimSpace().
		CollapseWhitespace().
		SplitTokens().
		DedupTokens().
		RemoveStopwords(Stopwords.English).
		JoinTokens(" ").
		Run("The quick brown fox jumps over the lazy fox")
	if err != nil {
		t.Fatalf("Pipeline() error = %v", err)
	}
	want := "quick brown fox jumps over lazy"
	if got != want {
		t.Fatalf("Pipeline() = %q, want %q", got, want)
	}
}

func TestStopwordsExportNonEmpty(t *testing.T) {
	if len(Stopwords.English) == 0 {
		t.Error("Stopwords.English is empty")
	}
	if len(Stopwords.French) == 0 {
		t.Error("Stopwords.French is empty")
	}
	if len(Stopwords.Italian) == 0 {
		t.Error("Stopwords.Italian is empty")
	}
}
```

- [ ] **Step 14.2: Confirm fails**

Run: `go test ./pkg/text/ -run "TestPipeline|TestStopwords" -v`
Expected: build error — `Pipeline undefined`, `Stopwords undefined`.

- [ ] **Step 14.3: Add `Pipeline` and `Stopwords` to `pkg/text/text.go`**

In `pkg/text/text.go`, change the import block to:

```go
import (
	"github.com/alessiosavi/GoGPUtils/stringutil"
	"github.com/alessiosavi/GoGPUtils/textnorm"
	"github.com/alessiosavi/GoGPUtils/textnorm/stopwords"
)
```

Append at the bottom of `pkg/text/text.go`:

```go
// Pipeline returns a fresh GoGPUtils textnorm.Pipeline so callers can
// compose custom flows (DedupTokens, RemoveStopwords, etc.) without
// importing GoGPUtils directly. Use this for embeddings hygiene,
// cross-field dedup composition, or any flow not covered by the four
// preset wrappers above.
func Pipeline() textnorm.Pipeline {
	return textnorm.New()
}

// Stopwords re-exports the GoGPUtils stopword sets so callers stay on the
// framework import path. The maps MUST NOT be mutated; use textnorm/
// stopwords.Union (or compose your own copy) if you need to combine sets.
var Stopwords = struct {
	English map[string]struct{}
	French  map[string]struct{}
	Italian map[string]struct{}
}{
	English: stopwords.English,
	French:  stopwords.French,
	Italian: stopwords.Italian,
}
```

- [ ] **Step 14.4: Confirm passes**

Run: `go test ./pkg/text/ -run "TestPipeline|TestStopwords" -v`
Expected: PASS for both tests.

- [ ] **Step 14.5: Run full pkg/text test suite**

Run: `go test ./pkg/text/...`
Expected: all tests pass.

- [ ] **Step 14.6: Commit**

```bash
git add pkg/text/text.go pkg/text/text_test.go
git commit -m "feat(text): add Pipeline accessor and Stopwords export"
```

### Task 15: `pkg/text/doc.go`

**Files:**
- Create: `pkg/text/doc.go`

- [ ] **Step 15.1: Create `pkg/text/doc.go`**

Write this entire file:

```go
// Package text is the canonical entry point for cleaning user-supplied or
// scraped strings in framework-golang services. It wraps the
// github.com/alessiosavi/GoGPUtils textnorm and stringutil packages so
// service code stays on a stable framework-versioned API.
//
// # Decision table
//
//	Input you're processing                    | Use
//	-------------------------------------------+----------------
//	Search query (matched against an index)    | SearchKey or
//	                                           |   request.GetQueryParamsSearchKey
//	Free text persisted in a DB column         | DBSafe
//	URL, filename, slug                        | Slug
//	General canonicalization                   | Clean
//	Embeddings input, dedup, stopword removal  | Pipeline + token stages
//
// # When NOT to clean
//
//   - JWT claims, IDs, UUIDs, enums, AWS resource ARNs — already
//     constrained, cleaning corrupts them.
//   - Fields you'll re-emit verbatim to a user — preserve their input.
//   - Numbers, timestamps, currency codes — wrong tool.
//
// # Custom flows
//
// Use Pipeline to access the underlying textnorm.Pipeline for composing
// dedup, stopword removal, or any other token-level flow:
//
//	clean, _ := text.Pipeline().
//	    SanitizeUTF8().
//	    NormalizeUnicode().
//	    FoldCase().
//	    TrimSpace().
//	    CollapseWhitespace().
//	    SplitTokens().
//	    DedupTokens().
//	    RemoveStopwords(text.Stopwords.English).
//	    JoinTokens(" ").
//	    Run(rawText)
package text
```

- [ ] **Step 15.2: Verify build**

Run: `go build ./pkg/text/`
Expected: clean build.

- [ ] **Step 15.3: Verify go doc renders**

Run: `go doc ./pkg/text`
Expected: package overview prints with the decision table and "When NOT to clean" sections.

- [ ] **Step 15.4: Commit**

```bash
git add pkg/text/doc.go
git commit -m "docs(text): add package overview with decision table"
```

### Task 16: `request.GetQueryParamsSearchKey` (TDD)

**Files:**
- Modify: `pkg/request/params_test.go`, `pkg/request/params.go`

- [ ] **Step 16.1: Inspect existing test layout for fixture conventions**

Run: `head -50 pkg/request/params_test.go`
Expected output gives you the existing `TestGetQueryParams` table-driven test you'll mirror.

- [ ] **Step 16.2: Append failing test to `pkg/request/params_test.go`**

Append:

```go
func TestGetQueryParamsSearchKey(t *testing.T) {
	tests := []struct {
		name     string
		params   map[string]string
		field    string
		want     string
		wantErr  bool
	}{
		{
			name:    "lowercases-strips-accents-and-punctuation",
			params:  map[string]string{"q": " Café, Hermès! "},
			field:   "q",
			want:    "cafe hermes",
			wantErr: false,
		},
		{
			name:    "missing-param-returns-error",
			params:  map[string]string{"other": "x"},
			field:   "q",
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty-param-returns-error",
			params:  map[string]string{"q": ""},
			field:   "q",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := events.APIGatewayProxyRequest{QueryStringParameters: tt.params}
			got, err := GetQueryParamsSearchKey(req, tt.field)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("got = %q, want %q", got, tt.want)
			}
		})
	}
}
```

If `events` is not already imported in this test file, add the import: `"github.com/aws/aws-lambda-go/events"`.

- [ ] **Step 16.3: Confirm fails**

Run: `go test ./pkg/request/ -run TestGetQueryParamsSearchKey -v`
Expected: build error — `GetQueryParamsSearchKey undefined`.

- [ ] **Step 16.4: Add `GetQueryParamsSearchKey` to `pkg/request/params.go`**

In `pkg/request/params.go`, change the import block to:

```go
import (
	"errors"

	"github.com/GreenCapitals/framework-golang/pkg/text"
	"github.com/aws/aws-lambda-go/events"
	"github.com/microcosm-cc/bluemonday"
)
```

Then append at the bottom of `pkg/request/params.go`:

```go
// GetQueryParamsSearchKey returns the named query parameter passed through
// the framework's search-key normalization (text.SearchKey). Use for
// endpoints that match against search-indexed columns. Existing callers
// of GetQueryParams are not affected — that function still returns the
// bluemonday-sanitized value.
func GetQueryParamsSearchKey(r events.APIGatewayProxyRequest, paramName string) (string, error) {
	paramValue, exists := r.QueryStringParameters[paramName]
	if !exists || paramValue == "" {
		return "", errors.New("query parameter is required")
	}
	return text.SearchKey(paramValue)
}
```

- [ ] **Step 16.5: Confirm passes**

Run: `go test ./pkg/request/...`
Expected: all tests pass — both the existing `TestGetQueryParams` and the new `TestGetQueryParamsSearchKey`.

- [ ] **Step 16.6: Commit**

```bash
git add pkg/request/params.go pkg/request/params_test.go
git commit -m "feat(request): add opt-in GetQueryParamsSearchKey wrapper"
```

### Task 17: `Usage.md` documentation

**Files:**
- Modify: `Usage.md`

- [ ] **Step 17.1: Append the `Package text` section to `Usage.md`**

Append at the bottom of `Usage.md` (after the existing "Handler Type Reference" table):

````markdown

---

## Package `text` - String Cleaning & Normalization

```go
import "github.com/GreenCapitals/framework-golang/pkg/text"
```

This package is the canonical entry point for cleaning user-supplied or
scraped text. **Do not invent ad-hoc cleaning** (no inline `strings.ToLower`
+ `strings.TrimSpace` + accent loops). Reach for one of these wrappers:

### When to use which

| You're processing…                        | Call                                   |
|-------------------------------------------|----------------------------------------|
| A search query (matching against an index) | `text.SearchKey` or `request.GetQueryParamsSearchKey` |
| Text you'll store in a user-visible column | `text.DBSafe`                          |
| A field for slugs, URLs, filenames         | `text.Slug`                            |
| Anything else needing a canonical form     | `text.Clean`                           |
| Embeddings input, dedup, stopword removal  | `text.Pipeline()` + token stages       |

### When NOT to clean

- **JWT claims, IDs, UUIDs, enums, AWS resource ARNs** — already constrained,
  cleaning corrupts them.
- **Fields you'll re-emit verbatim to a user** — preserve their input.
- **Numbers, timestamps, currency codes** — wrong tool.

### Examples

```go
// Search query from API Gateway
key, err := request.GetQueryParamsSearchKey(r, "q")
// "Café Hermès " → "cafe hermes"

// Persisting a free-text description
clean, _ := text.DBSafe(input.Description)

// Embeddings hygiene: strip duplicates and stopwords before sending to Titan
clean, _ := text.Pipeline().
    SanitizeUTF8().
    NormalizeUnicode().
    FoldCase().
    TrimSpace().
    CollapseWhitespace().
    SplitTokens().
    DedupTokens().
    RemoveStopwords(text.Stopwords.English).
    JoinTokens(" ").
    Run(rawText)
```

### Stopwords

`text.Stopwords.English` is populated. `text.Stopwords.French` and
`text.Stopwords.Italian` are starter sets — expected to grow. Combine with
the Union helper from the underlying package when needed:

```go
import "github.com/alessiosavi/GoGPUtils/textnorm/stopwords"

multi := stopwords.Union(text.Stopwords.English, text.Stopwords.French)
```
````

- [ ] **Step 17.2: Verify markdown renders without surprises**

Visual inspection — open `Usage.md` and confirm the new section sits at the bottom with proper headings.

- [ ] **Step 17.3: Commit**

```bash
git add Usage.md
git commit -m "docs: document pkg/text in Usage.md"
```

### Task 18: Top-level `CLAUDE.md`

**Files:**
- Create: `CLAUDE.md`

- [ ] **Step 18.1: Verify CLAUDE.md does not already exist**

Run: `ls CLAUDE.md 2>/dev/null && echo EXISTS || echo MISSING`
Expected: `MISSING`. (If `EXISTS`, append the new text-cleaning section instead of overwriting.)

- [ ] **Step 18.2: Create `CLAUDE.md`**

Write this entire file:

```markdown
# Framework-Golang — Claude Code Notes

Primary reference: `Usage.md`.

## Text cleaning

When a handler accepts user-supplied or scraped strings, prefer
`pkg/text` over inline cleaning. Decision table is in `Usage.md`
under "Package `text`". Do not import `github.com/alessiosavi/GoGPUtils`
directly from service code — go through `pkg/text` so framework upgrades
stay mechanical.

Quick reference:

- `text.SearchKey` — search queries / index matching.
- `text.DBSafe` — free text persisted to a column.
- `text.Slug` — URLs, filenames, IDs.
- `text.Clean` — generic canonicalization.
- `text.Pipeline()` — custom token-level flows (dedup, stopwords).

When **not** to clean: JWT claims, UUIDs, ARNs, enums, numbers, anything
the user re-reads verbatim.
```

- [ ] **Step 18.3: Commit**

```bash
git add CLAUDE.md
git commit -m "docs: add top-level CLAUDE.md pointing at pkg/text"
```

### Task 19: Final framework-golang integration check

- [ ] **Step 19.1: Run full framework-golang test suite**

Run: `go test ./...`
Expected: all tests pass.

- [ ] **Step 19.2: Vet, fmt, tidy**

Run:
```bash
go vet ./...
gofmt -l pkg/text/ pkg/request/
go mod tidy
```
Expected: vet exits 0; gofmt prints nothing; mod tidy makes no changes.

- [ ] **Step 19.3: Note for human handoff**

Phase 2 is now ready for review and a release tag. Tagging is a manual maintainer action.

---

## Verification Summary

After both phases land, the system has:

1. `github.com/alessiosavi/GoGPUtils/textnorm` — `DedupTokens()` + `RemoveStopwords(set)` token stages, fully tested.
2. `github.com/alessiosavi/GoGPUtils/textnorm/stopwords` — `English` (full), `French`/`Italian` (starter), `Union(...)`.
3. `github.com/GreenCapitals/framework-golang/pkg/text` — `Clean`, `SearchKey`, `DBSafe`, `Slug`, `Pipeline()`, `Stopwords` exposed; doc.go decision table available via `go doc`.
4. `github.com/GreenCapitals/framework-golang/pkg/request` — `GetQueryParamsSearchKey` available alongside the unchanged `GetQueryParams`.
5. Framework consumers have a top-level `CLAUDE.md` and `Usage.md` section telling Claude where to reach when handling raw input.
6. **No existing test changes its expected output. No service is migrated. veliu.com is untouched.**
