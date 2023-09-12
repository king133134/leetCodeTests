package leetcodeTests

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/king133134/leetCodeTests/parser"
	"log"
	"strings"
)

func InputRun() {
	var typeCode string
	fmt.Print("请输入您要生成测试代码的变量(eg: int:3, *TreeNode:[3,0,0]):")
	_, err := fmt.Scan(&typeCode)
	if err != nil {
		log.Fatal(err)
	}
	if typeCode == "" {
		panic("URL is empty!")
	}
	list := strings.Split(typeCode, ":")
	code := parser.InputToCode(list[0], list[1])
	err = clipboard.WriteAll(code)
	if err != nil {
		panic(err)
	}
	fmt.Println("代码:", code)
	fmt.Println(color2output("代码已复制到剪切板.", Green))
}

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
)

func color2output(output string, color string) string {
	// 输出带有颜色的文本
	return color + output + Reset
}
