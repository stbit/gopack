package utils

func Sum(a int, b int) (r int, err error) {
	err = nil
	return a + b, nil
}

type RTest struct{}

type SFunc = func(int) int
