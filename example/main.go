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

type ListNode struct {
	Val  int
	Next *ListNode
}

func main() {
	var mod = flag.String("mod", "http", "run mod")
	var port = flag.Int("port", 8080, "http port")
	var dir = flag.String("dir", ".", "file dir")

	flag.Parse()
	leetcodeTests.Run(leetcodeTests.Mod(*mod), leetcodeTests.Dir(*dir), leetcodeTests.Port(*port))
}
