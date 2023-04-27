package license

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var storeProd string
var storeClient string
var storeLC *LicenseCodeContent
var year string
var c chan bool

func init() {
	year = strconv.FormatInt(int64(time.Now().Year()), 10)
	c = make(chan bool, 1)
}

//PrintWelcome print the welcome information
func PrintWelcome() {
	<-c
	s := ""
	s += fmt.Sprintf("license to\n")
	s += fmt.Sprintf("%v\n", wrap(storeClient))
	s += fmt.Sprintf("exp date %v\n", storeLC.EndTime.Format("2006-01-02"))
	fmt.Println("12")
	colorPrint(box(s), wrap(storeClient))
}

func wrap(s string) string {
	defer func() {
		recover()
	}()
	l := strings.Split(s, " ")
	if len(l) == 1 {
		return s
	}
	st := make([]string, 0, len(l))
	end, n := 0, 0
	for i := range l {
		n += length(l[i])
		if n > 25 {
			n = 0
			end = i
			st = append(st, strings.Join(l[:i], " "))
		}
	}
	st = append(st, strings.Join(l[end:], " "))
	return strings.Trim(strings.Join(st, "\n"), "\n ")
}

func colorPrint(box, client string) {
	defer func() {
		recover()
	}()
	box = strings.Replace(box, "=", "\u001B[94m=\u001B[39m", -1)
	box = strings.Replace(box, "┃", "\u001B[94m┃\u001B[39m", -1)
	box = strings.Replace(box, "╋", "\u001B[94m╋\u001B[39m", -1)
	box = strings.Replace(box, client, "\u001B[31m"+client+"\u001B[39m", -1)
	box = strings.Replace(box,
		fmt.Sprintf("2016-%v (c) Hangzhou Qulian Technology Co.,Ltd.", year),
		fmt.Sprintf("\u001B[90m2016-%v (c) Hangzhou Qulian Technology Co.,Ltd.\u001B[39m", year),
		-1)
	reg, _ := regexp.Compile("[0-9]{4}-[0-9]{2}-[0-9]{2}")
	t := reg.FindString(box)
	if t != "" {
		box = strings.Replace(box, t, "\u001B[92m"+t+"\u001B[39m", -1)
	}
	fmt.Print(box)
}

func box(in string) string {
	defer func() {
		recover()
	}()
	if length(storeProd) < 11 && alpha(storeProd) {
		in = figure.NewFigure(storeProd, "drpepper", false).
			String() + in
	} else {
		in = figure.NewFigure("License", "drpepper", false).
			String() + in
	}

	in += fmt.Sprintf(" \n2016-%v (c) Hangzhou Qulian Technology Co.,Ltd.", year)
	in = strings.Trim(in, "\n")
	l := strings.Split(in, "\n")
	l = append([]string{" ", "WELCOME TO", " "}, l...)

	//get the Row height and max line width of WordArt
	var hb, he, lb, le, max1, max2, max int
	for i := range l {
		if strings.Contains(l[i], "WELCOME TO") {
			hb = i + 1
		}
		if strings.Contains(l[i], "license to") {
			he, lb = i, i
		}
		if strings.Contains(l[i], "exp date") {
			le = i
		}
	}
	for i := range l {
		if i > hb && i < he && length(l[i]) > max1 {
			max1 = length(l[i])
		}
		if i > lb && i < le && length(l[i]) > max2 {
			max2 = length(l[i])
		}
		if length(l[i]) > max {
			max = length(l[i])
		}
	}
	max += 4
	for i := range l {
		//left
		if strings.Contains(l[i], "license to") {
			l[i] += strings.Repeat(" ", max-length(l[i]))
			continue
		}
		//right
		if strings.Contains(l[i], "exp date") {
			l[i] = strings.Repeat(" ", max-length(l[i])) + l[i]
			continue
		}
		//middle
		var r int
		if i > hb && i < he {
			r = max - max1
			l[i] += strings.Repeat(" ", max1-length(l[i]))
		} else if i > lb && i < le {
			r = max - max2
			l[i] += strings.Repeat(" ", max2-length(l[i]))
		} else {
			r = max - length(l[i])
		}
		l[i] += strings.Repeat(" ", r/2)
		l[i] = strings.Repeat(" ", r-r/2) + l[i]
	}
	in = "    ┃    " + strings.Join(l, "    ┃\n    ┃    ")
	end := "    ╋" + strings.Repeat("=", max+8) + "╋\n"
	in = end + in + "    ┃\n" + end
	return in
}
func alpha(s string) bool {
	fmt.Sprintf(s)
	for i := range s {
		if !unicode.IsLetter(rune(s[i])) {
			return false
		}
	}
	return true
}

func length(s string) int {
	n := 0
	r := []rune(s)
	for i := range r {
		if r[i] > 0xff {
			n += 2
			continue
		}
		n++
	}
	return n
}
