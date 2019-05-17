Go Application [![Last release](https://img.shields.io/github/release/euskadi31/go-application.svg)](https://github.com/euskadi31/go-application/releases/latest) [![Documentation](https://godoc.org/github.com/euskadi31/go-application?status.svg)](https://godoc.org/github.com/euskadi31/go-application)
================

[![Go Report Card](https://goreportcard.com/badge/github.com/euskadi31/go-application)](https://goreportcard.com/report/github.com/euskadi31/go-application)

| Branch  | Status | Coverage |
|---------|--------|----------|
| master  | [![Build Status](https://img.shields.io/travis/euskadi31/go-application/master.svg)](https://travis-ci.org/euskadi31/go-application) | [![Coveralls](https://img.shields.io/coveralls/euskadi31/go-application/master.svg)](https://coveralls.io/github/euskadi31/go-application?branch=master) |

go-application is a HTTP micro-framework library for Go.

Example
-------

```go
package main

import (
    "github.com/euskadi31/go-application"
    "github.com/euskadi31/go-application/provider"
)

func main() {
    app := application.New()

    app.Register(provider.NewEventDispatcherServiceProvider())
    app.Register(provider.NewHTTPServiceProvider())

    defer func() {
		if err := app.Close(); err != nil {
			panic(err)
		}
	}()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
```

License
-------

go-application is licensed under [the MIT license](LICENSE.md).
