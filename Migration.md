# GoGPUtils v2

A collection of well-tested, idiomatic Go utilities for common programming tasks.

## Installation

```bash
go get github.com/alessiosavi/GoGPUtils/v2
```

Requires Go 1.24 or later.

## Packages

### sliceutil - Generic Slice Operations

```go
import "github.com/alessiosavi/GoGPUtils/v2/sliceutil"

// Filter
evens := sliceutil.Filter([]int{1, 2, 3, 4, 5}, func(n int) bool { return n%2 == 0 })
// [2, 4]

// Map
doubled := sliceutil.Map([]int{1, 2, 3}, func(n int) int { return n * 2 })
// [2, 4, 6]

// Reduce
sum := sliceutil.Reduce([]int{1, 2, 3, 4}, 0, func(acc, n int) int { return acc + n })
// 10

// Chunk
chunks := sliceutil.Chunk([]int{1, 2, 3, 4, 5}, 2)
// [[1, 2], [3, 4], [5]]

// Set operations
common := sliceutil.Intersect([]int{1, 2, 3}, []int{2, 3, 4})  // [2, 3]
unique := sliceutil.Unique([]int{1, 2, 2, 3, 1})               // [1, 2, 3]

// Group by
groups := sliceutil.GroupBy(nums, func(n int) string {
    if n%2 == 0 { return "even" }
    return "odd"
})
```

### stringutil - String Manipulation

```go
import "github.com/alessiosavi/GoGPUtils/v2/stringutil"

// Search
indices := stringutil.AllIndexes("banana", "an")  // [1, 3]
ok := stringutil.HasAnyPrefix(url, "http://", "https://")
ok := stringutil.ContainsAll(text, "error", "line")

// Transform
reversed := stringutil.Reverse("hello")           // "olleh"
truncated := stringutil.Truncate(text, 100, "...")
padded := stringutil.PadLeft("42", 5, '0')        // "00042"

// Case conversion
stringutil.SnakeCase("HelloWorld")    // "hello_world"
stringutil.CamelCase("hello_world")   // "helloWorld"
stringutil.PascalCase("hello_world")  // "HelloWorld"

// Similarity algorithms
distance := stringutil.LevenshteinDistance("kitten", "sitting")  // 3
score := stringutil.JaroWinklerSimilarity("martha", "marhta", 0.1)
coefficient := stringutil.DiceCoefficient("night", "nacht")
```

### mathutil - Mathematical Operations

```go
import "github.com/alessiosavi/GoGPUtils/v2/mathutil"

// Aggregations
sum := mathutil.Sum([]int{1, 2, 3, 4, 5})        // 15
avg := mathutil.Average([]float64{1, 2, 3, 4})  // 2.5

// Statistics
median := mathutil.Median([]int{1, 2, 3, 4, 5})  // 3.0
mode := mathutil.Mode([]int{1, 2, 2, 3})         // [2]
stddev := mathutil.StdDev(data)

// Number theory
if mathutil.IsPrime(17) { ... }
primes := mathutil.Primes(100)  // [2, 3, 5, 7, 11, ...]
gcd := mathutil.GCD(12, 18)     // 6

// Vector operations
dot := mathutil.DotProduct(a, b)
cosine := mathutil.CosineSimilarity(vec1, vec2)
distance := mathutil.EuclideanDistance(point1, point2)

// Matrix operations
result, err := mathutil.MatrixMultiply(a, b)
transposed := mathutil.MatrixTranspose(matrix)
```

### fileutil - File System Operations

```go
import "github.com/alessiosavi/GoGPUtils/v2/fileutil"

ctx := context.Background()

// Reading
lines, err := fileutil.ReadLines(ctx, "file.txt")
content, err := fileutil.ReadString(ctx, "file.txt")
data, err := fileutil.ReadBytes(ctx, "file.txt")

// File info
if fileutil.Exists(path) { ... }
if fileutil.IsDir(path) { ... }
size, err := fileutil.Size(path)

// Directory operations
files, err := fileutil.List(ctx, dir, fileutil.Recursive|fileutil.FilesOnly)
matches, err := fileutil.Find(ctx, dir, "*.go")
err := fileutil.EnsureDir(path, 0755)

// Line terminator handling
terminator := fileutil.DetectLineTerminator(data)
normalized := fileutil.NormalizeLineTerminators(data, fileutil.LF)
```

### cryptoutil - Secure Encryption

```go
import "github.com/alessiosavi/GoGPUtils/v2/cryptoutil"

// Generate a key
key, err := cryptoutil.GenerateKey(32)  // AES-256

// Or derive from password
key := cryptoutil.DeriveKey("password", "salt")

// Encrypt/Decrypt (AES-GCM)
ciphertext, err := cryptoutil.Encrypt(plaintext, key)
plaintext, err := cryptoutil.Decrypt(ciphertext, key)

// String convenience
encrypted, err := cryptoutil.EncryptString("secret", key)
decrypted, err := cryptoutil.DecryptString(encrypted, key)

// Hashing
hash := cryptoutil.HashString(data)  // SHA-256 hex
ok := cryptoutil.CompareHash(data, expectedHash)

// Random bytes
randomBytes, err := cryptoutil.RandomBytes(32)
```

### collection - Generic Data Structures

```go
import "github.com/alessiosavi/GoGPUtils/v2/collection"

// Stack
stack := collection.NewStack[int]()
stack.Push(1)
stack.Push(2)
val, ok := stack.Pop()  // 2, true

// Queue
queue := collection.NewQueue[string]()
queue.Enqueue("first")
queue.Enqueue("second")
val, ok := queue.Dequeue()  // "first", true

// Set
set := collection.NewSetFrom([]int{1, 2, 3})
set.Add(4)
set.Contains(2)  // true
union := set1.Union(set2)
intersection := set1.Intersection(set2)

// Binary Search Tree
tree := collection.NewBST[int]()
tree.Insert(5, 3, 7, 1, 4)
tree.Contains(3)        // true
sorted := tree.InOrder() // [1, 3, 4, 5, 7]
min, _ := tree.Min()     // 1
```

### randutil - Random Generation

```go
import "github.com/alessiosavi/GoGPUtils/v2/randutil"

// Secure random (crypto/rand)
bytes, err := randutil.SecureBytes(32)
token, err := randutil.SecureString(16, randutil.AlphaNumeric)
id, err := randutil.SecureID()  // 32-char hex

// Fast random (math/rand/v2)
rng := randutil.NewGenerator()
n := rng.Int(100)               // [0, 100)
n := rng.IntRange(10, 20)       // [10, 20]
f := rng.Float64()              // [0.0, 1.0)
choice := rng.Choice(items)     // Random element
rng.ShuffleSlice(items)         // Shuffle in place
sample := rng.Sample(items, 5)  // 5 random without replacement

// Deterministic (for tests)
rng := randutil.NewGeneratorWithSeed(42)

// Sequences
nums := randutil.Range(0, 10)       // [0, 1, 2, ..., 9]
odds := randutil.RangeStep(1, 10, 2) // [1, 3, 5, 7, 9]
```

## Design Principles

1. **Errors over panics** - All functions return errors instead of panicking
2. **No global state** - All state is explicit; no singletons or hidden dependencies
3. **Context-aware** - I/O operations accept `context.Context` for cancellation
4. **Generic where useful** - Uses generics to reduce duplication without over-abstraction
5. **Zero external dependencies** - Core library has no external dependencies
6. **Nil-safe** - Functions handle nil inputs gracefully

## Testing

Run all tests:

```bash
cd v2
go test ./...
```

Run with coverage:

```bash
go test -cover ./...
```

Run benchmarks:

```bash
go test -bench=. ./...
```

## Migration from v1

### Error Handling

```go
// v1 (panics)
client := awsutils.New()  // panics on error

// v2 (returns error)
key, err := cryptoutil.GenerateKey(32)
if err != nil {
    return err
}
```

### Context Support

```go
// v1 (uses context.Background() internally)
sqsClient.GetMessage(queueName)

// v2 (accepts context)
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
lines, err := fileutil.ReadLines(ctx, "file.txt")
```

### Generic Collections

```go
// v1 (float64 only)
stack := datastructure.NewStack()
stack.Push(1.5)
val := stack.Pop().(float64)

// v2 (generic)
stack := collection.NewStack[float64]()
stack.Push(1.5)
val, ok := stack.Pop()  // val is float64
```

### Secure Encryption

```go
// v1 (AES-ECB - insecure!)
crypt.EncryptAES(plaintext, key)

// v2 (AES-GCM - authenticated encryption)
ciphertext, err := cryptoutil.Encrypt(plaintext, key)
```

## License

MIT License - see LICENSE file for details.
