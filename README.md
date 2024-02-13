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
		SharedMap mtx.Map[string, int]
	}

	something := Something{
		Field1:    "",
		SharedMap: mtx.NewRWMap[string, int](nil),
	}

	for i := 0; i < 100; i++ {
		go something.SharedMap.SetKey("a", i)
	}

	fmt.Println(something.SharedMap.GetKey("a"))
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
		Field1        string
		SharedInt     mtx.Mtx[int]
		SharedFloat64 mtx.RWMtx[float64]
		SharedMap     mtx.Map[string, int]
		SharedRWMap   mtx.Map[string, int]
		SharedSlice   mtx.Slice[int]
		SharedRWSlice mtx.Slice[int]
	}
	something := Something{
		Field1:        "",
		SharedInt:     mtx.NewMtx(0),                        // uses sync.Mutex
		SharedFloat64: mtx.NewRWMtx(0.0),                    // uses sync.RWMutex
		SharedMap:     mtx.NewMap(map[string]int{"a": 1}),   // uses sync.Mutex
		SharedRWMap:   mtx.NewRWMap(map[string]int{"b": 2}), // uses sync.RWMutex
		SharedSlice:   mtx.NewSlice([]int{1, 2, 3}),         // uses sync.Mutex
		SharedRWSlice: mtx.NewRWSlice([]int{4, 5, 6}),       // uses sync.RWMutex
	}
	fmt.Println(something)
}
```

## Goal

It is not unusual in Go to see code like this,  
where the user has to not forget to use the mutex, and has to not make a mistake with the unlocking mechanism.
```go
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
    SharedMapMtx sync.RWMutex
    SharedMap    mtx.Map[string, int]
}

// This is the recommended way of setting a key
func someFn(s Something) {
    s.SharedMap.SetKey("foo", 1)
}

// But it can be as flexible as needed
func someOtherFn(s Something) {
    s.SharedMapMtx.Lock()
    defer s.SharedMapMtx.Unlock()
    sharedMapPtr := s.SharedMap.Val()
    (*sharedMapPtr)["foo"] = 1
}
```