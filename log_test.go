package log

import "testing"


func TestInfo(t *testing.T) {
	SetTopic("test_topic")
	Info("this is Info")
	Error("this is Error")
	Debug("this is Debug")
	Warn("this is Warn")

	var Teststruct = struct {
		Name string `json:"name"`
		Age int `json:"age"`
		Sex string `json:"sex"`
	}{
		Name:"日志测试",
		Age:18,
		Sex:"Male",
	}
	Info("this is struct","Teststruct",Teststruct)
}
