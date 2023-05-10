package leetcodeTests

import (
	"bytes"
	"fmt"
	"github.com/buger/jsonparser"
	"io"
	"log"
	"net/http"
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
	if url == "" {
		panic("URL is empty!")
	}
	if url == "all" {
		CreatAll(dir)
		return
	}
	CreatByUrl(dir, url)
}

const queryAll = `{
  "query": "\n    query problemsetQuestionList($categorySlug: String, $limit: Int, $skip: Int, $filters: QuestionListFilterInput) {\n  problemsetQuestionList(\n    categorySlug: $categorySlug\n    limit: $limit\n    skip: $skip\n    filters: $filters\n  ) {\n    hasMore\n    total\n    questions {\n      acRate\n      difficulty\n      freqBar\n      frontendQuestionId\n      isFavor\n      paidOnly\n      solutionNum\n      status\n      title\n      titleCn\n      titleSlug\n      topicTags {\n        name\n        nameTranslated\n        id\n        slug\n      }\n      extra {\n        hasVideoSolution\n        topCompanyTags {\n          imgUrl\n          slug\n          numSubscribed\n        }\n      }\n    }\n  }\n}\n    ",
  "variables": {
    "categorySlug": "algorithms",
    "skip": %d,
    "limit": 50,
    "filters": {}
  },
  "operationName": "problemsetQuestionList"
}`

func CreatAll(dir string) {
	page := 1
	for {
		data := queryPage(page)
		hasMore, _ := jsonparser.GetBoolean(data, "data", "problemsetQuestionList", "hasMore")
		if !hasMore {
			break
		}
		_, _ = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			id := ""
			defer func() {
				if err := recover(); err != nil {
					fmt.Println("error:", err)
					fmt.Println("url:", "https://leetcode.cn/problems/"+id+"/")
				}
			}()
			id, _ = jsonparser.GetString(value, "titleSlug")
			CreateById(dir, id)
		}, "data", "problemsetQuestionList", "questions")
		page++
	}
}

func queryPage(page int) []byte {
	url := "https://leetcode.cn/graphql/"
	client := &http.Client{}
	body := bytes.NewReader([]byte(fmt.Sprintf(queryAll, (page-1)*50)))
	req, _ := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		panic("请求失败:" + err.Error())
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	data, _ := io.ReadAll(resp.Body)
	return data
}

func CreatByUrl(dir, url string) {
	CreateById(dir, url2id(url))
}

func CreateById(dir, id string) {
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
