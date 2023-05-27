package hello

import (
	"fmt"
	"io"
	"os"

	"github.com/stbit/gopack/example/utils"
)

type Profile struct {
	Name       string `validate:"required"`
	AppId      string `json:",omitempty,string"`
	CountUsers int
	NestedObj  struct {
		Name         string
		PeoplesStudy int
	}
}

type IBase interface{}

func (p *Profile) save() error {
	s, _ := withError("dsf")
	fmt.Println(s)

	return nil
}

func withError(s string) (string, error) {
	return s + "1", nil
}

func WithStringError(s string) (string, error) {
	return s + "1", nil
}

func withInterface() error {
	withError("sdfds")

	return nil
}

func withOneError(s string) error {
	d, _ := os.OpenFile("", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755)
	// SetMaxIdleConns устанавливает максимальное количество соединений в пуле бездействия.
	// sqlDB.SetMaxIdleConns(50)

	// // SetMaxOpenConns устанавливает максимальное количество открытых соединений с БД.
	// sqlDB.SetMaxOpenConns(100)
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

func showHello() (string, func(int) int, Profile, IBase, *Profile, int, utils.RTest, utils.SFunc, error) {
	f := func() func() (int, error) {
		return func() (int, error) {
			s, _ := withError("hello11")

			fmt.Print("suka", s)

			return 1, nil
		}
	}

	nt := struct {
		Name       string
		CountUsers int
	}{"hello", 4}

	s, _ := withError("hello")
	q, _ := withError("hello test")
	err := withOneError("sdfs")
	g1, g2 := err, 3
	withOneError("simple with error")
	utils.Sum(1, 2)

	// if err := withOneError("simple with error"); err != nil {
	// 	return "", nil, Profile{}, nil, 0, nil
	// }

	_ = err
	k, _ := f()()
	r, _ := utils.Sum(1, 2)

	fmt.Println(s, q, k, r, g1, g2, nt)

	return "", nil, Profile{}, nil, nil, 0, utils.RTest{}, nil, nil
}
