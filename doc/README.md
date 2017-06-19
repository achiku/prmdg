## main resource

Simply put, main resource is `#/definitions/{main_resource}`. It is defined right under the first `definitions`, and will be treated differently from other definitions. Other definitions might be scalar, object, or array, but main resources are transformed into `struct` in go semantics. Furthermore, it can be reused in other main resources. For example, `type Todo struct` can contain `type User struct` although those two are separately defined as `struct`.


```golang
// User user resource
type User struct {
	LastName  string `json:"lastName"`
	FirstName string `json:"firstName"`
	Age       int    `json:"age"`
}

// Todo todo resource
type Todo struct {
	User  *User    `json:"user"`
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
}
```

## rel tag
