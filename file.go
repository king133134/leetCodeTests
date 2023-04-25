package leetcodeTests

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func FileRun(dir string) {
	var url string
	fmt.Print("请输入您要生成测试用例的题目的URL：")
	_, err := fmt.Scan(&url)
	if err != nil {
		log.Fatal(err)
	}
	id := url2id(url)
	file(dir, id)
}

func file(dir, id string) {
	data := question(id)
	dir = strings.TrimRight(dir, "/ ") + "/problem_" + data.id
	_, err := os.Stat(dir)
	if !os.IsNotExist(err) {
		panic(fmt.Sprintf("dir %s is exists.", dir))
	}
	// mkdir 创建目录，mkdirAll 可创建多层级目录
	_ = os.MkdirAll(dir, 0755)
	fmt.Println("create dir:", dir)

	funcName := data.tests.FuncName()
	codeFile, testFile := fmt.Sprintf("%s/%s.go", dir, funcName), fmt.Sprintf("%s/%s_test.go", dir, funcName)

	head := []byte("package problem_" + data.id + "\n\n")
	perm := os.FileMode(0644)
	err = writeFile(codeFile, perm, head, []byte(data.code))
	if err != nil {
		panic(err)
	}
	err = writeFile(testFile, perm, head, data.tests.ToBytes())
	if err != nil {
		panic(err)
	}
	fmt.Println("create file:", codeFile, testFile)
}

func writeFile(name string, perm os.FileMode, data ...[]byte) error {
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}

	for _, d := range data {
		_, err = f.Write(d)
	}

	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
	}
	return err
}
