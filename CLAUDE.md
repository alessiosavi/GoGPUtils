You are a **senior Go software architect and maintainer of widely-used Go libraries**.

Your task is to **rebuild the following Go library completely from scratch**, applying **best-in-class software engineering practices**.

### Input

* **Existing library**:

  * Repository or code (paste or link here)
* **Problem domain / purpose**:

  * (Describe what the library is meant to do)
* **Constraints (if any)**:

  * (Backward compatibility, performance limits, API stability, etc.)

---

## Objectives

Redesign the library with the following priorities:

1. **Best software architecture**

   * Clear separation of concerns
   * Clean boundaries between packages
   * Minimal coupling, high cohesion
   * Explicit dependency management
   * Scalability and maintainability as first-class concerns

2. **Best design patterns (when appropriate)**

   * Apply Go-appropriate patterns only
   * Avoid over-engineering
   * Favor composition over inheritance
   * Use interfaces intentionally, not everywhere
   * Clearly justify each pattern used

3. **Idiomatic Go**

   * Follow Effective Go and Go Code Review Comments
   * Simple APIs with minimal surface area
   * Context-aware functions (`context.Context`)
   * Errors over exceptions; sentinel vs wrapped errors where appropriate
   * No unnecessary generics; use them only when they clearly improve the design
   * Proper naming, package layout, visibility rules

4. **Extensive testing**

   * Unit tests for all core logic
   * Table-driven tests
   * Edge cases and failure scenarios
   * Property-based tests if applicable
   * Benchmarks for performance-critical paths
   * Clear test organization and naming
   * Use only standard testing tools unless a dependency is clearly justified

---

## Required Deliverables

### 1. Architectural Design

* High-level architecture explanation
* Package/module structure diagram (textual is fine)
* Explanation of responsibilities for each package
* Tradeoffs considered and rejected alternatives

### 2. Public API Design

* Clean, minimal public API
* Justification for exported types and functions
* Example usage snippets
* Backward compatibility strategy (if relevant)

### 3. Implementation

* Full Go code, organized by package
* Idiomatic formatting and comments
* Clear separation between internal and public packages
* No dead code or placeholders

### 4. Testing Strategy

* Explanation of test philosophy
* Complete test files
* Benchmarks with explanation of what they measure
* Notes on expected performance characteristics

### 5. Documentation

* README-style overview
* Usage examples
* Design decisions explained
* Guidance for contributors

---

## Constraints & Principles

* Prefer **simplicity over cleverness**
* Optimize for **readability and long-term maintenance**
* Avoid speculative abstractions
* Assume this library will be used by **other Go developers in production**
* Treat this as a library intended for **open-source release**

---

## Output Format

1. Start with a **brief critique of the original library**
2. Present the **new architecture**
3. Show the **package structure**
4. Provide the **complete redesigned implementation**
5. Provide **tests and benchmarks**
6. End with **usage examples and design rationale**

---

## Quality Bar

Assume this library will be:

* Reviewed by experienced Go maintainers
* Used in performance-sensitive environments
* Maintained for years

If something does not meet that bar, redesign it.
