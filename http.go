package leetcodeTests

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	. "leetcodeTests/cache"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var currentDir string
var httpOnce sync.Once

// indexHandle 用于处理/index的GET请求
func indexHandle(context *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			_, _ = fmt.Fprint(context.Writer, err)
		}
	}()
	context.HTML(http.StatusOK, "index.html", nil)
}

// questionHandle 用于处理/question的GET请求。它会获取问题的ID，并从问题数据中获取代码、翻译后的内容和测试用例数据
func questionHandle(c *gin.Context) {
	id := url2id(c.Query("url"))
	cache := NewCache()
	data := cache.Remember("question::"+id, func() interface{} {
		return question(id)
	}, 3600).(*QuestionData)
	c.HTML(http.StatusOK, "question.html", gin.H{
		"id":      id,
		"no":      data.id,
		"code":    data.code,
		"content": template.HTML(data.content),
		"tests":   template.HTML(data.tests.ToCode()),
	})
}

// http初始化
func httpInit() {
	currentDir = getCurrentAbPath()
}

// HttpStart 用于创建HTTP服务器并启动它。它会将HTML模板加载到gin中，并注册两个处理函数Index和Question。
func HttpStart(port int) {
	httpOnce.Do(httpInit)
	router := gin.Default()
	router.LoadHTMLFiles(currentDir+"/templates/index.html", currentDir+"/templates/question.html")
	router.GET("/index", indexHandle)
	router.GET("/question", questionHandle)
	router.StaticFS("/statics", http.Dir(currentDir+"/statics/"))
	_ = router.SetTrustedProxies(nil)
	_ = router.Run(fmt.Sprintf(":%d", port))
}

// 最终方案-全兼容
func getCurrentAbPath() string {
	dir := getCurrentAbPathByExecutable()
	tmpDir, _ := filepath.EvalSymlinks(os.TempDir())
	if strings.Contains(dir, tmpDir) {
		return getCurrentAbPathByCaller()
	}
	return dir
}

// 获取当前执行文件绝对路径
func getCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// 获取当前执行文件绝对路径（go run）
func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}
