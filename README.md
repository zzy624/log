# Golang 日志库
在 logrus 基础上改造
- 新增日志本地存储，支持灵活切分
- 封装日志调用方式

## 安装

1. go get github.com/sirupsen/logrus
2. go get github.com/zzy624/log

## 使用

```go
package main

import (
	"github.com/zzy624/log"
)

func main() {
	log.SetTopic("TEST_TOPIC")
	log.Info("this is Info")
	log.Info("this is Info KV","key","value")
	log.Error("this is Error")
	log.Info("this is Error KV","key","value")
	log.Debug("this is Debug")
	log.Info("this is Debug KV","key","value")
	log.Warn("this is Warn")
	log.Info("this is Warn KV","key","value")

	var TestStruct = struct {
		Name string `json:"name"`
		Age int `json:"age"`
		Sex string `json:"sex"`
	}{
		Name:"日志测试",
		Age:18,
		Sex:"Male",
	}
	log.Info("this is struct","TestStruct",TestStruct)
}
```

输出

```json
{"topic":"TEST_TOPIC","level":"info","func":"main.main","file":"zzy_log/main.go","line":9,"msg":"this is Info","timestamp":"2019-08-30 15:01:32"}
{"topic":"TEST_TOPIC","level":"info","func":"main.main","file":"zzy_log/main.go","line":10,"msg":"this is Info KV","data":{"key":"value"},"timestamp":"2019-08-30 15:01:32"}
{"topic":"TEST_TOPIC","level":"error","func":"main.main","file":"zzy_log/main.go","line":11,"msg":"this is Error","timestamp":"2019-08-30 15:01:32"}
{"topic":"TEST_TOPIC","level":"info","func":"main.main","file":"zzy_log/main.go","line":12,"msg":"this is Error KV","data":{"key":"value"},"timestamp":"2019-08-30 15:01:32"}
{"topic":"TEST_TOPIC","level":"debug","func":"main.main","file":"zzy_log/main.go","line":13,"msg":"this is Debug","timestamp":"2019-08-30 15:01:32"}
{"topic":"TEST_TOPIC","level":"info","func":"main.main","file":"zzy_log/main.go","line":14,"msg":"this is Debug KV","data":{"key":"value"},"timestamp":"2019-08-30 15:01:32"}
{"topic":"TEST_TOPIC","level":"warning","func":"main.main","file":"zzy_log/main.go","line":15,"msg":"this is Warn","timestamp":"2019-08-30 15:01:32"}
{"topic":"TEST_TOPIC","level":"info","func":"main.main","file":"zzy_log/main.go","line":16,"msg":"this is Warn KV","data":{"key":"value"},"timestamp":"2019-08-30 15:01:32"}
{"topic":"TEST_TOPIC","level":"info","func":"main.main","file":"zzy_log/main.go","line":27,"msg":"this is struct","data":{"TestStruct":{"name":"日志测试","age":18,"sex":"Male"}},"timestamp":"2019-08-30 15:01:32"}

```

> 默认：控制台输出为 TEXT 格式，本地目录下生成以 {Topic}+时间 的文件
日志方法接受两个参数：msg,kv ,kv参数可不传，若传入以 key,value 形式传入,数据会统计放在 data 字段下面

