# GoGPUtils

[![Go](https://github.com/alessiosavi/GoGPUtils/actions/workflows/go.yml/badge.svg)](https://github.com/alessiosavi/GoGPUtils/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/alessiosavi/GoGPUtils)](https://goreportcard.com/report/github.com/alessiosavi/GoGPUtils)
[![GoDoc](https://godoc.org/github.com/alessiosavi/GoGPUtils?status.svg)](https://godoc.org/github.com/alessiosavi/GoGPUtils)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> **v1 - Experimental**: This library is in its initial release phase. APIs may change in future versions.

A collection of well-tested, idiomatic Go utilities for common programming tasks. Zero external dependencies for core utilities.

## Installation

```bash
go get github.com/alessiosavi/GoGPUtils
```

## Design Philosophy

- **Errors over panics**: All functions return errors instead of panicking
- **Zero global state**: No singletons; all state is explicit
- **Generic when useful**: Uses generics to reduce duplication without over-abstraction
- **Minimal dependencies**: Core library has zero external dependencies
- **Context-aware**: Blocking operations accept `context.Context`

## Packages

| Package | Description |
|---------|-------------|
| [`sliceutil`](#sliceutil) | Generic slice operations (filter, map, reduce, chunk, etc.) |
| [`stringutil`](#stringutil) | String manipulation and similarity algorithms |
| [`textnorm`](#textnorm) | Deterministic text normalization pipelines |
| [`mathutil`](#mathutil) | Mathematical and statistical operations |
| [`fileutil`](#fileutil) | File system operations with proper error handling |
| [`cryptoutil`](#cryptoutil) | Secure AES-GCM encryption |
| [`randutil`](#randutil) | Cryptographically secure random generation |
| [`collection`](#collection) | Generic data structures (Stack, Queue, Set, BST) |
| [`aws`](./aws/README.md) | AWS SDK v2 helpers (S3, DynamoDB, SQS, SSM, Secrets Manager, Lambda) |

---

## sliceutil

Generic slice operations using Go generics.

```go
import "github.com/alessiosavi/GoGPUtils/sliceutil"
```

### Functions

| Function | Description |
|----------|-------------|
| `Filter[T](slice, predicate)` | Returns elements matching the predicate |
| `Map[T, U](slice, mapper)` | Transforms each element |
| `Reduce[T, U](slice, initial, reducer)` | Reduces slice to a single value |
| `Unique[T](slice)` | Returns unique elements |
| `Chunk[T](slice, size)` | Splits slice into chunks |
| `GroupBy[T, K](slice, keyFunc)` | Groups elements by key |
| `Partition[T](slice, predicate)` | Splits into matching/non-matching |
| `Flatten[T](slices)` | Flattens nested slices |
| `Contains[T](slice, value)` | Checks if slice contains value |
| `IndexOf[T](slice, value)` | Returns index of value (-1 if not found) |
| `Reverse[T](slice)` | Returns reversed slice |
| `FlatMap[T, U](slice, transform)` | Maps each element to a slice and flattens the result |
| `MapErr[T, U](slice, fn)` | Maps elements with error handling |
| `Intersection[T](a, b)` | Returns common elements |
| `Difference[T](a, b)` | Returns elements in a but not in b |
| `Union[T](a, b)` | Returns combined unique elements |

### Example

```go
numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

// Filter even numbers
evens := sliceutil.Filter(numbers, func(n int) bool {
    return n%2 == 0
})
// evens = [2, 4, 6, 8, 10]

// Double each number
doubled := sliceutil.Map(numbers, func(n int) int {
    return n * 2
})
// doubled = [2, 4, 6, 8, 10, 12, 14, 16, 18, 20]

// Sum all numbers
sum := sliceutil.Reduce(numbers, 0, func(acc, n int) int {
    return acc + n
})
// sum = 55

// Split into chunks
chunks := sliceutil.Chunk(numbers, 3)
// chunks = [[1, 2, 3], [4, 5, 6], [7, 8, 9], [10]]

// Group by even/odd
grouped := sliceutil.GroupBy(numbers, func(n int) string {
    if n%2 == 0 {
        return "even"
    }
    return "odd"
})
// grouped = map["even":[2, 4, 6, 8, 10], "odd":[1, 3, 5, 7, 9]]
```

---

## stringutil

String manipulation utilities and similarity algorithms.

```go
import "github.com/alessiosavi/GoGPUtils/stringutil"
```

### Functions

| Function | Description |
|----------|-------------|
| `Reverse(s)` | Reverses a string (Unicode-safe) |
| `Truncate(s, maxLen, suffix)` | Truncates to max length with suffix |
| `TruncateWords(s, maxLen, suffix)` | Truncates at a word boundary |
| `PadLeft(s, length, pad)` | Pads string on the left |
| `PadRight(s, length, pad)` | Pads string on the right |
| `PadCenter(s, length, pad)` | Centers string with padding |
| `IsEmpty(s)` | Checks if string is empty |
| `IsBlank(s)` | Checks if string is empty or whitespace |
| `IsAlpha(s)` | Checks if string is alphabetic |
| `IsNumeric(s)` | Checks if string is numeric |
| `SnakeCase(s)` | Converts to snake_case |
| `CamelCase(s)` | Converts to camelCase |
| `PascalCase(s)` | Converts to PascalCase |
| `KebabCase(s)` | Converts to kebab-case |
| `Words(s)` | Splits into words |
| `SplitAndTrim(s, sep)` | Splits and trims parts |

### Similarity Algorithms

| Function | Description |
|----------|-------------|
| `LevenshteinDistance(a, b)` | Edit distance between strings |
| `LevenshteinSimilarity(a, b)` | Normalized Levenshtein similarity |
| `DamerauLevenshteinDistance(a, b)` | Edit distance with transpositions |
| `JaroSimilarity(a, b)` | Jaro similarity score |
| `JaroWinklerSimilarity(a, b, prefixScale)` | Similarity score (0.0-1.0) |
| `DiceCoefficient(a, b)` | Dice coefficient similarity |
| `HammingDistance(a, b)` | Bit-level distance (same length strings) |

### Example

```go
// Case conversions
stringutil.SnakeCase("HelloWorld")      // "hello_world"
stringutil.CamelCase("hello_world")     // "helloWorld"
stringutil.PascalCase("hello_world")    // "HelloWorld"
stringutil.KebabCase("HelloWorld")      // "hello-world"

// Padding
stringutil.PadLeft("42", 5, '0')      // "00042"
stringutil.PadRight("Go", 5, '-')     // "Go---"
stringutil.Truncate("Hello World", 8, "...")

// Similarity
stringutil.LevenshteinDistance("kitten", "sitting")     // 3
stringutil.JaroWinklerSimilarity("hello", "hallo", 0.1) // ~0.88
```

---

## textnorm

Deterministic text normalization pipelines.

```go
import "github.com/alessiosavi/GoGPUtils/textnorm"
```

### Functions

| Function | Description |
|----------|-------------|
| `New()` | Creates an empty pipeline |
| `SearchPreset(opts...)` | Search-oriented normalization |
| `CanonicalPreset(opts...)` | Canonical normalization |
| `DBSafePreset(opts...)` | Persistence-safe normalization |
| `WithWidthFold()` | Option that folds full-width characters |

### Example

```go
textnorm.SearchPreset().Run("  Café, go!  ")
textnorm.CanonicalPreset().Run("  Hello, World!  ")
textnorm.DBSafePreset(textnorm.WithWidthFold()).Run("  Ｇｏ\x00  ")
```

---

## mathutil

Mathematical and statistical operations using generics.

```go
import "github.com/alessiosavi/GoGPUtils/mathutil"
```

### Functions

| Function | Description |
|----------|-------------|
| `Sum[T](values)` | Sum of all values |
| `Product[T](values)` | Product of all values |
| `Average[T](values)` | Arithmetic mean |
| `Min[T](values)` | Minimum value |
| `Max[T](values)` | Maximum value |
| `MinMax[T](values)` | Returns both min and max |
| `Clamp[T](value, min, max)` | Constrains value to range |
| `Abs[T](value)` | Absolute value |
| `Variance[T](values)` | Population variance |
| `StdDev[T](values)` | Standard deviation |
| `Median[T](values)` | Median value |
| `Mode[T](values)` | Most frequent value(s) |
| `Percentile[T](values, p)` | Percentile value |
| `GCD[T](a, b)` | Greatest common divisor |
| `LCM[T](a, b)` | Least common multiple |
| `IsPrime(n)` | Primality test |
| `Factorial(n)` | Factorial (n!) |
| `Fibonacci(n)` | Nth Fibonacci number |

### Example

```go
numbers := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

mathutil.Sum(numbers)        // 55.0
mathutil.Average(numbers)    // 5.5
mathutil.Min(numbers)        // 1.0
mathutil.Max(numbers)        // 10.0
mathutil.Median(numbers)     // 5.5
mathutil.StdDev(numbers)     // ~2.87

mathutil.Clamp(15, 0, 10)    // 10
mathutil.GCD(48, 18)         // 6
mathutil.IsPrime(17)         // true
```

---

## fileutil

File system operations with proper error handling.

```go
import "github.com/alessiosavi/GoGPUtils/fileutil"
```

### Functions

| Function | Description |
|----------|-------------|
| `Exists(path)` | Checks if path exists |
| `IsFile(path)` | Checks if path is a file |
| `IsDir(path)` | Checks if path is a directory |
| `IsSymlink(path)` | Checks if path is a symlink |
| `IsExecutable(path)` | Checks if path is executable |
| `ReadBytes(ctx, path)` | Reads entire file as bytes |
| `ReadString(ctx, path)` | Reads entire file as string |
| `ReadLines(ctx, path)` | Reads all lines from a file |
| `WriteBytes(path, data, perm)` | Writes data to file |
| `WriteString(path, content, perm)` | Writes string content to file |
| `AppendString(path, content, perm)` | Appends string content to file |
| `EnsureDir(path, perm)` | Creates directory if not exists |
| `List(ctx, dir, opts)` | Lists files in directory |
| `Find(ctx, dir, pattern)` | Finds files by pattern |
| `Copy(ctx, src, dst)` | Copies a file |
| `Move(ctx, src, dst)` | Moves a file |
| `Touch(path)` | Creates or updates a file timestamp |
| `Size(path)` | Returns file size |
| `ModTime(path)` | Returns modification time |

### Example

```go
// Check existence
if fileutil.Exists("/path/to/file") {
    // File exists
}

// Read file
content, err := fileutil.ReadString(context.Background(), "/path/to/file.txt")
if err != nil {
    return err
}

// Write file
err := fileutil.WriteString("/path/to/file.txt", "Hello", 0644)

// List files
files, err := fileutil.List(context.Background(), "/path/to/dir", 0)
for _, f := range files {
    fmt.Println(f)
}
```

---

## cryptoutil

Secure AES-GCM encryption utilities.

```go
import "github.com/alessiosavi/GoGPUtils/cryptoutil"
```

### Functions

| Function | Description |
|----------|-------------|
| `Encrypt(plaintext, key)` | Encrypts data using AES-GCM |
| `Decrypt(ciphertext, key)` | Decrypts AES-GCM encrypted data |
| `EncryptString(plaintext, key)` | Encrypts string, returns base64 |
| `DecryptString(ciphertext, key)` | Decrypts base64 encoded ciphertext |
| `DeriveKey(password, salt)` | Derives key bytes from password and salt |
| `GenerateKey(size)` | Generates a random AES key |

### Example

```go
// Generate a key
key, err := cryptoutil.GenerateKey(32)
if err != nil {
    return err
}

// Encrypt
plaintext := "Hello, World!"
ciphertext, err := cryptoutil.EncryptString(plaintext, key)
if err != nil {
    return err
}

// Decrypt
decrypted, err := cryptoutil.DecryptString(ciphertext, key)
if err != nil {
    return err
}
// decrypted == "Hello, World!"

// Password-based encryption
derivedKey := cryptoutil.DeriveKey("my-password", "my-salt")
_ = derivedKey
```

---

## randutil

Cryptographically secure random generation.

```go
import "github.com/alessiosavi/GoGPUtils/randutil"
```

### Functions

| Function | Description |
|----------|-------------|
| `SecureBytes(n)` | Generates n random bytes |
| `SecureString(length, charset)` | Generates random string |
| `SecureInt(max)` | Generates random int in range |
| `SecureInt64(max)` | Generates random int64 in range |
| `SecureID()` | Generates a random identifier |
| `SecureChoice[T](slice)` | Picks random element from slice |
| `NewGenerator()` | Creates a fast generator |
| `NewGeneratorWithSeed(seed)` | Creates a deterministic generator |

### Example

```go
// Generate random bytes
bytes, err := randutil.SecureBytes(32)

// Generate random string
token, err := randutil.SecureString(32, randutil.AlphaNumeric)

// Generate random integer
n, err := randutil.SecureInt(100)

// Pick random element
colors := []string{"red", "green", "blue"}
color, err := randutil.SecureChoice(colors)

// Generate identifier
id, err := randutil.SecureID()
```

---

## collection

Generic data structures.

```go
import "github.com/alessiosavi/GoGPUtils/collection"
```

### Stack

LIFO (Last In, First Out) data structure.

```go
stack := collection.NewStack[int]()
stack.Push(1)
stack.Push(2)
stack.Push(3)

value, ok := stack.Pop()  // 3, true
value, ok = stack.Peek()  // 2, true (doesn't remove)
stack.IsEmpty()           // false
stack.Len()               // 2
```

### Queue

FIFO (First In, First Out) data structure.

```go
queue := collection.NewQueue[string]()
queue.Enqueue("first")
queue.Enqueue("second")
queue.Enqueue("third")

value, ok := queue.Dequeue()  // "first", true
value, ok = queue.Peek()      // "second", true
queue.Len()                   // 2
```

### Set

Unordered collection of unique elements.

```go
set := collection.NewSet[int]()
set.Add(1)
set.Add(2)
set.Add(2)  // No effect, already exists

set.Contains(1)  // true
set.Contains(3)  // false
set.Len()        // 2
set.Remove(1)
values := set.Values()  // [2]

// Set operations
setA := collection.NewSetFrom([]int{1, 2, 3})
setB := collection.NewSetFrom([]int{2, 3, 4})

setA.Union(setB)        // {1, 2, 3, 4}
setA.Intersection(setB) // {2, 3}
setA.Difference(setB)   // {1}
```

### Binary Search Tree (BST)

Binary search tree for ordered data.

```go
bst := collection.NewBST[int]()
bst.Insert(5)
bst.Insert(3)
bst.Insert(7)
bst.Insert(1)

bst.Contains(3)   // true
min, _ := bst.Min()
max, _ := bst.Max()
bst.InOrder()     // [1, 3, 5, 7]
bst.Remove(3)
```

---

## AWS Utilities

For AWS SDK v2 helpers (S3, DynamoDB, SQS, SSM, Secrets Manager, Lambda), see the dedicated [AWS README](./aws/README.md).

---

## Testing

Run all tests:

```bash
go test ./...
```

Run with race detector:

```bash
go test -race ./...
```

Run benchmarks:

```bash
go test -bench=. ./...
```

---

## Contributing

Contributions are welcome! Please ensure:

1. All new code has tests
2. Tests pass with race detector enabled
3. Code follows Go conventions (`gofmt`, `golint`)
4. Public APIs are documented

---

## License

MIT License - see [LICENSE](LICENSE) for details.
