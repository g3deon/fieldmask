# Go FieldMask

[![Actions Status](https://github.com/g3deon/fieldmask/actions/workflows/go.yml/badge.svg)](https://github.com/g3deon/fieldmask/actions)
[![Go Reference](https://pkg.go.dev/badge/go.g3deon.com/fieldmask.svg)](https://pkg.go.dev/go.g3deon.com/fieldmask)
[![Go Version](https://img.shields.io/badge/go-1.22+-blue.svg)](https://golang.org/dl/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A lightweight Go library for safely updating struct fields using dot-notation paths. It is fast, dependency-free, 
supports nested structs and JSON tags, and is optimized for partial updates in APIs.

## Installation

```sh
  go get go.g3deon.com/fieldmask
```

## Usage

### Basic Example

```go
package main

import (
  "fmt"
  "go.g3deon.com/fieldmask"
)

type Profile struct {
  Age int `json:"age"`
}

type User struct {
  Name    string  `json:"name"`
  Email   string  `json:"email"`
  Profile Profile `json:"profile"`
}

func main() {
  user := &User{
    Name:  "John",
    Email: "john@example.com",
    Profile: Profile{
      Age: 30,
    },
  }

  mask := fieldmask.New("name", "profile.age")

  if err := mask.Apply(user); err != nil {
    panic(err)
  }

  fmt.Printf("%+v\n", user)
  // Output: {Name:John Email: Profile:{Age:30}}
}
```

### Getting Paths

Use `GetPaths()` instead of accessing the `Paths` field directly.

```go
mask := fieldmask.New("name", "profile.age")

paths := mask.GetPaths()
fmt.Println(paths)
// Output: FieldMask{Paths: name, profile.age}
```

### Removing Paths

Remove one or multiple paths using `RemovePaths()`.

```go
mask := fieldmask.New("name", "profile.age")

mask.RemovePaths("name", "unknown")

fmt.Println(mask)
// Output: FieldMask{Paths: profile.age}
```

### Normalizing Paths

The `New` function automatically normalizes paths by removing duplicates and empty values. To manually normalize, 
use `Normalize()`.

```go
mask := fieldmask.FieldMask{
  Paths: []string{"name", "name", "", "profile.age"},
}

mask.Normalize()

fmt.Println(mask)
// Output: FieldMask{Paths: name, profile.age}
```

---

## License

MIT © 2025 G3deon, Inc.
