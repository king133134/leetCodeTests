package parser

import (
	"html"
	"strings"
)

type param struct {
	key string
	val string
}

// 通过题目的测试用例给的参数转成代码
func str2code(p *param, val string) *Code {
	val = html.UnescapeString(val)
	i, n := 0, len(val)
	for val[i] == ' ' {
		i++
	}
	res := &Code{}
	switch {
	case p.val == "[]byte":
		res.Append([]byte(p.val)...)
		for i < n {
			b := val[i]
			if b == '[' {
				b = '{'
			}
			if b == ']' {
				b = '}'
			}
			if b == '"' {
				b = '\''
			}
			res.Append(b)
			i++
		}

	case p.val == "[][]int" || p.val == "[]int" || p.val == "[]string" || p.val == "[]bool":
		res.AppendBytes([]byte(p.val))
		for i < n {
			b := val[i]
			if b == '[' {
				b = '{'
			}
			if b == ']' {
				b = '}'
			}
			res.Append(b)
			i++
		}
	case p.val == "int" || p.val == "float32" || p.val == "float64" || p.val == "bool" || p.val == "string":
		res.AppendBytes([]byte(val[i:]))
	case p.val == "*TreeNode":
		res.AppendBytes(toTreeNode(val[i:]))
	case p.val == "*ListNode":
		res.AppendBytes(toListNode(val[i:]))
	default:
		panic("str2code func type:" + p.val + " is not match!")
	}
	return res
}

func toListNode(str string) []byte {
	i := 0
	for str[i] != '[' {
		i++
	}

	var items []string
	pre := i + 1
	for i++; str[i] != ']'; i++ {
		if str[i] == ',' {
			items = append(items, str[pre:i])
			pre = i + 1
		}
	}

	items = append(items, str[pre:i])
	n := len(items)

	null := []byte("nil")
	head := []byte("&ListNode{")
	tail := []byte("}")
	val := []byte("Val:")
	next := []byte("Next:")

	var dfs func(i int) []byte
	dfs = func(i int) []byte {
		if i == n {
			return null
		}
		res := make([]byte, 0)
		res = append(res, head...)
		res = append(res, append(val, items[i]...)...)
		res = append(res, ',')
		res = append(res, append(next, dfs(i+1)...)...)
		res = append(res, tail...)
		return res
	}
	return dfs(0)
}

// 把测试用例的数据转成TreeNode的代码
func toTreeNode(str string) []byte {
	i := 0
	for str[i] != '[' {
		i++
	}

	var items []string
	pre := i + 1
	for i++; str[i] != ']'; i++ {
		if str[i] == ',' {
			items = append(items, str[pre:i])
			pre = i + 1
		}
	}

	items = append(items, str[pre:i])
	null := []byte("nil")
	head := []byte("&TreeNode{")
	tail := []byte("}")
	val := []byte("Val:")
	left := []byte("Left:")
	right := []byte("Right:")
	n := len(items)
	offset := make([]int, n)
	preSum := 0
	for j, item := range items {
		if item == "null" {
			preSum += 2
		}
		offset[j] = preSum
	}
	var dfs func(i int) []byte
	dfs = func(i int) []byte {
		if items[i] == "null" {
			return null
		}
		res := make([]byte, 0)
		res = append(res, head...)
		res = append(res, append(val, items[i]...)...)
		res = append(res, ',')
		if r := i*2 + 2 - offset[i]; r < n {
			res = append(res, append(left, dfs(r-1)...)...)
			res = append(res, ',')
			res = append(res, append(right, dfs(r)...)...)
		} else if r-1 < n {
			res = append(res, append(left, dfs(r-1)...)...)
			res = append(res, ',')
			res = append(res, append(right, null...)...)
		} else {
			res = append(res, append(left, null...)...)
			res = append(res, ',')
			res = append(res, append(right, null...)...)
		}
		res = append(res, tail...)
		return res
	}

	return dfs(0)
}

// CreateTests 根据代码片段和题目内容的描述生成测试用例
func CreateTests(code string, con *string) *Tests {
	p, f := params(code)
	return createTests(p, f, getContentIO(con))
}

func createTests(params []*param, fName string, tests []testIO) *Tests {
	res := make([]*Test, len(tests))
	for i, io := range tests {
		t := &Test{}
		j, item := 0, ""
		for j, item = range strings.Split(io.input, ", ") {
			idx := strings.IndexByte(item, '=')
			*t = append(*t, str2code(params[j], item[idx+1:]))
		}
		*t = append(*t, str2code(params[j+1], io.output))
		res[i] = t
	}
	return newTests(res, fName, params)
}

func params(code string) (res []*param, fName string) {
	var f []byte
	j := 0
	for code[j] != '(' {
		b := code[j]
		if j > 3 && b == ' ' && code[j-4:j] == "func" {
			j++
			for c := code[j]; (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9'); {
				f = append(f, c)
				j++
				c = code[j]
			}
			j--
		}
		j++
	}

	var key, val []byte
	for j++; code[j] != ')'; {
		if code[j] == ' ' || code[j] == ',' {
			j++
			continue
		}
		for code[j] != ' ' {
			key = append(key, code[j])
			j++
		}
		for code[j] == ' ' {
			j++
		}
		for code[j] != ',' && code[j] != ')' {
			val = append(val, code[j])
			j++
		}
		res = append(res, &param{string(key), string(val)})
		key = key[:0]
		val = val[:0]
	}

	key = append(key, "output"...)
	for code[j] == ' ' || code[j] == ')' {
		j++
	}

	for code[j] != ' ' {
		val = append(val, code[j])
		j++
	}
	res = append(res, &param{string(key), string(val)})
	return res, string(f)
}

type testIO struct {
	input  string
	output string
}

func buildIO(con *string, i int) (start int, bytes []byte) {
	i += 10
	for (*con)[i] == ' ' || (*con)[i] == '\n' {
		i++
	}
	for ; (*con)[i] != '<'; i++ {
		b := (*con)[i]
		if b == '\n' || b == '\r' {
			continue
		}
		bytes = append(bytes, b)
	}
	start = i
	return
}

func getContentIO(con *string) []testIO {
	var res []testIO
	n := len(*con)
	var cur testIO

	for i := 0; i < n; i++ {
		b := (*con)[i]
		if b != ':' {

		} else if (*con)[i-5:i] == "Input" {
			start, bytes := buildIO(con, i)
			i = start
			cur.input = string(bytes)
		} else if (*con)[i-6:i] == "Output" {
			start, bytes := buildIO(con, i)
			i = start
			cur.output = string(bytes)
			res = append(res, cur)
			cur = testIO{}
		}
	}
	return res
}
