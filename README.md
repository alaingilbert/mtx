Golang generic mutex helpers

```go
shared := mtx.NewMtx("some value")
fmt.Println(shared.Get())
shared.Set("new value")
```

```go
type Something struct {
	Field1    string
	SharedMap mtx.Map[string, int]
}
something := Something{
	Field1:    "",
	SharedMap: mtx.NewMap[string, int](),
}
something.SharedMap.SetKey("a", 1)
fmt.Println(something.SharedMap.GetKey("a"))
```