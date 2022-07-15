package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type WorkWeiXin struct {
	CorpID     string `json:"corpId"`
	AgentID    string `json:"agentId"`
	CorpSecret string `json:"corpSecret"`
	Token      string `json:"token"`
}

func (w WorkWeiXin) Send(receivers []string, title string, content []byte, others ...string) error {

	if w.Token == "" {
		token, err := w.getToken()
		if err != nil {
			return err
		}
		w.Token = token
	}

	reqBody := make(map[string]interface{})
	reqBody["msgtype"] = "markdown"

	toUser := ""
	for _, v := range receivers {
		toUser = v + "|"
	}
	reqBody["touser"] = toUser
	reqBody["agentid"] = w.AgentID
	markdown := make(map[string]string)
	markdown["content"] = string(content)
	reqBody["markdown"] = markdown
	data, _ := json.Marshal(reqBody)
	body := strings.NewReader(string(data))
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", w.Token),
		body,
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	re, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	} else {
		result := make(map[string]interface{})
		if err := json.Unmarshal([]byte(re), &result); err != nil {
			return err
		}
		if result["errcode"].(float64) == 0 {
			return nil
		} else {
			return errors.New(result["errmsg"].(string))
		}
	}
}

func (w WorkWeiXin) getToken() (string, error) {
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", w.CorpID, w.CorpSecret)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	re, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	} else {
		result := make(map[string]interface{})
		if err := json.Unmarshal([]byte(re), &result); err != nil {
			return "", err
		}
		if result["errcode"].(float64) == 0 {
			return result["access_token"].(string), nil
		} else {
			return "", errors.New(result["errmsg"].(string))
		}
	}
}
