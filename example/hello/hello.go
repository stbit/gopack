package hello

import "fmt"

type Profile struct{}

func withError(s string) (string, error) {
	return s + "1", nil
}

func withOneError(s string) error {
	fmt.Println(s)
	return nil
}

func showHello() (string, func(int) int, Profile, *Profile, int, error) {
	f := func() func() error {
		return func() error {
			s, _ := withError("hello11")

			fmt.Print("suka", s)

			return nil
		}
	}

	s, _ := withError("hello")
	q, _ := withError("hello test")
	err := withOneError("sdfs")

	_ = err

	fmt.Println(s, q, f()())

	return "", nil, Profile{}, nil, 0, nil
}
