# Feature Landscape

**Domain:** Go text-normalization pipeline library  
**Researched:** 2026-04-02

## Table Stakes

Features mature normalization libraries almost always provide. Missing these makes the package feel like a toy.

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| Unicode normalization (NFC/NFKC/NFD/NFKD) | Canonical text handling starts here; mature ecosystem libs expose it directly. | Med | Foundation for search, dedupe, and stable comparisons. |
| Accent/diacritic removal | Search and canonicalization workflows need accent-insensitive matching. | Low/Med | Depends on Unicode decomposition first. |
| Case folding / lowercasing | Normalized search keys should not care about case. | Low/Med | Prefer full folding semantics over naive casing when possible. |
| Whitespace trim + collapse | Scraped and user-entered text is almost always messy. | Low | Usually applied after token cleanup or punctuation stripping. |
| Invalid UTF-8 + NUL sanitization | Dirty input and database safety are routine requirements. | Low | Must happen before length truncation. |
| Rune/class filtering and mapping | Pipelines need to drop punctuation, keep only selected scripts, or rewrite characters. | Med | Core primitive for most cleaning stages. |
| Split / map / filter / join token pipeline | This is the actual fluent composition model users expect. | Med | Token stages depend on a stable split boundary and deterministic order. |
| Deterministic, reusable stage order | Normalization must be predictable for tests and indexing. | Low | No hidden reordering, no global state. |
| Canonical-key helpers | Users need one-step output for dedupe, map keys, and search index keys. | Low/Med | Derived from normalization + case folding + separator rules. |

## Differentiators

Features that make the library better than a bag of helpers, without turning it into a full text framework.

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| Immutable fluent pipelines | Callers can build and reuse named pipelines safely. | Med | Strong fit for Go utility style and avoids shared mutable config. |
| Opinionated presets (`Search`, `Slug`, `Canonical`, `DBSafe`) | Gives users a high-quality default instead of forcing them to assemble everything. | Low/Med | Presets should be thin wrappers over the same core stages. |
| Per-token processing (`Split -> Map/Filter -> Join`) | Enables word-level cleanup, stopword-like filtering, and token rewrites. | Med | This is the clearest differentiator for a pipeline API. |
| Width folding / fullwidth-halfwidth canonicalization | Important for East Asian text and mixed-width scraped data. | Med | Use as an explicit stage, not a silent default. |
| Streaming adapters (`io.Reader` / `io.Writer`) | Helpful for large inputs and integration with scrapers or batch jobs. | Med | Nice-to-have; not required for the MVP. |
| Optional ASCII transliteration stage | Useful for slugs and legacy systems that only accept ASCII. | High | Keep separate from accent stripping; transliteration is a different operation. |
| Stage tracing / explain mode | Helps users debug why an input changed. | Med | Useful for adoption, but not a first-release requirement. |
| Configurable token boundaries | Lets callers choose whitespace, punctuation, or custom delimiters without rewriting stages. | Med | Nice extension once split/join exists. |

## Anti-Features

Things this package should deliberately not become.

| Anti-Feature | Why Avoid | What to Do Instead |
|--------------|-----------|-------------------|
| Full NLP stack | Stemming, lemmatization, and language modeling explode scope fast. | Keep the library to normalization and simple token cleanup. |
| HTML parsing / DOM extraction | Markup cleanup is a different problem than text normalization. | Require callers to extract text before running the pipeline. |
| Search indexing / ranking / query execution | This library should prepare text, not store or score it. | Return normalized output only. |
| Heavy transliteration databases by default | Large mapping tables and locale quirks bloat a small utility library. | Make transliteration optional, if included at all. |
| PRECIS / username-policy enforcement | Security profiles are specialized and pull in validation semantics. | Stay focused on normalization, not identity policy. |
| Global config or hidden mutable state | Breaks determinism and makes pipelines hard to reuse safely. | Keep all configuration on the pipeline value. |
| Giant regex-driven DSL | Too much expressiveness becomes hard to reason about and test. | Prefer explicit stages and small helpers. |
| Locale magic by default | Implicit locale rules make output surprising. | Make locale-specific behavior explicit and opt-in. |

## Feature Dependencies

```text
Unicode normalization → Accent/diacritic removal
Unicode normalization → Canonical-key helpers
Case folding → Canonical-key helpers
Split → Map/Filter/Join token pipeline
Split + Join → Token-stage presets
Invalid UTF-8/NUL sanitization → Truncation / DB-safe output
Rune/class filtering → Slug / canonical cleanup stages
Unicode normalization + width folding → Search-friendly presets
```

## MVP Recommendation

Prioritize:
1. Unicode normalization + diacritic removal
2. Whitespace cleanup + invalid UTF-8/NUL sanitization
3. Fluent split/map/filter/join pipeline

Defer: ASCII transliteration, streaming adapters, stage tracing, and locale-specific rules until the core pipeline is proven.

## Sources

- `golang.org/x/text/unicode/norm` docs (Unicode normalization forms, `String`, `Transform`, `Chain`-style iteration): https://pkg.go.dev/golang.org/x/text/unicode/norm
- `golang.org/x/text/transform` docs (composable transformers, `Chain`, `String`, `Reader`, `Writer`): https://pkg.go.dev/golang.org/x/text/transform
- `golang.org/x/text/runes` docs (rune sets, `Map`, `Remove`, `If`, `ReplaceIllFormed`): https://pkg.go.dev/golang.org/x/text/runes
- `golang.org/x/text/width` docs (canonical width folding and wide/narrow transforms): https://pkg.go.dev/golang.org/x/text/width
- `golang.org/x/text/cases` docs (language-aware case mapping/folding): https://pkg.go.dev/golang.org/x/text/cases
- `golang.org/x/text/secure/precis` docs (what a broader string-normalization/validation framework looks like, and why it is out of scope): https://pkg.go.dev/golang.org/x/text/secure/precis
- `go-unidecode` README (ASCII transliteration as a separate concern): https://raw.githubusercontent.com/mozillazg/go-unidecode/master/README.md
- Current repo context: `.planning/PROJECT.md`, `.planning/codebase/ARCHITECTURE.md`, `stringutil/clean.go`, `stringutil/stringutil.go`
