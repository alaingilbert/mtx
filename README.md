**From Error-Prone to Safe: How `mtx` Transforms Unsafe Patterns**

#### **The Problem: Manual Mutex Management**
In Go, managing shared resources with mutexes is error-prone. Developers must remember to:
1. **Lock** the mutex before accessing shared data.
2. **Unlock** it afterward (often using `defer`).
3. **Document** which fields are protected by the mutex.

Here’s a typical example:

```go
type SomeStruct struct {
    // WARNING: mx protects the shared slice and the timestamp
    // You shall never use the sharedSlice nor the timestamp without holding mx
    mx          sync.Mutex
    sharedSlice []int
    timestamp   time.Time
}

func (s *SomeStruct) DoSomething(el int) {
    s.mx.Lock()
    defer s.mx.Unlock()
    s.doSomethingInternally(el)
}

func (s *SomeStruct) doSomethingInternally(el int) {
    // WARNING: Caller must not forget to lock the shared resources before coming here
    s.sharedSlice = append(s.sharedSlice, el)
    s.timestamp = time.Now()
}
```

**Issues:**
- Forgetting to lock/unlock can lead to race conditions.
- Documentation (`WARNING` comments) is easy to miss or ignore.
- The mutex and its protected fields are loosely coupled, making it hard to enforce safety.

---

#### **The Solution: `mtx` for Enforced Safety**
The `mtx` library eliminates these pitfalls by:
1. **Encapsulating the mutex and its protected data** in a single type (`mtx.Mutex[container]`).
2. **Enforcing lock acquisition** before accessing the data via `With(func(c *container) { ... })`.

Here’s the improved version:

```go
type SomeStruct struct {
    // No longer need to explain what is being protected, this is self-explanatory
    inner mtx.Mutex[container]
}

type container struct {
    sharedSlice []int
    timestamp   time.Time
}

func (s *SomeStruct) DoSomething(el int) {
    // This construct makes it much harder to accidentally forget to release the lock
    s.inner.With(func(c *container) {
        doSomethingInternally(c, el)
    })
}

// It is not possible to come here without having a pointer to container,
// which you can only get by holding the lock.
// So the caller cannot forget to lock the resources before calling this function. 
func doSomethingInternally(c *container, el int) {
    (*c).sharedSlice = append((*c).sharedSlice, el)
    (*c).timestamp = time.Now()
}
```

**Benefits:**
- **No More Forgotten Locks**: The `With` method ensures the lock is held for the duration of the callback.
- **Self-Documenting**: The `container` type clearly groups protected fields, eliminating the need for `WARNING` comments.
- **Compiler-Enforced Safety**: The `doSomethingInternally` function can’t be called without holding the lock, as `c *container` is only accessible with the mutex being acquired.

---

#### **Key Takeaways**
- **Eliminate Boilerplate**: No more manual `Lock()`/`Unlock()` calls.
- **Reduce Human Error**: The compiler enforces correct usage.
- **Cleaner Code**: Protected data is explicitly grouped, making the design intent clear.

By adopting `mtx`, you trade manual mutex management for a safer, more maintainable approach. Less room for mistakes, more time for solving real problems!