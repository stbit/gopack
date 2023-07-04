package main

import (
	"encoding/json"
	"fmt"

	"github.com/rotisserie/eris"
	"github.com/stbit/gopack/example/hello"
)

func main() {
	_, err := hello.WithStringError("hi")
	if err != nil {
		panic(err)
	}
	err = hello.GetError()
	js, err := json.Marshal(eris.ToJSON(err, true))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(js))
}
