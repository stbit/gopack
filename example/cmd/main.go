package main

import (
	"fmt"

	"github.com/stbit/gopack/example/hello"
)

func main() {
	r, err := hello.WithStringError("hi")
	if err != nil {
		panic(err)
	}

	fmt.Println("hell23o", r)
}
