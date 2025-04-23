[![Go Reference](https://pkg.go.dev/badge/github.com/alaingilbert/mtx.svg)](https://pkg.go.dev/github.com/alaingilbert/mtx)

### Golang generic mutex helpers

```go
package main

import (
    "fmt"
    "github.com/alaingilbert/mtx"
)

func main() {
    // go run -race main.go
    type Something struct {
        Field1    string
		SharedMap mtx.MapMutex[string, int]
    }

    something := Something{
        Field1:    "this memory is not being shared, no mutex needed on Field1",
		SharedMap: mtx.NewMapMutex(make(map[string]int)),
    }

    for i := 0; i < 100; i++ {
        go something.SharedMap.Insert("a", i)
    }

    fmt.Println(something.SharedMap.Get("a"))
}
```

```go
package main

import (
    "fmt"
    "github.com/alaingilbert/mtx"
)

func main() {
    type Something struct {
        Field1           string
        SharedInt        mtx.Mtx[int]
        SharedFloat64    mtx.Mtx[float64]
        SharedMap1       mtx.Map[string, int]
        SharedMap2       mtx.Map[string, int]
        SharedSlice1     mtx.Slice[int]
        SharedSlice2     mtx.Slice[int]
        SharedSlicePtr1  *mtx.Slice[int]
        SharedSlicePtr2  *mtx.Slice[int]
    }
    something := Something{
        Field1:           "",
        SharedInt:        mtx.NewMtx(0),                        // uses sync.Mutex
        SharedFloat64:    mtx.NewRWMtx(0.0),                    // uses sync.RWMutex
        SharedMap1:       mtx.NewMap(map[string]int{"a": 1}),   // uses sync.Mutex
        SharedMap2:       mtx.NewRWMap(map[string]int{"b": 2}), // uses sync.RWMutex
        SharedSlice1:     mtx.NewSlice([]int{1, 2, 3}),         // uses sync.Mutex
        SharedSlice2:     mtx.NewRWSlice([]int{4, 5, 6}),       // uses sync.RWMutex
        SharedSlicePtr1:  mtx.NewSlicePtr([]int{7, 8, 9}),      // uses sync.Mutex
        SharedSlicePtr2:  mtx.NewRWSlicePtr([]int{10, 11, 12}), // uses sync.RWMutex
    }
    fmt.Println(something)
}
```

## Goal

It is not unusual in Go to see code like this,  
where the user has to not forget to use the mutex, and has to not make a mistake with the unlocking mechanism.
```go
// NOTE: This code block is NOT an example on how to use the library!

type Something struct {
    SharedMapMtx sync.RWMutex
    SharedMap    map[string]int
}

func someFn(s Something) {
    s.SharedMapMtx.Lock()
    defer s.SharedMapMtx.Unlock()
    s.SharedMap["foo"] = 1
}
```

This library ensure that a field which is protected by a mutex will be used properly without compromising on flexibility.
```go
type Something struct {
    SharedMap    mtx.Map[string, int]
}

// This is the recommended way of setting a key
func someFn(s *Something) {
    s.SharedMap.Insert("foo", 1)
}

// This is also good
func someOtherFn(s *Something) {
    s.SharedMap.With(func(sharedMapPtr *map[string]int) {
        (*sharedMapPtr)["foo"] = 1
    })
}

// But it can be as flexible as needed
func anotherOneFn(s *Something) {
    s.SharedMap.Lock()
    defer s.SharedMap.Unlock()
    sharedMapPtr := s.SharedMap.GetPointer()
    (*sharedMapPtr)["foo"] = 1
}
```