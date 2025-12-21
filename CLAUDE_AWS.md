## Claude Prompt: Go AWS Utilities v2 Redesign

You are a **senior Go engineer and cloud SDK designer**, with deep expertise in:

* Idiomatic Go
* Large utility libraries
* AWS SDK for Go (v2 preferred)
* Clean architecture and long-term API design

Your task is to **design and implement a new AWS helper utilities suite** in a **new `aws` subfolder**, based on analysis of existing code and a previous implementation.

---

## Context

This repository contains:

1. A **large collection of general-purpose Go utilities**, including (but not limited to):

   * `fileutils`
   * `arrayutils`
   * `datastructure`
   * `mathutils`
   * `cryptoutils`
   * etc.

2. A historical **`aws` folder** in a **specific commit** which provides an earlier AWS helper implementation.

3. A **current branch (`main`)** where the codebase has already been migrated and cleaned up.

Your job is to create a **new versioned implementation** in the folder `/aws`

---

## Objectives

### 1. Read & Understand Existing Code

* Read the current utility packages to understand:

  * Naming conventions
  * Error-handling style
  * Package layout
  * API philosophy
* Identify reusable patterns or anti-patterns to avoid
* Do **not** modify existing utility packages

---

### 2. Analyze the Historical AWS Implementation

* Read the `aws` folder from **this specific commit**:

  * b7a0843b35f6f458abd382dcb448ca2c1c748683
* Identify:

  * Which AWS services are supported
  * Common helper patterns used
  * Configuration strategies
  * Repeated logic or inconsistencies
  * Gaps or incomplete coverage

---

### 3. Design the New `aws` Package

Create a **complete AWS helper utilities suite**, inspired by the old one but **fully redesigned**.

#### Design Requirements

* Use **AWS SDK for Go v2**
* Follow **idiomatic Go conventions**
* Prefer **explicit configuration over magic**
* Favor **small, focused packages**
* No global state
* Context-aware APIs
* Clear error wrapping and propagation
* Minimal public API surface
* Safe defaults

#### Package Structure (Example – adapt as needed)

```
aws/
├── config/
│   └── config.go
├── clients/
│   ├── s3.go
│   ├── dynamodb.go
│   ├── sqs.go
│   └── ...
├── helpers/
│   ├── retry.go
│   ├── paginator.go
│   └── ...
├── types/
│   └── options.go
├── internal/
│   └── validation/
└── README.md
```

You may adjust this structure if you justify it.

---

### 4. Implement AWS Helpers

* Provide helpers for the AWS services already present in the old implementation
* Expand coverage where logical (e.g., missing but commonly needed helpers)
* Each helper must:

  * Be composable
  * Be testable without real AWS calls (interfaces, mocks)
  * Avoid leaky abstractions over the AWS SDK
* Do **not** re-export the AWS SDK blindly

---

### 5. Testing & Benchmarks

* Provide:

  * Unit tests (table-driven)
  * Mocked AWS clients
  * Edge-case coverage
  * CI/CD files for github-actions
* Add benchmarks for:

  * Serialization/deserialization
  * Helper logic that may be performance-critical
* Tests must not require real AWS credentials

---

### 6. Documentation

Include:

* `aws/README.md`
* Clear usage examples
* Design philosophy
* Migration notes from old `aws` helpers
* Guidance on extending the helpers

---

## Output Requirements

### Structure your response as follows

1. **Analysis of current utility packages**
2. **Critique of the old AWS implementation**
3. **New AWS v2 architecture**
4. **Package layout explanation**
5. **Full implementation (Go code)**
6. **Tests and benchmarks**
7. **README and usage examples**

---

## Constraints

* Treat this as a **long-term maintained open-source module**
* Assume usage in:

  * Production systems
  * Concurrent environments
  * Infrastructure automation
* Avoid speculative abstractions
* Every exported symbol must be justified

---

## Quality Bar

Assume this code will be:

* Reviewed by AWS SDK contributors
* Used in latency-sensitive systems
* Maintained across multiple major versions

If any part does not meet that bar, redesign it.

---

If needed, ask **clarifying questions only about**:

* The commit hash
* Target AWS services
* Required backward compatibility

Otherwise, proceed autonomously.

---

If you want, I can also:

* Add **strict backward-compatibility constraints**
* Generate **migration adapters**
* Produce **architecture diagrams**
* Convert this into a **multi-step execution prompt**

Just tell me.
