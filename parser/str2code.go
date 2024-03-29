package parser

import (
	"bytes"
	"golang.org/x/net/html"
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
	case p.val == "[]byte" || p.val == "[][]byte":
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
	case p.val == "[][]int" || p.val == "[]int" || p.val == "[]int64" || p.val == "[]string" || p.val == "[][]string" || p.val == "[]bool" || p.val == "[][]bool" || p.val == "[]float64" || p.val == "[][]float64" || p.val == "[]float32" || p.val == "[][]float32":
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
	case p.val == "int" || p.val == "int64" || p.val == "int32" || p.val == "float32" || p.val == "float64" || p.val == "bool" || p.val == "string":
		res.AppendBytes([]byte(val[i:]))
	case p.val == "*TreeNode":
		res.AppendBytes(toTreeNode(val[i:]))
	case p.val == "*ListNode":
		res.AppendBytes(toListNode(val[i:]))
	case p.val == "[]*ListNode":
		res.AppendBytes(toListNodeSlice(val[i:]))
	case p.val == "[]*TreeNode":
		res.AppendBytes(toTreeNodeSlice(val[i:]))
	default:
		panic("str2code func type:" + p.val + " is not match!")
	}
	return res
}

// InputToCode eg: input = [10,0,0],paramType = *TreeNode, return = &TreeNode{10, &TreeNode{0,nil,nil}, &TreeNode{0,nil,nil}}
func InputToCode(paramType, input string) string {
	res := str2code(&param{val: paramType}, input)
	return res.String()
}

// eg:[[1,4,5],[1,3,4],[2,6]]
func toListNodeSlice(str string) []byte {
	res := []byte("[]*ListNode{")
	i, j, n := 0, 0, len(str)
	for str[i] != '[' {
		i++
	}
	i++
	for i < n {
		for i < n && str[i] != '[' {
			i++
		}
		j = i
		for j < n && str[j] != ']' {
			j++
		}
		if j-i < 2 {
			res = append(res, ',')
			break
		}
		res = append(res, toListNode(str[i:j+1])...)
		res = append(res, ',')
		i = j + 1
	}
	res[len(res)-1] = '}'
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

	if pre < i {
		items = append(items, str[pre:i])
	}
	n := len(items)
	if n == 0 {
		return []byte("nil")
	}

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

func toTreeNodeSlice(str string) []byte {
	res := []byte("[]*TreeNode{")
	i, j, n := 0, 0, len(str)
	for str[i] != '[' {
		i++
	}
	i++
	for i < n {
		for i < n && str[i] != '[' {
			i++
		}
		j = i
		for j < n && str[j] != ']' {
			j++
		}
		if j-i < 2 {
			res = append(res, ',')
			break
		}
		res = append(res, toTreeNode(str[i:j+1])...)
		res = append(res, ',')
		i = j + 1
	}
	res[len(res)-1] = '}'
	return res
}

// 把测试用例的数据转成TreeNode的代码(eg:[1,2,3])
func toTreeNode(str string) []byte {
	items := strings.Split(strings.Trim(str, "[] "), ",")
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
	pn := 0
	for _, v := range params {
		if v.key != "output" {
			pn++
		}
	}
	for i, io := range tests {
		t := &Test{}
		j, item := 0, ""
		start := 0
		for k := 0; k < pn; k++ {
			item, j = indexParam(io.input[start:], params[k+1].key)
			start += j
			idx := strings.IndexByte(item, '=')
			*t = append(*t, str2code(params[k], item[idx+1:]))
		}
		*t = append(*t, str2code(params[pn], io.output))
		res[i] = t
	}
	return newTests(res, fName, params)
}

func indexParam(input, key string) (string, int) {
	if key == "output" {
		return input, len(input)
	}
	idx := strings.Index(input, key+" = ")
	for idx > 0 && input[idx] != ',' {
		idx--
	}
	return input[:idx], idx
}

func params(code string) (res []*param, fName string) {
	n := len(code)
	var f []byte
	j := strings.LastIndex(code, "func ")
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

	end := strings.IndexByte(code[j:], ')')
	items := strings.Split(code[j+1:j+end], ",")
	// root, p, q *TreeNode
	valStr := ""
	res = make([]*param, len(items))
	for i := len(items) - 1; i >= 0; i-- {
		kv := strings.Split(strings.Trim(items[i], " "), " ")
		if len(kv) != 1 {
			valStr = kv[1]
		}
		res[i] = &param{kv[0], valStr}
	}

	var key, val []byte
	j += end
	key = append(key, "output"...)
	for code[j] == ' ' || code[j] == ')' {
		j++
	}

	for j < n && code[j] != ' ' && code[j] != '{' {
		val = append(val, code[j])
		j++
	}
	if len(val) == 0 {
		panic("output is empty.")
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
	for ; (*con)[i] != '<' || ((*con)[i:i+8] != "<strong>" && (*con)[i:i+8] != "<strong " && (*con)[i:i+6] != "</pre>" && (*con)[i:i+3] != "<p>"); i++ {
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
			start, text := buildIO(con, i)
			i = start
			cur.input = removeHTMLTags(text)
		} else if (*con)[i-6:i] == "Output" {
			start, text := buildIO(con, i)
			i = start
			cur.output = removeHTMLTags(text)
			res = append(res, cur)
			cur = testIO{}
		}
	}
	return res
}

func removeHTMLTags(text []byte) string {
	doc, _ := html.Parse(bytes.NewReader(text))
	var result strings.Builder

	var getText func(*html.Node)
	getText = func(n *html.Node) {
		if n.Type == html.TextNode {
			result.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			getText(c)
		}
	}

	getText(doc)
	return strings.Trim(result.String(), " ")
}
