Golang generic mutex helpers

```go
shared := mtx.NewMtx("some value")
fmt.Println(shared.Get())
shared.Set("new value")
```

```go
package main

import (
	"fmt"
	"github.com/alaingilbert/mtx"
)

func main() {
	type Something struct {
		Field1    string
		SharedMap mtx.Map[string, int]
	}

	something := Something{
		Field1:    "",
		SharedMap: mtx.NewRWMap[string, int](nil),
	}

	for i := 0; i < 100; i++ {
		go func(i int) {
			something.SharedMap.SetKey("a", i)
		}(i)
	}

	fmt.Println(something.SharedMap.GetKey("a"))
}
```