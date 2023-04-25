package parser

import (
	"strings"
)

type Test []*Code

func (_this *Test) Value() []*Code {
	return *_this
}

func (_this *Test) ToString() string {
	res := make([]byte, 0)
	for _, code := range *_this {
		res = append(res, code.Value()...)
	}
	return string(res)
}

func (_this *Test) ToCodeHtml() string {
	b := strings.Builder{}
	for _, code := range *_this {
		b.WriteString(code.ToCodeHtml())
	}
	return b.String()
}
