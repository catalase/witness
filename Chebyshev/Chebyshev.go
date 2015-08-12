package main

import (
	"fmt"
	"unicode"
	"github.com/atotto/clipboard"
)

// 정수 계수 다항식
// 첫번째 원소는 상수항을 나타낸다. 차수는 순차적으로 증가한다.
type Po []int

func maxint(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// x + y
func Add(x, y Po) Po {
	z := make(Po, maxint(len(x), len(y)))
	copy(z, x)
	
	for i, r := range y {
		z[i] += r
	}

	return z
}

// x - y
func Sub(x, y Po) Po {
	z := make(Po, maxint(len(x), len(y)))
	copy(z, x)
	
	for i, r := range y {
		z[i] -= r
	}

	return z	
}

// Shift 는 모든 차수를 증가시킨다. 
func Shift(x Po, n int) Po {
	z := make(Po, len(x) + n)
	copy(z[n:], x)

	return z
}

type byteWriter []byte

func (bytearray *byteWriter) Write(data []byte) (int, error) {
	*bytearray = append(*bytearray, data...)
	return len(data), nil
}

func Latex(x Po) string {
	var w byteWriter

	for i := len(x) - 1; i >= 0; i-- {
		r := x[i]
		if r == 0 {
			continue
		}

		if i > 0 {
			switch r {
			case 1:
				w = append(w, '+')
			case -1:
				w = append(w, '-')
			default:
				fmt.Fprintf(&w, "%+d", r)
			}

			if i >= 10 {
				fmt.Fprintf(&w, "x^{%d}", i)
			} else {
				fmt.Fprintf(&w, "x^%d", i)
			}
		} else {
			fmt.Fprintf(&w, "%+d", r)
		}
	}

	if len(w) > 0 {
		if w[0] == '+' {
			w = w[1:]
		}
	}

	return string(w)
}

var cacheCosN = make(map[int]Po)

func CosN(n int) Po {
	if n < 0 {
		panic("given negative value")
	}

	if n == 1 {
		return Po{0, 1}
	}

	if n == 2 {
		return Po{-1, 0, 2}
	}

	if z, ok := cacheCosN[n]; ok {
		return z
	}

	z := Shift(CosN(n - 1), 1)
	for i := range z {
		z[i] *= 2
	}

	z = Sub(z, CosN(n - 2))
	cacheCosN[n] = z

	return z
}

func want() bool {
	var cp rune
	fmt.Scanf("%c\n", &cp)
	return unicode.ToLower(cp) == 'y'
}

func main() {
	var n int

	fmt.Print("cos(nx) where n = ")
	fmt.Scanf("%d\n", &n)

	z := CosN(n)

	fmt.Println("cos(nx) is:")
	fmt.Println(z)

	fmt.Print("do you want to copy latex of cos(nx) into clipboard? (y or n): ")
	if want() {
		clipboard.WriteAll(Latex(z))
	}
}