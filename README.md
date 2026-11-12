# golang

My Golang studies! A practical, idiomatic path through Go for developers coming
from TypeScript, JavaScript, Java, or C. It skips "what is a loop" and focuses on
what is *different* in Go. Everything is written in English, including the
exercises, to practice technical English along the way.

## Table of Contents

- [Stage 0: Installing Go and the Toolchain](#stage-0-installing-go-and-the-toolchain)
- [How to Use This Repository](#how-to-use-this-repository)
- [Stage 1: Syntax, Types & Pointers](#stage-1-syntax-types--pointers)
- [Stage 2: Data Structures in Go](#stage-2-data-structures-in-go)
- [Stage 3: Implicit Interfaces & Error Handling](#stage-3-implicit-interfaces--error-handling)
- [Stage 4: Native Concurrency](#stage-4-native-concurrency)
- [Stage 5: Practical Capstone Project](#stage-5-practical-capstone-project)
- [Stage 6: I/O, Networking & the Standard Library](#stage-6-io-networking--the-standard-library)
- [Stage 7: Production Patterns](#stage-7-production-patterns)
- [Further Reading](#further-reading)

---

## Stage 0: Installing Go and the Toolchain

These are manual installation steps; nothing in this repository downloads or
installs software for you.

### 1. Install Go (macOS)

Pick **one** option.

**Option A — Official installer (no extra tools needed)**

1. Open <https://go.dev/dl/>.
2. Download the macOS package that matches your CPU:
   - Apple Silicon (M1/M2/M3/...) → `darwin-arm64` `.pkg`
   - Intel → `darwin-amd64` `.pkg`
   - Not sure? Run `uname -m` in the terminal: `arm64` or `x86_64`.
3. Open the `.pkg` and follow the wizard. It installs Go to `/usr/local/go` and
   adds `/usr/local/go/bin` to your `PATH` through `/etc/paths.d/go`.
4. **Close and reopen your terminal** so the new `PATH` is loaded.

**Option B — Homebrew**

```bash
brew install go
```

### 2. Verify the installation

```bash
go version        # e.g. go version go1.2x.y darwin/arm64
go env GOROOT     # where Go itself lives
go env GOPATH     # where downloaded modules and `go install` binaries go (default: ~/go)
```

If you get `command not found: go`, your `PATH` is missing Go's `bin` directory.
For the official installer, add this to `~/.zshrc` and reopen the terminal:

```bash
export PATH="$PATH:/usr/local/go/bin"
```

To run binaries installed with `go install`, also add:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

### 3. Editor setup (optional but recommended)

- **VS Code**: install the extension "Go" (by the Go Team at Google). When it
  asks to install tools such as `gopls`, accept; that is the language server.
- **GoLand** (JetBrains) works out of the box.

### 4. The commands you will use every day

| Command | What it does | Closest analogy |
| --- | --- | --- |
| `go mod init <name>` | Creates `go.mod` (already done here: `go-studies`) | `npm init` / `pom.xml` |
| `go run main.go` / `go run .` | Compiles to a temp binary and runs it | `ts-node` / `java Main.java` |
| `go build ./...` | Compiles everything; errors if code doesn't compile | `tsc` / `javac` |
| `go build -o app .` | Produces a single static binary named `app` | `gcc -o` |
| `go test ./...` | Runs every `_test.go` in every package | `npm test` / JUnit |
| `go test -v ./exercises/stage1` | Verbose output for one package | |
| `go test -run TestSwap ./exercises/stage1` | Runs only tests whose name matches the regex | `-t` in Jest |
| `go test -race ./...` | Runs tests with the **race detector** (vital in Stage 4) | no direct equivalent |
| `go test -count=1 ./...` | Disables result caching | |
| `go fmt ./...` | Formats code (`gofmt`); there is exactly one style | Prettier, but not optional |
| `go vet ./...` | Static analysis for suspicious code | a light linter |
| `go doc strings.Fields` | Shows documentation in the terminal | `man` |
| `go mod tidy` | Syncs `go.mod` with actual imports | `npm install` / prune |

Useful facts:

- A **package** = one directory. All `.go` files in a directory share a package name.
- A **module** = a tree of packages with one `go.mod`. Import paths start with the module name: `go-studies/exercises/stage1`.
- `package main` + `func main()` makes an executable. Everything else is a library.
- Test files end in `_test.go`; test functions are `func TestXxx(t *testing.T)`.
- Unused imports and unused local variables are **compile errors**.

### 5. Your first program (to confirm it works)

Create a scratch file **outside** the exercise folders, for example `scratch/main.go`:

```go
package main

import "fmt"

func main() {
	fmt.Println("Hello, Go!")
}
```

```bash
go run ./scratch
```

---

## How to Use This Repository

```text
go.mod                 module go-studies
README.md              this guide
exercises/
  stage1/ ... stage4/  one package per stage: starter code + tests
  stage5/
    workerpool/        capstone option A
    api/               capstone option B
  stage6/              I/O, networking, encoding, text, filesystem
  stage7/              production patterns: generics, middleware, resilience
```

Workflow for every stage:

1. Read the stage section below.
2. Open the starter `.go` files in `exercises/stageN` and read each `// TODO:`.
3. Run the tests and watch them fail: `go test ./exercises/stage1`
4. Implement, re-run, repeat until green.
5. Run `go fmt ./...` and `go vet ./...` before moving on.

The tests **fail** until you implement each exercise. Run them all with
`go test ./...` once you have finished.

---

## Stage 1: Syntax, Types & Pointers

### Minimalism: what Go removes

| Feature | Java / TypeScript | Go |
| --- | --- | --- |
| Classes & inheritance | `class A extends B` | None. Structs + methods + embedding |
| Exceptions | `try/catch/throw` | Errors are values (Stage 3) |
| Overloading | yes | no: one function per name |
| Ternary `?:` | yes | no: use `if` |
| `while` | yes | only `for` |
| Semicolons | yes | inserted by the lexer; never write them |
| Visibility keywords | `public/private` | Capitalized name = exported, lowercase = package-private |
| Generics | yes | yes, since 1.18 (simpler, used sparingly) |

### Declarations and zero values

```go
var count int          // zero value: 0. There is no "undefined" and no uninitialized memory.
name := "gopher"       // short declaration, only inside functions
const maxRetries = 3   // untyped constant

var (
	s string // ""
	b bool   // false
	p *int   // nil
)
```

Every type has a usable **zero value**. Idiomatic Go designs types so the zero
value is ready to use (`sync.Mutex`, `bytes.Buffer`).

### Types are strict

```go
var i int = 10
var f float64 = float64(i) // explicit conversion, no implicit widening like in Java/C
// var g float64 = i       // compile error
```

`int` is platform-sized (64-bit on modern machines). `byte` = `uint8`,
`rune` = `int32` (a Unicode code point). Strings are immutable UTF-8 byte
sequences: `len("é") == 2`, but `[]rune("é")` has length 1.

### Functions: multiple returns, variadics, closures

```go
func divmod(a, b int) (q, r int) { // named results
	return a / b, a % b
}

func sum(nums ...int) (total int) { // nums is []int inside
	for _, n := range nums {
		total += n
	}
	return
}

func counter() func() int { // closure, like a JS closure
	n := 0
	return func() int { n++; return n }
}
```

Functions are first-class values, exactly like in TypeScript.

### `defer`

```go
f, err := os.Open(path)
if err != nil {
	return err
}
defer f.Close() // runs when the surrounding function returns, LIFO order
```

Think `finally`, but attached right next to the acquisition. Deferred call
arguments are evaluated **immediately**; the call itself runs later.

### Control flow

```go
for i := 0; i < 3; i++ {}        // classic
for cond {}                      // while
for {}                           // infinite
for i, v := range []string{"a"} {} // index + value
for i := range 3 {}              // Go 1.22+: range over an integer

switch x := v.(type) {           // type switch
case int:    _ = x + 1
case string: _ = len(x)
default:
}
```

`switch` cases **don't fall through** (no `break` needed; the opposite of C/Java).
Use `fallthrough` explicitly if you really want it.

### Value vs. reference semantics

Go is **always pass-by-value**. The question is *what* the value is:

| Type | What gets copied | Behaves like |
| --- | --- | --- |
| `int`, `bool`, `string` | the data | primitives |
| array `[3]int` | the whole array | a C struct, not a Java array! |
| `struct` | all fields | a C struct, not a Java object |
| pointer `*T` | the address | a C pointer / a Java reference |
| slice, map, channel, func | a small header pointing at shared data | a Java reference, but copied by value |

```go
type Point struct{ X, Y int }

func (p Point) Add(q Point) Point { return Point{p.X + q.X, p.Y + q.Y} } // value receiver: p is a copy
func (p *Point) Scale(k int)      { p.X *= k; p.Y *= k }               // pointer receiver: mutates caller
```

Rule of thumb: use a pointer receiver if the method mutates, if the struct is
large, or if *any* method of the type already uses one (stay consistent).

### Safe pointers compared to C

```go
x := 42
p := &x      // address-of, like C
*p = 43      // dereference, like C
// p++       // compile error: NO pointer arithmetic
```

- No pointer arithmetic, no manual `malloc/free`. A garbage collector frees memory.
- Returning the address of a local variable is **safe**: the compiler's escape
  analysis moves it to the heap.

```go
func newPoint() *Point {
	p := Point{1, 2}
	return &p // fine in Go; undefined behavior in C
}
```

- `new(T)` allocates a zeroed `T` and returns `*T`; `&T{...}` is more common.
- Dereferencing `nil` panics (a defined runtime error, not a segfault/UB).
- `unsafe` exists but you shouldn't need it.

### Enumerations with `iota`

```go
type Weekday int

const (
	Sunday Weekday = iota // 0
	Monday                // 1
	Tuesday               // 2
)

func (d Weekday) String() string { /* implements fmt.Stringer */ return "" }
```

**Exercises:** `exercises/stage1` — `go test ./exercises/stage1`

---

## Stage 2: Data Structures in Go

### Arrays vs. slices

An **array** has its length in its type and is a value:

```go
a := [3]int{1, 2, 3}
b := a        // full copy
b[0] = 99     // a is unchanged
// [3]int and [4]int are different, incompatible types
```

A **slice** is a tiny struct (a *header*): `{ptr, len, cap}` pointing into a
**backing array**.

```text
s := []int{10, 20, 30, 40, 50}
t := s[1:3]

 s: ptr ─┐ len=5 cap=5
         ▼
 array: [10][20][30][40][50]
              ▲
 t:      ptr ─┘ len=2 cap=4   (cap counts to the end of the backing array)
```

- `len(s)` = elements you can see; `cap(s)` = room before reallocation.
- Slicing **shares** memory: writing through `t` changes `s`.
- `append` writes in place if `len < cap`; otherwise it allocates a bigger
  array (roughly 2x for small slices), copies, and returns a new header.
  **Always use the return value:** `s = append(s, x)`.

```go
t = append(t, 100) // cap(t) is 4, so this OVERWRITES s[3]!
```

The classic trap, and the fix:

```go
clone := slices.Clone(s)      // independent backing array (Go 1.21+)
clone2 := append([]int(nil), s...)
clean := s[1:3:3]             // full slice expression: cap capped at 2, next append reallocates
```

Other idioms:

```go
s := make([]int, 0, 10)            // len 0, cap 10: preallocate when you know the size
s = append(s[:i], s[i+1:]...)      // delete index i (keeps order)
var nilSlice []int                 // nil slice: len 0, cap 0, append works fine
```

A nil slice and an empty slice both have `len == 0`; prefer checking `len(s) == 0`.
Compare with Java's `ArrayList` (growth is explicit here) and JS arrays (always shared by reference).

### Maps

```go
m := map[string]int{"a": 1}
m["b"]++                      // missing key reads as the zero value (0)
v, ok := m["zzz"]             // "comma ok": distinguishes missing from zero
delete(m, "a")                // no-op if absent
for k, v := range m {}        // order is intentionally RANDOM
```

- A nil map can be read but **writing panics**. Create with `make` or a literal.
- Not safe for concurrent writes (use a mutex or `sync.Map`, see Stage 4).
- Keys must be comparable (`==`): no slices, maps, or funcs as keys; structs and arrays are fine.
- Like `HashMap` / JS `Map`, but with value semantics for struct keys and random iteration order.

### Structs

```go
type User struct {
	Name  string
	Email string `json:"email"` // struct tag, used by encoding/json (like a Java annotation)
}

u := User{Name: "Ana"}   // always prefer named fields
p := &u                  // p.Name works: automatic dereference
```

There are no constructors; by convention, write `func NewUser(...) *User` when
you need to enforce invariants. Structs are comparable with `==` if all fields are.

### Composition over inheritance (embedding)

```go
type Person struct{ Name string }
func (p Person) Greet() string { return "Hi, I'm " + p.Name }

type Employee struct {
	Person        // embedded: fields and methods are PROMOTED
	Role string
}

e := Employee{Person: Person{Name: "Ana"}, Role: "Dev"}
e.Greet()   // works: promoted from Person
e.Name      // works too
```

This is *not* inheritance:

- There is no subtyping: an `Employee` is not a `Person`; you can't pass it where a `Person` is expected.
- No virtual dispatch: `Person.Greet` never calls an `Employee` override. The receiver is always the embedded value.
- It's sugar for a field plus auto-generated forwarding methods. In Java you would
  write a `Person person` field and delegate manually.
- Embedding an *interface* in a struct is a common trick for partial implementations and decorators.

**Exercises:** `exercises/stage2` — `go test ./exercises/stage2`

---

## Stage 3: Implicit Interfaces & Error Handling

### Implicit interfaces: static duck typing

```go
type Shape interface {
	Area() float64
}

type Circle struct{ R float64 }
func (c Circle) Area() float64 { return math.Pi * c.R * c.R }

var s Shape = Circle{R: 2} // OK: Circle has Area(), so it satisfies Shape
```

There is no `implements`. If the method set matches, the type satisfies the
interface; checked **at compile time**, unlike TypeScript's structural typing
being erased at runtime and unlike JS's runtime-only duck typing. It's close
to TypeScript's structural types, but for *methods* and with real runtime values.

Idioms:

- **Define interfaces where they are used (the consumer), not where implemented.**
  Java does the opposite.
- **Keep them tiny.** `io.Reader` has one method, `Read([]byte) (int, error)`;
  files, network connections, strings, and gzip streams all satisfy it.
- "Accept interfaces, return structs."
- Compile-time assertion that a type satisfies an interface: `var _ Shape = Circle{}`.
- `any` (alias of `interface{}`) holds anything; recover the type with a type
  assertion `v, ok := x.(int)` or a type switch.
- **Nil trap:** an interface value is a `(type, value)` pair. An interface holding
  a nil `*MyErr` is **not** `== nil`. Return a literal `nil` for "no error".

### Errors are values

```go
type error interface {
	Error() string
}
```

There are no exceptions for ordinary failures. Functions return an `error` as
the last result, and the caller decides what to do:

```go
f, err := os.Open("data.txt")
if err != nil {
	return fmt.Errorf("open data: %w", err) // add context, keep the cause
}
defer f.Close()
```

| | Java / TS | Go |
| --- | --- | --- |
| Failure path | invisible: any call may throw | visible in every signature and call site |
| Handling | `try { } catch (e) { }` | `if err != nil { }` right after the call |
| Cleanup | `finally` | `defer` |
| Error type | class hierarchy | any type with `Error() string` |
| Cause chain | `new X(msg, cause)` | `%w` + `errors.Unwrap` |

`if err != nil` is verbose by design: error paths are explicit and local. Handle
the error once: either return it (with context) or log it, not both.

### Sentinel errors, wrapping, and custom types

```go
var ErrNotFound = errors.New("not found") // sentinel

func find(id int) error {
	return fmt.Errorf("find %d: %w", id, ErrNotFound) // wrap with %w
}

err := find(7)
errors.Is(err, ErrNotFound) // true: walks the wrap chain (never compare with ==)

type ValidationError struct{ Field string }
func (e *ValidationError) Error() string { return "invalid " + e.Field }

var ve *ValidationError
if errors.As(err, &ve) { // extracts a specific TYPE from the chain
	fmt.Println(ve.Field)
}

errors.Join(err1, err2) // combine several errors into one
```

### `panic` and `recover`

`panic` is for programmer errors and impossible states (nil map write, index out
of range), not for control flow. Libraries should return errors. `recover()` only
works inside a `defer` and is mostly used at boundaries (an HTTP server recovering
per request so one bug doesn't kill the process).

**Exercises:** `exercises/stage3` — `go test ./exercises/stage3`

---

## Stage 4: Native Concurrency

> "Don't communicate by sharing memory; share memory by communicating."

### Goroutines

```go
go doWork() // starts a goroutine: a lightweight, runtime-scheduled coroutine
```

- Start with ~2 KB of stack (grows dynamically); you can run hundreds of thousands.
- Multiplexed over OS threads by Go's scheduler (M:N). Not 1:1 with threads like Java's
  platform threads, and truly parallel unlike Node's single event loop.
- When `main` returns, **all goroutines die**. You must wait for them explicitly.
- No `async/await` and no "colored functions": any function can be run concurrently.

### Channels

```go
ch := make(chan int)      // unbuffered: send blocks until a receiver is ready (rendezvous)
buf := make(chan int, 8)  // buffered: send blocks only when full

go func() { ch <- 42 }()
v := <-ch

close(ch)                 // only the SENDER closes; signals "no more values"
for v := range ch {}      // loops until closed and drained
v, ok := <-ch             // ok == false when closed and empty
```

Directional types document intent: `chan<- int` (send-only), `<-chan int` (receive-only).

Rules to memorize:

| Operation | nil channel | closed channel |
| --- | --- | --- |
| send | blocks forever | **panics** |
| receive | blocks forever | returns zero value immediately (`ok == false`) |
| close | **panics** | **panics** |

### `select`

```go
select {
case v := <-ch1:
	fmt.Println("got", v)
case ch2 <- 7:
	fmt.Println("sent")
case <-time.After(time.Second):
	fmt.Println("timeout")
case <-ctx.Done():
	return ctx.Err()
default: // makes the select non-blocking
}
```

`select` waits on several channel operations at once; if several are ready,
one is chosen at random. This is the idiomatic way to do timeouts and cancellation.

### `sync.WaitGroup` and `Mutex`

```go
var wg sync.WaitGroup
for i := 0; i < 3; i++ {
	wg.Add(1)             // BEFORE starting the goroutine
	go func() {
		defer wg.Done()
		work(i)           // Go 1.22+: each iteration has its own i
	}()
}
wg.Wait()
```

```go
type SafeCounter struct {
	mu sync.Mutex
	n  int
}

func (c *SafeCounter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.n++
}
```

Never copy a `sync.Mutex` or `WaitGroup` after first use: pass pointers.
Channels are for ownership transfer and coordination; mutexes are for protecting
shared state. Both are idiomatic; choose the simpler one.
Also see `sync/atomic` and `sync.Once`.

### Race conditions

A **data race** is two goroutines accessing the same memory concurrently with at
least one write and no synchronization. Results are undefined (like C), not
merely "a bit wrong".

```go
// BUG: unsynchronized counter
var n int
for i := 0; i < 1000; i++ {
	go func() { n++ }()
}
```

Find races with the built-in detector:

```bash
go test -race ./exercises/stage4
go run -race main.go
```

Prevent them by: not sharing (confine data to one goroutine), channels, mutexes, or atomics.

### Patterns

- **Pipeline:** stages connected by channels; each stage `range`s its input and closes its output.
- **Fan-out / fan-in:** several workers read one channel; a merge function combines outputs.
- **Worker pool:** N goroutines consuming a jobs channel (Stage 5, option A).
- **Cancellation:** pass a `context.Context` as the first parameter; select on `ctx.Done()`.
- **Avoid goroutine leaks:** every goroutine needs a guaranteed way to exit (buffered result channel, closed channel, or context).

**Exercises:** `exercises/stage4` — `go test -race ./exercises/stage4`

---

## Stage 5: Practical Capstone Project

Two independent options; do one or both. Both use only the standard library.

### Option A: Concurrent worker pool

`exercises/stage5/workerpool`: process a list of jobs with at most N concurrent
workers, return results in input order, and collect per-job errors without
stopping the others.

Skeleton:

```go
jobs := make(chan int)
var wg sync.WaitGroup
results := make([]Result, len(input)) // each worker writes its own slots: no race

for w := 0; w < workers; w++ {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range jobs {
			results[i] = run(input[i])
		}
	}()
}
for i := range input {
	jobs <- i
}
close(jobs)
wg.Wait()
```

```bash
go test -race ./exercises/stage5/workerpool
```

### Option B: REST API with `net/http`

`exercises/stage5/api`: an in-memory todo API with JSON, status codes, and a
concurrency-safe store. Go 1.22+ `http.ServeMux` supports methods and path
wildcards, so no router library is needed.

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /todos/{id}", func(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"id": id})
})
// http.ListenAndServe(":8080", mux)
```

Handlers run on one goroutine **per request**, so any shared state needs a
mutex. Tests use `net/http/httptest`, so no real port is opened.
Compare with Express (`app.get('/todos/:id')`) and Spring (`@GetMapping`).

```bash
go test -race ./exercises/stage5/api
```

---

## Stage 6: I/O, Networking & the Standard Library

Most real Go programs move bytes between places: files, sockets, HTTP bodies.
This stage covers the standard-library pieces behind tools such as web crawlers,
search indexes, file-transfer utilities, and network services.

### `io.Reader` / `io.Writer`: the core abstraction

```go
type Reader interface{ Read(p []byte) (n int, err error) }
type Writer interface{ Write(p []byte) (n int, err error) }
```

Files, TCP connections, HTTP bodies, gzip streams, and hashes all implement them,
so they compose like Unix pipes (or Java streams, but far smaller):

```go
io.Copy(dst, src)                    // streams in chunks; constant memory
io.ReadAll(r)                        // loads everything: only for small inputs
io.LimitReader(r, 1<<20)             // cap how much you read from untrusted input
io.TeeReader(r, hash)                // read and hash at the same time
io.ReadFull(r, buf)                  // exactly len(buf) bytes or an error
```

Contract details that trip people up:

- `Read` may return `n > 0` **and** an error (even `io.EOF`) in the same call: process the bytes first.
- `io.EOF` means a clean end; `io.ErrUnexpectedEOF` means the stream ended mid-read.
- A reader that wraps another reader is the **decorator pattern** (`bufio.NewReader(conn)`).

### Files and paths

```go
f, err := os.Open(path)           // read-only; os.Create truncates/creates
if err != nil { return err }
defer f.Close()

data, _ := os.ReadFile(path)      // small files only
filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error { return nil })
```

Use `path/filepath` for OS paths and `path` for URLs. `t.TempDir()` gives tests a
sandbox that is deleted automatically.

### TCP with `net`

TCP is a **byte stream**: there are no message boundaries. Two writes may arrive
as one read, and one write may arrive as two reads. The standard fix is **framing**,
such as a length prefix:

```text
+----------------+------------------+
| 4-byte length  | payload (length) |   big-endian, then raw bytes
+----------------+------------------+
```

```go
// server
ln, _ := net.Listen("tcp", ":9000")
for {
	conn, err := ln.Accept()
	if err != nil { return err }
	go handle(conn)               // one goroutine per connection
}

// client
conn, _ := net.DialTimeout("tcp", "host:9000", 5*time.Second)
conn.SetDeadline(time.Now().Add(10 * time.Second)) // ALWAYS set deadlines on network I/O
```

Why it matters:

- Without deadlines a silent peer blocks a goroutine forever (goroutine leak).
- Never trust a length read from the network: cap it before allocating.
- `net.Conn` is an `io.Reader` + `io.Writer`, so `io.Copy(conn, file)` sends a file with constant memory.
- Plain TCP is unencrypted; `crypto/tls` wraps a connection with the same interfaces.
- Peers behind home routers (NAT) cannot accept inbound connections without port
  forwarding, a relay server, or NAT traversal. Networking, not Go, is the hard part.

Compare with Node's `net.Socket` (events/callbacks) and Java's `Socket` + threads:
Go's code looks synchronous and blocking, but the runtime multiplexes it efficiently.

### HTTP as a client

```go
client := &http.Client{Timeout: 10 * time.Second}   // the zero-value default client has NO timeout
req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
req.Header.Set("User-Agent", "my-tool/1.0")
resp, err := client.Do(req)
if err != nil { return err }
defer resp.Body.Close()                              // always; otherwise connections leak
if resp.StatusCode < 200 || resp.StatusCode > 299 { /* a 404 is NOT an err */ }
```

- A non-2xx status is **not** a Go error; you must check `StatusCode` yourself (unlike Axios).
- Reuse one `http.Client`: it pools connections.
- Respect the servers you talk to: set a `User-Agent`, limit concurrency, honor `robots.txt`.
- `httptest.NewServer` spins up a real local server for tests.

### `context`: cancellation and deadlines

```go
ctx, cancel := context.WithTimeout(parent, 2*time.Second)
defer cancel()                                  // always release resources

select {
case <-ctx.Done():
	return ctx.Err()                            // context.Canceled or context.DeadlineExceeded
case v := <-results:
	_ = v
}
```

Convention: `ctx` is the **first parameter**, never stored in a struct. It
propagates cancellation through every layer (HTTP request, DB call, goroutines).
This replaces `AbortController` (JS) and thread interruption (Java).

### JSON (`encoding/json`)

```go
type Message struct {
	Type string `json:"type"`
	Body string `json:"body,omitempty"`   // omitted when empty
}

b, _ := json.Marshal(m)
dec := json.NewDecoder(r)
dec.DisallowUnknownFields()                // strict decoding
err := dec.Decode(&m)
```

Only exported fields are encoded. `json.NewDecoder(r)` streams from any `io.Reader`;
`json.Unmarshal` works on a `[]byte` you already hold in memory. Other encodings,
`encoding/gob` and `encoding/csv`, work the same way.

### Text, URLs, and indexing

- `net/url`: `url.Parse`, `base.ResolveReference(ref)` for relative links, `u.Query()`.
- `regexp`: RE2 syntax (linear time, no backreferences). Compile once at package level with `regexp.MustCompile`.
- The standard library has **no HTML parser**; use `golang.org/x/net/html` (an official, separately versioned module) for real documents.
- `strings`, `unicode`, `sort`, and `slices` cover tokenizing and sorting.
- An **inverted index** (`map[token][]docID`) is the core data structure of a search engine; queries are intersections of posting lists.

**Exercises:** `exercises/stage6` — `go test ./exercises/stage6`

---

## Stage 7: Production Patterns

What shows up in real Go codebases beyond the language itself. Each item has an
exercise except those marked *(read-only)*.

### Generics (Go 1.18+)

```go
func Map[T, U any](s []T, f func(T) U) []U { /* ... */ return nil }

type Number interface{ ~int | ~int64 | ~float64 }   // type constraint
func SumOf[T Number](xs []T) T { var s T; for _, x := range xs { s += x }; return s }

type Set[T comparable] map[T]struct{}
```

Compare with Java generics (no erasure here: the compiler generates code per shape)
and TypeScript generics (checked at compile time only). Idiom: reach for generics
for containers and utilities (`slices`, `maps` packages); keep interfaces for behavior.

### Functional options

```go
func NewServer(addr string, opts ...Option) *Server
srv := NewServer(":8080", WithTimeout(5*time.Second), WithMaxConns(50))
```

A stable, extensible constructor API without overloads or builders. Used by
many popular Go libraries.

### HTTP middleware

```go
type Middleware func(http.Handler) http.Handler
handler := Chain(mux, RequestID(), Logging(logger), Recover())
```

Cross-cutting concerns (logging, auth, panic recovery, request IDs, metrics) wrap
handlers. Same idea as Express `app.use(...)` or Servlet filters. Request-scoped
values travel in `context.Context` using an **unexported key type**.

### Structured logging (`log/slog`, Go 1.21+)

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
logger.Info("request", "method", r.Method, "path", r.URL.Path, "status", 200)
```

Key/value logs that machines can query, instead of formatted strings.

### Resilience: retries, rate limits, timeouts

- **Retry with exponential backoff** (+ jitter) for transient failures; stop on `ctx.Done()`.
- **Token bucket** rate limiter to protect your service or respect someone else's.
- Timeouts everywhere: HTTP clients, servers (`ReadHeaderTimeout`), DB calls, connections.
- **Inject time** (`now func() time.Time`) so time-dependent code is testable without `Sleep`.

### Graceful shutdown

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()
// serve in a goroutine, wait for ctx.Done(), then srv.Shutdown(shutdownCtx)
```

Required for containers and orchestrators (Kubernetes sends SIGTERM): finish
in-flight requests, then exit.

### Error groups

`golang.org/x/sync/errgroup` runs goroutines that can fail: the first error cancels
a shared context and `Wait()` returns it. The exercise rebuilds a minimal version so you
understand `sync.Once`, `WaitGroup`, and `context.WithCancel` working together.

### Testing like a Go developer

```go
func TestSomething(t *testing.T) {
	tests := []struct{ name string; in, want int }{ /* ... */ }
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { /* subtest */ })
	}
}
```

- **Table-driven tests + subtests** are the norm; `t.Helper()`, `t.TempDir()`, `t.Cleanup()` keep them tidy.
- Fakes over mocking frameworks: small interfaces make fakes trivial.
- `httptest` for HTTP; `-race` in CI; `go test -cover`.
- Benchmarks: `go test -bench=. -benchmem ./exercises/stage7`.
- Fuzzing (`func FuzzXxx(f *testing.F)`) *(read-only)*.

### Project layout, tooling, and deployment *(read-only)*

- Typical layout: `cmd/<app>/main.go` for entry points, `internal/` for private packages (the compiler forbids outside imports), small focused packages, no `utils`.
- Dependencies: `go get pkg@version`, `go mod tidy`, `go mod vendor` (optional), `go.sum` for integrity.
- Linting: `go vet`, `staticcheck`, `golangci-lint`.
- Profiling: `net/http/pprof` and `go tool pprof`; `go tool trace` for scheduler traces.
- Build: `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/app` gives one static binary,
  so Docker images can be a tiny multi-stage build (`golang` builder, `scratch` or `distroless` runtime).
- Databases: `database/sql` + a driver (`pgx`, ...) or `sqlc`; always pass `ctx`, always `rows.Close()`.
- Configuration: environment variables (12-factor) via `os.Getenv`; flags via the `flag` package.
- Common frameworks you will meet: `chi`, `gin`, `echo` (routers), `cobra` (CLIs), `grpc-go`, `zap`/`zerolog` (logging). The standard library covers more than you expect.

**Exercises:** `exercises/stage7` — `go test -race ./exercises/stage7`

---

## Further Reading

- [A Tour of Go](https://go.dev/tour/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)
- [Go by Example](https://gobyexample.com/)
- [The Go Blog: Slices intro](https://go.dev/blog/slices-intro), [Errors are values](https://go.dev/blog/errors-are-values), [Pipelines](https://go.dev/blog/pipelines)
- Standard library docs: <https://pkg.go.dev/std>
