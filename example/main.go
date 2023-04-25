package main

import (
	"flag"
	"leetcodeTests"
)

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func main() {
	var mod = flag.String("mod", "http", "run mod")
	var port = flag.Int("port", 8080, "http port")
	var dir = flag.String("dir", ".", "file dir")

	flag.Parse()
	leetcodeTests.Run(leetcodeTests.Mod(*mod), leetcodeTests.Dir(*dir), leetcodeTests.Port(*port))
	// data := toTreeNode("[1,null,2,null,0,3,null,4,5]")
	// fmt.Println(string(data))
}

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
