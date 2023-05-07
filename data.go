package leetcodeTests

import (
	"bytes"
	"fmt"
	"github.com/buger/jsonparser"
	"io"
	. "leetcodeTests/parser"
	"net/http"
	"strings"
)

/**
{"operationName":"question","variables":{"titleSlug":"largest-local-values-in-a-matrix"},"query":"query question($titleSlug: String!) {\n  question(titleSlug: $titleSlug) {\n    questionId\n    questionFrontendId\n    categoryTitle\n    boundTopicId\n    title\n    titleSlug\n    content\n    translatedTitle\n    translatedContent\n    isPaidOnly\n    difficulty\n    likes\n    dislikes\n    isLiked\n    similarQuestions\n    contributors {\n      username\n      profileUrl\n      avatarUrl\n      __typename\n    }\n    langToValidPlayground\n    topicTags {\n      name\n      slug\n      translatedName\n      __typename\n    }\n    companyTagStats\n    codeSnippets {\n      lang\n      langSlug\n      code\n      __typename\n    }\n    stats\n    hints\n    solution {\n      id\n      canSeeDetail\n      __typename\n    }\n    status\n    sampleTestCase\n    metaData\n    judgerAvailable\n    judgeType\n    mysqlSchemas\n    enableRunCode\n    envInfo\n    book {\n      id\n      bookName\n      pressName\n      source\n      shortDescription\n      fullDescription\n      bookImgUrl\n      pressImgUrl\n      productUrl\n      __typename\n    }\n    isSubscribed\n    isDailyQuestion\n    dailyRecordStatus\n    editorType\n    ugcQuestionId\n    style\n    exampleTestcases\n    jsonExampleTestcases\n    __typename\n  }\n}\n"}
*/

// GraphQL查询模板
const queryPattern = "{\"operationName\":\"question\",\"variables\":{\"titleSlug\":\"$id\"},\"query\":\"query question($titleSlug: String!) {\\n  question(titleSlug: $titleSlug) {\\n    questionId\\n    questionFrontendId\\n    categoryTitle\\n    boundTopicId\\n    title\\n    titleSlug\\n    difficulty\\n    translatedContent\\n    codeSnippets {\\n      lang\\n      langSlug\\n      code\\n      __typename\\n    }\\n    content\\n  }\\n}\\n\"}"

// 根据ID获取GraphQL
func getRequestBody(id string) (res []byte) {
	return []byte(strings.Replace(queryPattern, "$id", id, 1))
}

type QuestionData struct {
	id      string // 题目编号
	code    string // 代码片段
	tests   *Tests // 测试用例
	content string // 题目内容
}

// question 获取题目编号，代码片段，测试用例，题目内容
func question(id string) *QuestionData {
	url := "https://leetcode.cn/graphql/"
	client := &http.Client{}
	body := bytes.NewReader(getRequestBody(id))
	req, _ := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败:", err.Error())
		return nil
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	data, _ := io.ReadAll(resp.Body)
	content, _ := jsonparser.GetString(data, "data", "question", "content")
	if content == "" {
		panic("content is empty.")
	}
	translatedContent, _ := jsonparser.GetString(data, "data", "question", "translatedContent")
	questionFrontendId, _ := jsonparser.GetString(data, "data", "question", "questionFrontendId")
	code := ""
	_, _ = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if lang, _ := jsonparser.GetString(value, "lang"); lang == "Go" {
			code, _ = jsonparser.GetString(value, "code")
		}
	}, "data", "question", "codeSnippets")
	return &QuestionData{tests: CreateTests(code, &content), code: code, content: translatedContent, id: questionFrontendId}
}

func url2id(url string) (id string) {
	if len(url) == 0 {
		panic("problem url is empty")
	}
	fmt.Println("get content url", url)
	arr := strings.Split(url, "/problems/")
	return strings.Trim(arr[1], "/")
}

func Question(url string) *QuestionData {
	return question(url2id(url))
}
