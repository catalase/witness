package main

import (
	"fmt"
	"github.com/atotto/clipboard"
	"io"
	"os"
	"unicode"
	"log"
)

import (
	"image"
	_ "image/gif"
	"image/png"
	"net/http"
	"net/url"
	"io/ioutil"
)

var _ = png.Encode

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
	z := make(Po, len(x)+n)
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

	z := Shift(CosN(n-1), 1)
	for i := range z {
		z[i] *= 2
	}

	z = Sub(z, CosN(n-2))
	cacheCosN[n] = z

	return z
}

func RenderStringAsURL(str string) (string, error) {
	v := make(url.Values)
	v.Set("val", str)

	u, _ := url.Parse("http://www.numberempire.com/texequationeditor/get_tex_image_url.php")
	u.RawQuery = v.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return "", err
	}

	url, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}

	return string(url), nil
}

func RenderString(str string) (image.Image, error) {
	url, err := RenderStringAsURL(str)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	m, _, err := image.Decode(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return m, nil
}

func minint(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func LatexPo(z Po, n, sep int) string {
	var w byteWriter

	io.WriteString(&w, `\begin{align}`)
	fmt.Fprintf(&w, `\cos(%dx) = `, n)

	var arrange [][2]int

	for r := len(z) - 1; r >= 0; r-- {
		x := z[r]
		if x != 0 {
			arrange = append(
				arrange,
				[2]int{r, x},
			)
		}
	}

	for i := 0; i < len(arrange); i += sep {
		io.WriteString(&w, "&")
		for _, term := range arrange[i:minint(i+sep, len(arrange))] {
			r, x := term[0], term[1]
			if r > 0 {
				switch r {
				case 1:
					fmt.Fprintf(&w, "+cos^{%d}(x)", r)
				case -1:
					fmt.Fprintf(&w, "-cos^{%d}(x)", r)
				default:
					fmt.Fprintf(&w, "%+dcos^{%d}(x)", x, r)
				}
			} else {
				fmt.Fprintf(&w, "%+d", x)
			}
		}
		io.WriteString(&w, "\\\\\n")
	}

	io.WriteString(&w, `\end{align}`)

	return string(w)
}

func want() bool {
	var cp rune
	fmt.Scanf("%c\n", &cp)
	return unicode.ToLower(cp) == 'y'
}

func test2() {
	var n int

	fmt.Print("cos(nx) where n = ")
	fmt.Scanf("%d\n", &n)

	z := CosN(n)

	fmt.Println("cos(nx) is:")
	fmt.Println(z)

	url, err := RenderStringAsURL(LatexPo(z, n, 4))
	if err == nil {
		fmt.Print("do you want to copy latex address into clipboard? (y or n): ")
		if want() {
			clipboard.WriteAll(url)
		}
	} else {
		fmt.Println(err)
	}

}

func main() {
	for i := 1; i <= 100; i++ {
		z := CosN(i)
		im, err := RenderString(LatexPo(z, i, 4))
		if err != nil {
			log.Fatal(err)
		}

		w, err := os.Create(fmt.Sprintf("cos(%dx).png", i))
		if err != nil {
			log.Fatal(err)
		}

		if err := png.Encode(w, im); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Done cos(%dx)\n", i)
	}
}