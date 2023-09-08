package funcs

import "fmt"

//* why no strconv allowed??
func Atoi(s string) (int, error) {
	st := []rune(s)
	n := 0
	sign := 1
	if len(st) > 1 {
		if st[0] == '+' {
			st = st[1:]
		} else if st[0] == '-' {
			sign = -1
			st = st[1:]
		}
	}
	for i := 0; i < len(st); i++ {
		if st[i] < '0' || st[i] > '9' || st[i] == ' ' {
			return 0, fmt.Errorf("invalid argument for Atoi")
		}
		n = (n*10 + int(st[i]) - '0')
	}
	return sign * n, nil
}

