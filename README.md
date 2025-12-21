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
| `Sort[T](slice)` | Returns sorted slice (ordered types) |
| `SortBy[T](slice, less)` | Returns sorted slice with custom comparator |
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
| `Truncate(s, maxLen)` | Truncates to max length |
| `TruncateWithSuffix(s, maxLen, suffix)` | Truncates with custom suffix |
| `PadLeft(s, length, pad)` | Pads string on the left |
| `PadRight(s, length, pad)` | Pads string on the right |
| `IsEmpty(s)` | Checks if string is empty |
| `IsBlank(s)` | Checks if string is empty or whitespace |
| `ToSnakeCase(s)` | Converts to snake_case |
| `ToCamelCase(s)` | Converts to camelCase |
| `ToPascalCase(s)` | Converts to PascalCase |
| `ToKebabCase(s)` | Converts to kebab-case |
| `CountWords(s)` | Counts words in string |
| `RemoveDuplicateSpaces(s)` | Collapses multiple spaces |

### Similarity Algorithms

| Function | Description |
|----------|-------------|
| `LevenshteinDistance(a, b)` | Edit distance between strings |
| `JaroWinklerSimilarity(a, b)` | Similarity score (0.0-1.0) |
| `DiceSimilarity(a, b)` | Dice coefficient similarity |
| `HammingDistance(a, b)` | Bit-level distance (same length strings) |

### Example

```go
// Case conversions
stringutil.ToSnakeCase("HelloWorld")     // "hello_world"
stringutil.ToCamelCase("hello_world")    // "helloWorld"
stringutil.ToPascalCase("hello_world")   // "HelloWorld"
stringutil.ToKebabCase("HelloWorld")     // "hello-world"

// Padding
stringutil.PadLeft("42", 5, '0')   // "00042"
stringutil.PadRight("Go", 5, '-')  // "Go---"

// Similarity
stringutil.LevenshteinDistance("kitten", "sitting")  // 3
stringutil.JaroWinklerSimilarity("hello", "hallo")   // ~0.88
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
| `ReadFile(path)` | Reads entire file as bytes |
| `ReadFileString(path)` | Reads entire file as string |
| `WriteFile(path, data, perm)` | Writes data to file |
| `AppendFile(path, data)` | Appends data to file |
| `CopyFile(src, dst)` | Copies a file |
| `MoveFile(src, dst)` | Moves a file |
| `DeleteFile(path)` | Deletes a file |
| `ListFiles(dir)` | Lists files in directory |
| `ListFilesRecursive(dir)` | Lists files recursively |
| `EnsureDir(path)` | Creates directory if not exists |
| `FileSize(path)` | Returns file size |
| `FileModTime(path)` | Returns modification time |

### Example

```go
// Check existence
if fileutil.Exists("/path/to/file") {
    // File exists
}

// Read file
content, err := fileutil.ReadFileString("/path/to/file.txt")
if err != nil {
    return err
}

// Write file
err := fileutil.WriteFile("/path/to/file.txt", []byte("Hello"), 0644)

// List files
files, err := fileutil.ListFiles("/path/to/dir")
for _, f := range files {
    fmt.Println(f.Name)
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
| `DeriveKey(password, salt)` | Derives key from password using Argon2 |
| `GenerateKey()` | Generates random 256-bit key |

### Example

```go
// Generate a key
key, err := cryptoutil.GenerateKey()
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
salt := make([]byte, 16)
rand.Read(salt)
key, err := cryptoutil.DeriveKey("my-password", salt)
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
| `SecureString(n)` | Generates random alphanumeric string |
| `SecureHex(n)` | Generates random hex string |
| `SecureBase64(n)` | Generates random base64 string |
| `SecureInt(min, max)` | Generates random int in range |
| `SecureFloat64()` | Generates random float64 [0.0, 1.0) |
| `SecureChoice[T](slice)` | Picks random element from slice |
| `SecureShuffle[T](slice)` | Shuffles slice in-place |
| `SecureUUID()` | Generates UUID v4 |

### Example

```go
// Generate random bytes
bytes, err := randutil.SecureBytes(32)

// Generate random string
token, err := randutil.SecureString(32)
// e.g., "a7Bk9mNpQr2sT5uV8wXy1zA3bC6dE4fG"

// Generate random integer
n, err := randutil.SecureInt(1, 100)

// Pick random element
colors := []string{"red", "green", "blue"}
color, err := randutil.SecureChoice(colors)

// Generate UUID
uuid, err := randutil.SecureUUID()
// e.g., "550e8400-e29b-41d4-a716-446655440000"
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

Self-balancing binary search tree for ordered data.

```go
bst := collection.NewBST[int]()
bst.Insert(5)
bst.Insert(3)
bst.Insert(7)
bst.Insert(1)

bst.Contains(3)   // true
bst.Min()         // 1, true
bst.Max()         // 7, true
bst.InOrder()     // [1, 3, 5, 7]
bst.Delete(3)
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
