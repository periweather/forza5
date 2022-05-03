package main

import (
	"context"
	"github.com/periweather/forza5"
	"strings"
)

func main() {
	ctx := context.Background()
	err := forza5.Client(ctx, "0.0.0.0:3000", strings.NewReader("a string reader"))
	if err != nil {
		panic(err)
	}
}
