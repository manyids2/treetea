# tree

Given:

- List of items ( i.e. key, children )

Constructions:

```go
// Assumed item interface ( private to avoid name clashes with struct fields )
type Item interface {
 key() string
 children() []string
 desc(string) string // For printing
}

type Items []Item // Mapping to list indices
```

| Field    | Type              | Comments                                                            |
| -------- | ----------------- | ------------------------------------------------------------------- |
| Items    | Items             | Should be satisfied by Tasks, Projects, Contexts, Tags and Historys |
| Levels   | map[string]int    | Needed to track indent                                              |
| Parents  | map[string]string | Reverse tree                                                        |
| Children | []string          | Top level items                                                     |
| Order    | []string          | Current viewing order                                               |
