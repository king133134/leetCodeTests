package main

import (
	"flag"
	"github.com/king133134/leetCodeTests"
)

func main() {
	var mod = flag.String("mod", "http", "run mod")
	var port = flag.Int("port", 8080, "http port")
	var dir = flag.String("dir", ".", "file dir")

	flag.Parse()
	leetcodeTests.Run(leetcodeTests.Mod(*mod), leetcodeTests.Dir(*dir), leetcodeTests.Port(*port))
}
