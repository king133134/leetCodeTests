# 自动生成LeetCode Go单元测试用例

用于基于 LeetCode 提供的示例自动生成单元测试用例。该工具目前并不能百分之百地支持所有题目，支持绝大部分数据类型但是不太普遍或没有返回值的问题无法自动生成单元测试。

## 功能

- 基于 LeetCode 提供的示例自动生成 Go 单元测试用例
- 能够自动生成超过 98% 的 LeetCode 题目的单元测试用例
- 不能生成的题目包括环形单向链表和双向链表等

## 使用方法

1. 把代码clone下来或者直接下载代码的压缩包
   ```shell
   git clone git@github.com:king133134/leetCodeTests.git
   ```
2. 进入 leetCodeTests 文件夹 
3. 有两种方式来生成测试用例，只需要提供题目的URL无论是.cn还是.com都支持
   1. 直接生成代码文件（推荐使用），下面命令就会在当前./tests目录生成对应的题目URL的单元测试用例
      ```shell
      go run example/main.go -mod=file -dir=./tests
      ```
      会提示你输入URL，如下图：
      ![示例图片](https://github.com/king133134/leetCodeTests/blob/master/images/1.jpg)
   2. 通过网页表单表单交互生成
      ```shell
      go run example/main.go -mod=http -port=8080
      ```
      然后打开网址[http://localhost:8080/index](http://localhost:8080/index)，输入URL即可

LICENSE
---

[MIT](https://github.com/king133134/leetCodeTests/blob/master/LICENSE)