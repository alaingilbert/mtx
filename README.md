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