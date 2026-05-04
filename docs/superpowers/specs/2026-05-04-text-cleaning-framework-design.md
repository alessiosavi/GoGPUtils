# Text Cleaning — GoGPUtils → framework-golang Alignment

**Date:** 2026-05-04
**Status:** Approved (brainstorming complete)
**Repos touched:** `github.com/alessiosavi/GoGPUtils`, `github.com/GreenCapitals/framework-golang`
**Repos NOT touched (this work):** `veliu.com`

## Goal

Give `framework-golang` consumers a single, opinionated entry point for cleaning user-supplied or scraped text — backed by the existing `GoGPUtils` `textnorm` pipeline. Add the missing token-level primitives (`DedupTokens`, `RemoveStopwords`) so domain code in services like `lib/services/normalizer/` and `lib/services/embeddings/` can compose recipes (e.g. cross-field dedup, embeddings hygiene) on top of stable, tested primitives instead of reinventing them inline.

## Non-Goals

- Replace `models.NormalizeText()` in `veliu.com` as part of this work. The existing recipe keeps running. Migration is opportunistic — the next PR that touches normalizer text moves to `pkg/text`.
- Cross-field dedup (e.g. strip `brand` tokens from `title`). That's domain logic and stays in the consumer; the framework only ships primitives.
- Streaming / `io.Reader` adapters (deferred per `textnorm/doc.go`).
- Change the behavior of the existing `request.GetQueryParams`. It keeps the bluemonday `StrictPolicy` it has today.
- Remove anything from `GoGPUtils/stringutil/clean.go` (independent surface, other consumers may use it).

## Architecture

```
┌──────────────────────────────────────────────────────────────┐
│  GoGPUtils  (github.com/alessiosavi/GoGPUtils)               │
│                                                              │
│  textnorm/                                                   │
│   ├── pipeline.go         (existing)                         │
│   ├── stages.go           (existing)                         │
│   ├── tokens.go           ✚ DedupTokens, RemoveStopwords     │
│   ├── presets.go          (existing — unchanged)             │
│   └── stopwords/          ✚ NEW subpackage                   │
│        ├── stopwords.go     (English/French/Italian + Union) │
│        └── stopwords_test.go                                 │
└──────────────────────────────────────────────────────────────┘
                           ▲
                           │  go.mod require
                           │
┌──────────────────────────┴───────────────────────────────────┐
│  framework-golang  (github.com/GreenCapitals/framework-golang)│
│                                                              │
│  pkg/text/               ✚ NEW package                       │
│   ├── doc.go              (overview + decision tree)         │
│   ├── text.go             (Clean, SearchKey, DBSafe, Slug,   │
│   │                        Pipeline, Stopwords)              │
│   ├── presets.go          (preset wiring)                    │
│   └── text_test.go                                           │
│                                                              │
│  pkg/request/                                                │
│   ├── params.go           ✚ GetQueryParamsSearchKey          │
│   └── params_test.go      ✚ TestGetQueryParamsSearchKey      │
│                                                              │
│  CLAUDE.md                ✚ NEW file (top-level pointer)     │
│  Usage.md                 ✚ "Package `text`" section added   │
└──────────────────────────────────────────────────────────────┘
                           ▲
                           │  consumed by
                           │
┌──────────────────────────┴───────────────────────────────────┐
│  veliu.com  (NOT modified in this work)                      │
│   - lib/services/normalizer/                                 │
│   - lib/services/embeddings/                                 │
│   These will migrate piecewise on next text-touching PR,     │
│   guided by the new CLAUDE.md / Usage.md sections.           │
└──────────────────────────────────────────────────────────────┘
```

## GoGPUtils Changes

### 1. `textnorm/tokens.go` — two new TokenStages

**`DedupTokens()`** — drops repeated tokens, preserves first-seen order. Plain string equality (case-sensitive). Callers wanting case-insensitive dedup must call `FoldCase()` upstream in the pipeline.

```go
func (tp TokenPipeline) DedupTokens() TokenPipeline {
    return tp.Then(func(tokens []string) ([]string, error) {
        seen := make(map[string]struct{}, len(tokens))
        out := make([]string, 0, len(tokens))
        for _, t := range tokens {
            if _, ok := seen[t]; ok {
                continue
            }
            seen[t] = struct{}{}
            out = append(out, t)
        }
        return out, nil
    })
}
```

**`RemoveStopwords(set map[string]struct{})`** — drops tokens present in `set`. A `nil` set is a no-op (returns the stage unchanged so callers can pass results from a config knob without branching).

```go
func (tp TokenPipeline) RemoveStopwords(set map[string]struct{}) TokenPipeline {
    if set == nil {
        return tp
    }
    return tp.Then(func(tokens []string) ([]string, error) {
        out := make([]string, 0, len(tokens))
        for _, t := range tokens {
            if _, drop := set[t]; drop {
                continue
            }
            out = append(out, t)
        }
        return out, nil
    })
}
```

### 2. `textnorm/stopwords/` — new subpackage

```
textnorm/stopwords/
├── stopwords.go
└── stopwords_test.go
```

**Public surface:**

```go
package stopwords

// English is a default English stopword set (~150 common words).
var English = map[string]struct{}{ ... }

// French is a starter French stopword set. Expected to grow.
var French = map[string]struct{}{ ... }

// Italian is a starter Italian stopword set. Expected to grow.
var Italian = map[string]struct{}{ ... }

// Union returns a new set containing every key in any of the input sets.
// Convenience for multilingual filtering.
func Union(sets ...map[string]struct{}) map[string]struct{}
```

Sets are exported as `map[string]struct{}` rather than slices: O(1) membership check at call time, and the `RemoveStopwords` stage takes the same type. Callers MUST treat the variables as read-only (Go has no `const map`); the package documents this.

`English` ships fully populated. `French` and `Italian` ship with starter content (~30–50 entries) — the API is the contract; data grows over time without API churn.

### 3. `textnorm/presets.go` — godoc only

No code changes. Update godoc on `SearchPreset` to point at the new token-level primitives:

```go
// SearchPreset builds a search-key pipeline.
// To dedup repeated tokens or strip stopwords, compose on top:
//
//   textnorm.SearchPreset().
//       SplitTokens().
//       DedupTokens().
//       RemoveStopwords(stopwords.English).
//       JoinTokens(" ")
```

### 4. Tests

- `textnorm/tokens_test.go` — extend with cases:
  - `DedupTokens`: empty input, single token, all-duplicate, mixed, order preserved on first-seen, case-sensitive (`"Red red"` → `["Red","red"]` unless `FoldCase` runs first).
  - `RemoveStopwords`: nil set is no-op, empty input, every token removed, no token removed, partial removal, case-sensitive matching.
- `textnorm/stopwords/stopwords_test.go`:
  - `English` contains expected anchors (`"the"`, `"a"`, `"is"`, `"are"`, `"of"`, `"with"`).
  - `Union` merges correctly with no duplicates and is independent of input order.
  - `Union` of zero inputs returns a non-nil empty map.

Existing fuzz/bench/regression tests in `textnorm/` continue to run; the new stages compose into the existing pipeline without modification to runners.

## framework-golang Changes

### 1. `go.mod` — add dependency

```
require github.com/alessiosavi/GoGPUtils <version-tagged-after-step-1>
```

The version is the tag cut after the GoGPUtils additions land.

### 2. `pkg/text/` — new package

**`pkg/text/doc.go`** — package-level godoc. Mirrors the Usage.md decision table so callers see it on `go doc`.

**`pkg/text/text.go`:**

```go
package text

import (
    "github.com/alessiosavi/GoGPUtils/textnorm"
    "github.com/alessiosavi/GoGPUtils/textnorm/stopwords"
)

// Clean applies a single fixed-order canonical clean: UTF-8 sanitize →
// Unicode normalize (diacritics stripped) → fold case → trim → collapse
// whitespace. Use for free-text fields stored or compared as-is.
// Wraps textnorm.CanonicalPreset.
func Clean(s string) (string, error)

// SearchKey produces a search-friendly form: Clean + diacritic strip +
// non-letter/digit/space removal + token split + rejoin. Use for search
// indexing, query matching, deterministic comparison. Wraps textnorm.SearchPreset.
func SearchKey(s string) (string, error)

// DBSafe produces persistence-safe text: UTF-8 sanitize, Unicode normalize
// (NFD + drop combining marks + NFC — diacritics are stripped), trim, collapse
// whitespace. Does NOT lowercase. Use for storing scraper/normalizer output
// where you want stable byte-for-byte storage but want to keep original case.
// Wraps textnorm.DBSafePreset.
func DBSafe(s string) (string, error)

// Slug produces a URL-friendly slug. Use for URLs, filenames, IDs.
// Wraps GoGPUtils stringutil.Slugify.
func Slug(s string) (string, error)

// Pipeline returns a fresh GoGPUtils textnorm.Pipeline so callers can compose
// custom flows (DedupTokens, RemoveStopwords, etc.) without importing GoGPUtils
// directly. Use this for embedding-input cleaning, cross-field dedup
// composition, or any flow not covered by the four presets above.
func Pipeline() textnorm.Pipeline { return textnorm.New() }

// Stopwords re-exports the GoGPUtils stopwords sets so callers stay on the
// framework import path.
var Stopwords = struct {
    English, French, Italian map[string]struct{}
}{
    English: stopwords.English,
    French:  stopwords.French,
    Italian: stopwords.Italian,
}
```

**`pkg/text/presets.go`** — preset construction helpers (private; called by `Clean`/`SearchKey`/`DBSafe`). Single source of truth for which `textnorm` preset each public function uses, so future preset swaps are one-line changes.

### 3. `pkg/request/params.go` — additive

Existing `GetQueryParams` and `sanitizeQuery` unchanged. Add:

```go
// GetQueryParamsSearchKey returns the named query parameter passed through
// the framework's search-key normalization (text.SearchKey). Use for
// endpoints that match against search-indexed columns.
func GetQueryParamsSearchKey(r events.APIGatewayProxyRequest, paramName string) (string, error) {
    paramValue, exists := r.QueryStringParameters[paramName]
    if !exists || paramValue == "" {
        return "", errors.New("query parameter is required")
    }
    return text.SearchKey(paramValue)
}
```

### 4. Tests

- `pkg/text/text_test.go` — one happy-path test per public function (`Clean`, `SearchKey`, `DBSafe`, `Slug`). Light because GoGPUtils owns the deep coverage; these tests verify the routing wires up correctly.
- `pkg/request/params_test.go` — add `TestGetQueryParamsSearchKey` mirroring the existing `TestGetQueryParams` pattern.

### 5. Documentation

**`Usage.md`** — append a new section, "Package `text` - String Cleaning & Normalization". Contains the decision table, the "when NOT to clean" list, and worked examples for search-query, DB persistence, and embeddings hygiene. Full content is in the brainstorming transcript and reproduced in the implementation plan.

**`CLAUDE.md`** (new top-level file) — short pointer:

```markdown
# Framework-Golang — Claude Code Notes

Primary reference: `Usage.md`.

## Text cleaning

When a handler accepts user-supplied or scraped strings, prefer
`pkg/text` over inline cleaning. Decision table is in `Usage.md`
under "Package `text`". Do not import GoGPUtils directly from
service code — go through `pkg/text` so framework upgrades stay
mechanical.
```

## Data Flow Examples

**Search query, end to end:**

```
APIGatewayProxyRequest                                       events.APIGatewayProxyRequest
        │
        ▼
request.GetQueryParamsSearchKey(r, "q")                      "Café Hermès "
        │
        ▼
text.SearchKey                                               "cafe hermes"
        │
        ▼
textnorm.SearchPreset()                                      pipeline:
                                                              SanitizeUTF8 → NormalizeUnicode →
                                                              FoldCase → FilterRunes(letter|digit|space) →
                                                              TrimSpace → CollapseWhitespace →
                                                              SplitTokens → JoinTokens(" ")
```

**Embeddings hygiene (consumer composes):**

```go
// veliu.com/lib/services/embeddings/<future migration>
clean, err := text.Pipeline().
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

## Error Handling

- All preset and pipeline functions return `(string, error)`. Errors come exclusively from `transform.Chain` failures — extremely rare for valid Go strings, but propagated faithfully.
- `request.GetQueryParamsSearchKey` returns the underlying error from `text.SearchKey` so handlers can return a `400 Bad Request` if normalization fails (it won't, in practice, but the contract is preserved).
- `RemoveStopwords(nil)` is a deliberate no-op — not an error. This lets callers wire a stopword set from config without branching.

## Rollout

1. **GoGPUtils PR** — `textnorm/tokens.go` additions + `textnorm/stopwords/` subpackage + tests. Tag a release.
2. **framework-golang PR** — `pkg/text/` package + `pkg/request/` addition + `Usage.md` section + new `CLAUDE.md`. Pin the GoGPUtils dep to the tag from step 1. Tag a release.
3. **veliu.com migration is out of scope.** Future PRs touching normalizer/embeddings text handling migrate to `pkg/text` opportunistically, guided by the new docs.

Each step is fully additive: no existing caller breaks, no existing test changes its expected output.

## Open Questions

None. All design questions resolved in brainstorming:

- Dependency direction: `framework-golang` depends on `GoGPUtils` (Q1 answered "A").
- Dedup semantics scope: GoGPUtils ships generic primitives only; cross-field dedup is consumer logic (Q2 answered "B").
- `DedupTokens` casing: case-sensitive primitive, pair with `FoldCase()` upstream (Q3a).
- Stopwords: subpackage with English baked in, French/Italian starter sets (Q3b).
- Package location: `pkg/text/` (Q4a).
- Default normalization: opt-in `GetQueryParamsSearchKey`, existing `GetQueryParams` unchanged (Q4b).
- Surface shape: thin re-export under framework names (Q4c).
