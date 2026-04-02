# Architecture Patterns

**Domain:** Go text-normalization pipeline library
**Researched:** 2026-04-02

## Recommended Architecture

Use one dedicated top-level package, e.g. `textnorm/`, with an immutable fluent pipeline at the center. The package should expose small, composable stages that each take a string and return a string plus error; the pipeline just stores stage order and executes them left-to-right. Keep `stringutil` as a separate legacy/helper package; do not make the new pipeline depend on it.

For Unicode work, lean on `golang.org/x/text/transform` and `golang.org/x/text/unicode/norm` for normalization stages, because the official docs describe them as the canonical byte/string transformers and a `Chain`-style composition model. Use that same mental model, but at the public API level keep the pipeline fluent and string-oriented.

### Component Boundaries

| Component | Responsibility | Communicates With |
|-----------|---------------|-------------------|
| `Pipeline` | Owns ordered stages, executes them, returns final string/error | Public stages, presets, adapters |
| `Stage` / `Transformer` abstraction | A single normalization step with explicit input/output contract | Pipeline executor only |
| Core text stages | Unicode normalization, case folding, whitespace cleanup, trimming, replacement, filtering | Pipeline |
| Token stages | Split, map, drop, join, normalize per-token | Tokenizer + pipeline |
| Unicode adapter | Bridges pipeline stages to `x/text` normalization primitives | `transform` and `norm` |
| Presets | Reusable canned pipelines for search/indexing/cleanup | Core stages and pipeline builder |
| Compatibility helpers | Optional wrappers for existing `stringutil` behavior | New package only; no reverse dependency |

### Data Flow

1. Caller creates a pipeline with a fluent constructor, e.g. `textnorm.New().Normalize(...).TrimSpace().CollapseWhitespace()`.
2. Pipeline stores stages immutably; each call returns a new pipeline value.
3. `Run` (or similar) takes raw input text and feeds it through stages in order.
4. Text stages operate on the whole string first: Unicode normalization, entity/markup cleanup if allowed, case folding, punctuation filtering, whitespace normalization.
5. Token stages run after whole-string cleanup when token boundaries matter: split → per-token map/filter → join.
6. Final stage emits a deterministic normalized string; errors short-circuit immediately.

### Suggested Build Order

1. **Pipeline contract first** — define the stage type, immutable pipeline value, and `Run` execution path.
2. **Core string cleanup next** — whitespace, trimming, filtering, replacement, and simple case transforms.
3. **Unicode normalization layer** — add `x/text`-backed NFC/NFKC/NFKD and diacritic removal stages.
4. **Token pipeline** — split/map/filter/join once the whole-string stages are stable.
5. **Presets and adapters** — publish opinionated ready-made pipelines and optional `transform.Transformer` bridging last.

### Recommended Package Structure

```text
textnorm/
├── doc.go
├── pipeline.go        # Pipeline type, fluent builder, Run method
├── stages.go          # Stage interface/function types and common helpers
├── unicode.go         # norm/transform-backed normalization stages
├── whitespace.go      # trim/collapse/dedent-like cleanup
├── tokens.go          # split/map/filter/join token stages
├── presets.go         # opinionated ready-made pipelines
├── adapter.go         # optional transform.Transformer bridge
└── *_test.go
```

Keep helpers private unless they are clearly reusable and stable. If the package grows, prefer more files in the same package over new subpackages; this repo’s style is one package per top-level directory.

### Patterns to Follow

### Pattern 1: Immutable fluent builder
**What:** Each method returns a new pipeline value with one more stage appended.
**When:** Always; it avoids shared mutable state and makes reuse safe.
**Example:**
```go
pipe := textnorm.New().NormalizeUnicode().RemoveAccents().CollapseWhitespace()
```

### Pattern 2: Single-pass ordered execution
**What:** Execute stages exactly in declaration order, with no hidden reordering.
**When:** For all normalization flows.
**Example:**
```go
out, err := pipe.Run(input)
```

### Pattern 3: Token-stage isolation
**What:** Only tokenize once, then operate per token, then join once.
**When:** When a stage needs word-level cleanup or filtering.
**Example:**
```go
textnorm.New().SplitWords().DropEmpty().Map(strings.ToLower).Join(" ")
```

### Pattern 4: Adapter, not duplicate Unicode logic
**What:** Wrap `x/text` primitives instead of reimplementing normalization tables.
**When:** For NFC/NFKC/NFKD and diacritic removal.
**Example:**
```go
// internally use transform.Chain(norm.NFKD, runes.Remove(...), norm.NFC)
```

## Anti-Patterns to Avoid

### Anti-Pattern 1: Global normalization mode
**What:** Package-level flags that change behavior for all callers.
**Why bad:** Breaks determinism and tests; impossible to reason about reuse.
**Instead:** Keep all configuration on the pipeline instance.

### Anti-Pattern 2: Mixed token and string semantics in one stage
**What:** A stage that sometimes rewrites whole strings and sometimes rewrites tokens.
**Why bad:** Creates order bugs and unclear invariants.
**Instead:** Separate whole-string stages from token stages.

### Anti-Pattern 3: Reusing `stringutil` as the public API
**What:** Extending the old helper bag with pipeline behavior.
**Why bad:** The existing package is function-centric; the new feature needs a distinct abstraction.
**Instead:** Introduce a new package and keep `stringutil` as a compatibility source only.

## Scalability Considerations

| Concern | At 100 users | At 10K users | At 1M users |
|---------|--------------|--------------|-------------|
| Allocation cost | Simple string copies are fine | Reuse builders and avoid intermediate slices in hot paths | Add benchmarks and fast paths for ASCII/no-op input |
| Unicode cost | `x/text` transforms are acceptable | Cache preset pipelines, not mutable state | Consider exposing streaming adapters for large inputs |
| Tokenization | Basic split/join is enough | Keep tokenizer configurable but deterministic | Avoid per-token allocations where possible |
| API stability | One fluent API is enough | Keep stage names narrow and explicit | Preserve backwards-compatible stage semantics |

## Sources

- Go package organization and one-package-per-directory guidance: https://go.dev/doc/code
- `golang.org/x/text/transform` docs (`Transformer`, `Chain`, `String`, `Reader`, `Writer`): https://pkg.go.dev/golang.org/x/text/transform
- `golang.org/x/text/unicode/norm` docs (`Form`, `String`, `Transform`, `Chain`-compatible behavior): https://pkg.go.dev/golang.org/x/text/unicode/norm
- Existing repo conventions and current string helpers: `.planning/codebase/ARCHITECTURE.md`, `.planning/codebase/STRUCTURE.md`, `stringutil/clean.go`, `stringutil/stringutil.go`
