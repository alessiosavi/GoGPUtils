# Technology Stack

**Project:** Text Normalization Pipeline
**Researched:** 2026-04-02

## Recommended Stack

### Core Framework
| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| Go | 1.25.0 | Library runtime, stdlib text handling, tests | Already pinned in `go.mod`; keeps the package lightweight, dependency-free at runtime, and aligned with the repo’s current toolchain. |
| `golang.org/x/text` | v0.35.0 | Unicode normalization and rune transforms | This is the right non-stdlib dependency for a serious text-normalization pipeline. It gives you `transform.Chain`, `unicode/norm`, and `runes` with tagged, Go-native APIs. |

### Database
| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| None | — | — | This package should normalize text *before* persistence/search, not own storage. |

### Infrastructure
| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| None | — | — | No server, queue, or background worker is needed for a pure library package. |

### Supporting Libraries
| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| `golang.org/x/text/transform` | v0.35.0 | Compose pipeline stages | Use for the fluent chain execution layer; `transform.Chain` is the canonical way to sequence normalization steps. |
| `golang.org/x/text/unicode/norm` | v0.35.0 | Unicode NFC/NFD/NFKC/NFKD normalization | Use for canonicalization, diacritic stripping, and stable search keys. Prefer this over ad hoc rune hacks. |
| `golang.org/x/text/runes` | v0.35.0 | Rune filtering/mapping/removal | Use for step-level transforms like removing combining marks, filtering control chars, or custom predicates. |
| `golang.org/x/text/width` | v0.35.0 | Width folding | Optional only if you need half-width/full-width normalization for East Asian or scraped text. |
| Standard library: `strings`, `unicode`, `unicode/utf8`, `html`, `testing` | Go 1.25.0 | Whitespace cleanup, character classes, UTF-8 validation, entity decoding, tests | Keep all non-Unicode-heavy work in stdlib; it is enough for whitespace collapse, printable filtering, HTML entity decoding, and validation. |

## Alternatives Considered

| Category | Recommended | Alternative | Why Not |
|----------|-------------|-------------|---------|
| Unicode normalization | `golang.org/x/text/unicode/norm` | Hand-rolled accent stripping or regex-based cleanup | Incorrect Unicode coverage, brittle behavior, and poor long-term maintainability. |
| Pipeline composition | `golang.org/x/text/transform` | Custom one-off string functions only | You want reusable stages and explicit sequencing; `transform.Chain` already solves the composition problem well. |
| Rune removal/filtering | `golang.org/x/text/runes.Remove` / `Map` | `transform.RemoveFunc` | `RemoveFunc` is deprecated; docs direct users to `runes.Remove` instead. |
| HTML cleanup | `html.UnescapeString` + caller-side parsing | Building HTML parsing into this library | Parsing/extraction is out of scope; normalization should operate on text after upstream cleanup. |
| Heavy text stack | Stdlib + `x/text` | ICU/cgo or large NLP libs | Too heavy for a small library, adds deployment complexity, and goes beyond normalization into linguistic processing. |

## Installation

```bash
# Core dependency set
go get golang.org/x/text@v0.35.0

# Keep module metadata tidy
go mod tidy
```

## Sources

- `/opt/SP/Workspace/Go/GoGPUtils/go.mod` — repo already pins `go 1.25.0` and `golang.org/x/text v0.35.0`
- `/opt/SP/Workspace/Go/GoGPUtils/stringutil/clean.go` — existing normalization helpers already use `transform`, `norm`, `runes`, `html`, `strings`, `unicode`, and `utf8`
- https://pkg.go.dev/golang.org/x/text/transform — `transform.Chain`, `String`, `Bytes`, `Reader`, `Writer`
- https://pkg.go.dev/golang.org/x/text/unicode/norm — Unicode normalization forms and `Form.String`/`Form.Transform`
- https://pkg.go.dev/golang.org/x/text/runes — rune transforms and `Remove`, `Map`, `If`
- https://go.dev/doc/go1.25 — Go 1.25 release notes and toolchain baseline
- `go list -m -versions golang.org/x/text` — confirms `v0.35.0` is the current tagged module version in this environment
