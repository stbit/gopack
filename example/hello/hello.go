package hello

import (
	"fmt"
	"io"
	"os"

	"github.com/stbit/gopack/example/utils"
)

type Profile struct{}

func withError(s string) (string, error) {
	return s + "1", nil
}

func withOneError(s string) error {
	d, _ := os.OpenFile("", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755)
	defer d.Close()
	fmt.Println(s)
	return nil
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, _ := os.Stat(src)
	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, _ := os.Open(src)
	defer source.Close()

	destination, _ := os.Create(dst)
	defer destination.Close()

	nBytes, _ := io.Copy(destination, source)
	return nBytes, nil
}

func showHello() (string, func(int) int, Profile, *Profile, int, error) {
	f := func() func() (int, error) {
		return func() (int, error) {
			s, _ := withError("hello11")

			fmt.Print("suka", s)

			return 1, nil
		}
	}

	s, _ := withError("hello")
	q, _ := withError("hello test")
	err := withOneError("sdfs")
	withOneError("simple with error")
	utils.Sum(1, 2)

	_ = err
	k, _ := f()()
	r, _ := utils.Sum(1, 2)

	fmt.Println(s, q, k, r)

	return "", nil, Profile{}, nil, 0, nil
}
