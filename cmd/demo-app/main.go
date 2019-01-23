// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"github.com/euskadi31/go-application"
)

func main() {
	app := application.New()

	defer func() {
		if err := app.Close(); err != nil {
			panic(err)
		}
	}()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
