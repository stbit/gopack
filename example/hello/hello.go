package hello

import "fmt"

func withError(s string) (string, error) {
	return s + "1", nil
}

func showHello() {
	s, _ := withError("hello")

	fmt.Println(s)
}
