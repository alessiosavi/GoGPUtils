---
layout: default
title: randutil
parent: Packages
nav_order: 7
---

# randutil

Cryptographically secure random generation utilities for Go.
{: .fs-6 .fw-300 }

The `randutil` package provides a dual-interface approach to random number generation: **cryptographically secure** functions for security-sensitive applications, and **fast** functions for general-purpose use cases like testing, simulations, and sampling.

---

## Overview

Security is the primary design goal of this package. All functions prefixed with `Secure` use Go's `crypto/rand` package, which reads from the operating system's cryptographically secure random number generator (`/dev/urandom` on Unix, `CryptGenRandom` on Windows). This makes them suitable for:

- Password and token generation
- Cryptographic keys and nonces
- Session identifiers
- API keys and secrets
- Any security-sensitive random data

For non-security purposes like unit tests, Monte Carlo simulations, or game mechanics, the package provides a `Generator` type backed by `math/rand/v2` for significantly better performance.

---

## Installation

```go
import "github.com/alessiosavi/GoGPUtils/randutil"
```

---

## Character Sets

Predefined character sets for string generation:

| Constant       | Value                           |
| -------------- | ------------------------------- |
| `Digits`       | `"0123456789"`                  |
| `Lowercase`    | `"abcdefghijklmnopqrstuvwxyz"`  |
| `Uppercase`    | `"ABCDEFGHIJKLMNOPQRSTUVWXYZ"`  |
| `Letters`      | `Lowercase + Uppercase`         |
| `AlphaNumeric` | `Letters + Digits`              |
| `Hex`          | `"0123456789abcdef"`            |
| `Symbols`      | `"!@#$%^&*()_+-=[]{}\|;:,.<>?"` |
| `All`          | `AlphaNumeric + Symbols`        |

---

## Errors

| Error              | Description                                    |
| ------------------ | ---------------------------------------------- |
| `ErrInvalidLength` | Returned when length parameter is not positive |
| `ErrEmptyCharset`  | Returned when charset string is empty          |
| `ErrEmptySlice`    | Returned when operating on an empty slice      |

---

## Secure Random Generation (crypto/rand)

Use these functions for all security-sensitive applications. They provide cryptographically secure randomness suitable for passwords, tokens, keys, and identifiers.

### SecureBytes

```go
func SecureBytes(n int) ([]byte, error)
```

Returns `n` cryptographically secure random bytes.

**Example:**

```go
key, err := randutil.SecureBytes(32)
if err != nil {
    return err
}
// key is 32 random bytes suitable for AES-256
```

---

### SecureString

```go
func SecureString(length int, charset string) (string, error)
```

Returns a cryptographically secure random string of the specified length using characters from the given charset.

**Example:**

```go
// Generate a 16-character alphanumeric password
password, err := randutil.SecureString(16, randutil.AlphaNumeric)
if err != nil {
    return err
}
// password: "k9mP2vLqR5nX8wJt"

// Generate a 32-character hex token
token, err := randutil.SecureString(32, randutil.Hex)
if err != nil {
    return err
}
// token: "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4"
```

---

### SecureInt

```go
func SecureInt(max int) (int, error)
```

Returns a cryptographically secure random integer in the range `[0, max)`.

**Example:**

```go
n, err := randutil.SecureInt(100)
if err != nil {
    return err
}
// n is a random integer from 0 to 99
```

---

### SecureInt64

```go
func SecureInt64(max int64) (int64, error)
```

Returns a cryptographically secure random int64 in the range `[0, max)`.

**Example:**

```go
n, err := randutil.SecureInt64(1_000_000_000)
if err != nil {
    return err
}
```

---

### SecureID

```go
func SecureID() (string, error)
```

Generates a cryptographically secure random ID. Returns a 32-character hex string representing 128 bits of entropy.

**Example:**

```go
id, err := randutil.SecureID()
if err != nil {
    return err
}
// id: "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4"
```

---

### SecureChoice

```go
func SecureChoice[T any](s []T) (T, error)
```

Returns a cryptographically secure random element from a slice.

**Example:**

```go
colors := []string{"red", "green", "blue", "yellow"}
color, err := randutil.SecureChoice(colors)
if err != nil {
    return err
}
// color is randomly selected from the slice
```

---

## Fast Random Generation (math/rand/v2)

Use the `Generator` type for non-security purposes where performance matters, such as tests, simulations, sampling, or game mechanics. The `Generator` is automatically seeded with cryptographically secure randomness on creation.

### Creating a Generator

```go
// Automatic seeding with crypto-secure randomness
rng := randutil.NewGenerator()

// Deterministic seed for reproducible sequences (useful in tests)
rng := randutil.NewGeneratorWithSeed(42)
```

---

### Generator Methods

#### Int

```go
func (g *Generator) Int(max int) int
```

Returns a random integer in the range `[0, max)`. Panics if `max <= 0`.

```go
rng := randutil.NewGenerator()
n := rng.Int(100)  // Random integer from 0 to 99
```

---

#### IntRange

```go
func (g *Generator) IntRange(min, max int) int
```

Returns a random integer in the range `[min, max]`. Panics if `min > max`.

```go
rng := randutil.NewGenerator()
dice := rng.IntRange(1, 6)  // Random dice roll: 1 to 6
```

---

#### Int64

```go
func (g *Generator) Int64(max int64) int64
```

Returns a random int64 in the range `[0, max)`.

---

#### Int64Range

```go
func (g *Generator) Int64Range(min, max int64) int64
```

Returns a random int64 in the range `[min, max]`.

---

#### Float64

```go
func (g *Generator) Float64() float64
```

Returns a random float64 in the range `[0.0, 1.0)`.

```go
rng := randutil.NewGenerator()
f := rng.Float64()  // Random float between 0.0 and 1.0
```

---

#### Float64Range

```go
func (g *Generator) Float64Range(min, max float64) float64
```

Returns a random float64 in the range `[min, max)`.

```go
rng := randutil.NewGenerator()
temp := rng.Float64Range(-10.0, 40.0)  // Random temperature
```

---

#### Float32

```go
func (g *Generator) Float32() float32
```

Returns a random float32 in the range `[0.0, 1.0)`.

---

#### Bool

```go
func (g *Generator) Bool() bool
```

Returns a random boolean value.

```go
rng := randutil.NewGenerator()
if rng.Bool() {
    // 50% chance of executing
}
```

---

#### Bytes

```go
func (g *Generator) Bytes(n int) []byte
```

Returns `n` random bytes. **Note:** Not cryptographically secure.

---

#### String

```go
func (g *Generator) String(length int, charset string) string
```

Returns a random string of the specified length using the given charset. **Note:** Not cryptographically secure.

---

#### AlphaNumericString

```go
func (g *Generator) AlphaNumericString(length int) string
```

Returns a random alphanumeric string. **Note:** Not cryptographically secure.

---

#### ChoiceInt

```go
func (g *Generator) ChoiceInt(s []int) int
```

Returns a random element from an int slice. Panics if slice is empty.

---

#### ChoiceString

```go
func (g *Generator) ChoiceString(s []string) string
```

Returns a random element from a string slice. Panics if slice is empty.

---

#### ShuffleInts

```go
func (g *Generator) ShuffleInts(s []int)
```

Shuffles an int slice in place.

---

#### ShuffleStrings

```go
func (g *Generator) ShuffleStrings(s []string)
```

Shuffles a string slice in place.

---

#### SampleInts

```go
func (g *Generator) SampleInts(s []int, n int) []int
```

Returns `n` random ints from the slice without replacement. If `n > len(s)`, returns all elements shuffled.

---

#### WeightedChoice

```go
func (g *Generator) WeightedChoice(weights []float64) int
```

Returns a random index based on weights. Weight values don't need to sum to 1.0.

```go
rng := randutil.NewGenerator()
weights := []float64{0.7, 0.2, 0.1}  // 70%, 20%, 10%
idx := rng.WeightedChoice(weights)
```

---

#### Probability

```go
func (g *Generator) Probability(p float64) bool
```

Returns `true` with the given probability (0.0 to 1.0).

```go
rng := randutil.NewGenerator()
if rng.Probability(0.75) {
    // 75% chance of executing
}
```

---

## Generic Functions

These generic functions work with any type and operate on the `Generator`:

### Choice

```go
func Choice[T any](g *Generator, s []T) T
```

Returns a random element from a slice of any type.

```go
rng := randutil.NewGenerator()
items := []string{"apple", "banana", "cherry"}
fruit := randutil.Choice(rng, items)
```

---

### ChoiceN

```go
func ChoiceN[T any](g *Generator, s []T, n int) []T
```

Returns `n` random elements with replacement.

---

### Sample

```go
func Sample[T any](g *Generator, s []T, n int) []T
```

Returns `n` random elements without replacement.

```go
rng := randutil.NewGenerator()
numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
pick5 := randutil.Sample(rng, numbers, 5)  // 5 unique random numbers
```

---

### Shuffle

```go
func Shuffle[T any](g *Generator, s []T)
```

Shuffles a slice in place.

```go
rng := randutil.NewGenerator()
cards := []string{"A♠", "K♠", "Q♠", "J♠", "10♠"}
randutil.Shuffle(rng, cards)
```

---

### ShuffleCopy

```go
func ShuffleCopy[T any](g *Generator, s []T) []T
```

Returns a new shuffled slice without modifying the original.

```go
rng := randutil.NewGenerator()
original := []int{1, 2, 3, 4, 5}
shuffled := randutil.ShuffleCopy(rng, original)
// original is unchanged
```

---

## Sequence Generation

Utility functions for generating integer sequences:

### Sequence

```go
func Sequence(start, n int) []int
```

Returns a slice of `n` sequential integers starting from `start`.

```go
randutil.Sequence(0, 5)   // [0, 1, 2, 3, 4]
randutil.Sequence(10, 3)  // [10, 11, 12]
```

---

### Range

```go
func Range(start, end int) []int
```

Returns a slice of integers from `start` to `end` (exclusive).

```go
randutil.Range(0, 5)   // [0, 1, 2, 3, 4]
randutil.Range(3, 7)   // [3, 4, 5, 6]
```

---

### RangeStep

```go
func RangeStep(start, end, step int) []int
```

Returns a slice of integers from `start` to `end` with the given step.

```go
randutil.RangeStep(0, 10, 2)   // [0, 2, 4, 6, 8]
randutil.RangeStep(10, 0, -2)  // [10, 8, 6, 4, 2]
```

---

## Usage Examples

### Generating Secure Tokens

```go
package main

import (
    "fmt"
    "github.com/alessiosavi/GoGPUtils/randutil"
)

func main() {
    // Generate a 32-character API token
    token, err := randutil.SecureString(32, randutil.AlphaNumeric)
    if err != nil {
        panic(err)
    }
    fmt.Println("API Token:", token)

    // Generate a 128-bit random ID
    id, err := randutil.SecureID()
    if err != nil {
        panic(err)
    }
    fmt.Println("ID:", id)
}
```

---

### Generating Secure Passwords

```go
package main

import (
    "fmt"
    "github.com/alessiosavi/GoGPUtils/randutil"
)

func main() {
    // Generate a strong 20-character password with symbols
    password, err := randutil.SecureString(20, randutil.All)
    if err != nil {
        panic(err)
    }
    fmt.Println("Password:", password)

    // Generate a PIN code
    pin, err := randutil.SecureString(6, randutil.Digits)
    if err != nil {
        panic(err)
    }
    fmt.Println("PIN:", pin)
}
```

---

### Random Sampling and Shuffling

```go
package main

import (
    "fmt"
    "github.com/alessiosavi/GoGPUtils/randutil"
)

func main() {
    rng := randutil.NewGenerator()

    // Shuffle a deck of cards
    cards := []string{"A♠", "K♠", "Q♠", "J♠", "10♠", "9♠", "8♠", "7♠"}
    randutil.Shuffle(rng, cards)
    fmt.Println("Shuffled:", cards)

    // Draw 5 random cards without replacement
    hand := randutil.Sample(rng, cards, 5)
    fmt.Println("Hand:", hand)

    // Weighted random selection
    prizes := []string{"Grand Prize", "Second Prize", "Consolation"}
    weights := []float64{0.05, 0.15, 0.80}
    winner := prizes[rng.WeightedChoice(weights)]
    fmt.Println("Winner:", winner)
}
```

---

### Deterministic Randomness for Tests

```go
package main

import (
    "fmt"
    "github.com/alessiosavi/GoGPUtils/randutil"
)

func main() {
    // Same seed produces the same sequence every time
    rng := randutil.NewGeneratorWithSeed(42)

    for i := 0; i < 5; i++ {
        fmt.Println(rng.Int(100))
    }
    // Always prints the same 5 numbers
}
```

---

## crypto/rand vs math/rand: When to Use Each

| Use Case                      | Recommended API               | Reason                                   |
| ----------------------------- | ----------------------------- | ---------------------------------------- |
| Passwords, tokens, API keys   | `SecureString`, `SecureBytes` | Cryptographically secure, unpredictable  |
| Session IDs, CSRF tokens      | `SecureID`, `SecureString`    | Must be unpredictable to prevent attacks |
| Cryptographic keys, nonces    | `SecureBytes`                 | Required for security guarantees         |
| Lottery, gambling, giveaways  | `SecureInt`, `SecureChoice`   | Must be fair and unpredictable           |
| Unit tests, simulations       | `Generator` methods           | Faster, reproducible with seed           |
| Game mechanics (non-monetary) | `Generator` methods           | Good enough, much faster                 |
| Data sampling, shuffling      | `Generator` methods           | Performance matters                      |
| Load balancing, A/B testing   | `Generator` methods           | Statistical randomness is sufficient     |

### Key Differences

**`crypto/rand`** (Secure functions):

- Uses OS-level cryptographically secure random number generator
- Suitable for security-sensitive applications
- Slower due to system calls
- Can fail (returns errors)
- Unpredictable and non-reproducible

**`math/rand/v2`** (Generator methods):

- Uses deterministic pseudo-random number generator
- Much faster (no system calls)
- Cannot fail (no errors returned)
- Reproducible with `NewGeneratorWithSeed`
- **Not suitable for security applications**

### Security Warning

{: .warning }

> Never use `Generator` methods or `math/rand` for security-sensitive operations like password generation, token creation, or cryptographic purposes. Always use the `Secure*` functions which rely on `crypto/rand`.

---

## Performance Considerations

The secure functions are significantly slower than the fast generator due to system calls to the OS random source. Benchmark results on a typical machine:

| Operation                            | Approximate Time |
| ------------------------------------ | ---------------- |
| `SecureBytes(32)`                    | ~5-10 µs         |
| `SecureString(32, AlphaNumeric)`     | ~10-15 µs        |
| `Generator.Int(1000)`                | ~50 ns           |
| `Generator.String(32, AlphaNumeric)` | ~200 ns          |
| `Generator.ShuffleInts(1000)`        | ~5 µs            |

For high-throughput scenarios requiring secure randomness, consider batching secure random byte generation and deriving multiple values from a single batch.
