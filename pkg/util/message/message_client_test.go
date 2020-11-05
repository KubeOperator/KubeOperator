package message

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	vars := make(map[string]interface{})

	vars["type"] = "DING_TALK"
	vars["DING_TALK_WEBHOOK"] = ""
	vars["DING_TALK_SECRET"] = ""
	client, err := NewMessageClient(vars)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	vars["RECEIVERS"] = "18366136383"
	vars["CONTENT"] = "<font color='info'>本消息由KubeOperator自动发送</font>"
	vars["TITLE"] = "测试消息"
	err = client.SendMessage(vars)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
