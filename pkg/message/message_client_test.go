package message

import (
	"fmt"
	client2 "github.com/KubeOperator/KubeOperator/pkg/message/client"
	"testing"
)

func Test(t *testing.T) {
	vars := make(map[string]interface{})
	//vars["type"] = "EMAIL"
	//vars["smtp"] = "smtp.sina.com"
	//vars["port"] = "25"
	//vars["fromMail"] = "zhengkunwang123@sina.com"
	//vars["password"] = "Wzk17611460961"
	//client, err := NewMessageClient(vars)
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	//vars["toUsers"] = "zhengkun@fit2cloud.com"
	//vars["subject"] = "测试"
	//vars["body"] = "from:zhengkunwang123@sina.com 测试邮件"
	//err = client.SendMessage(vars)
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}

	//vars["type"] = "DING_TALK"
	//vars["webhook"] = "https://oapi.dingtalk.com/robot/send?access_token=0e960f1b44ed6e5f7cb5874e709cc24d6c37827de9a63ba1b18403f0a47ea87d"
	//vars["secret"] = "SEC047fd6938eec907fdbefa2120639a9a178d1ff769a674c4027e7d41f2ea3dd13"
	//client, err := NewMessageClient(vars)
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	//vars["receivers"] = "13521797236"
	//vars["content"] = "<font color='info'>本消息由KubeOperator自动发送</font>"
	//vars["title"] = "测试消息"
	//err = client.SendMessage(vars)
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}

	vars["type"] = "WORK_WEIXIN"
	vars["corpId"] = "ww918354e3468dc0cc"
	vars["corpSecret"] = "9-5LJqncSKK6xYlYuiTkycz9q5RGpQsVwiFyf68En6M"
	vars["agentId"] = "1000005"
	client, err := NewMessageClient(vars)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	res, err := client2.GetToken(vars)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(res)

	vars["receivers"] = "370785199309261812"
	vars["content"] = "<font color='info'>本消息由KubeOperator自动发送</font>"
	vars["title"] = "测试消息"
	vars["token"] = res
	err = client.SendMessage(vars)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}
