package main

import (
	"fmt"
	"time"

	"github.com/stbit/gopack/example/hello"
)

func main() {
	r, err := hello.WithStringError("hi")
	if err != nil {
		panic(err)
	}
	fmt.Println("sdfsdf244")

	time.Sleep(5 * time.Second)

	fmt.Println("hell23o234", r)
}
