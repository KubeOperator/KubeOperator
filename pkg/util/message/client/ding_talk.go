package client

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type dingTalk struct {
	Vars map[string]interface{}
}

func NewDingTalkClient(vars map[string]interface{}) (*dingTalk, error) {
	if _, ok := vars["DING_TALK_WEBHOOK"]; !ok {
		return nil, errors.New(ParamEmpty)
	}
	if _, ok := vars["DING_TALK_SECRET"]; !ok {
		return nil, errors.New(ParamEmpty)
	}
	return &dingTalk{
		Vars: vars,
	}, nil
}

func (d dingTalk) SendMessage(vars map[string]interface{}) error {
	var title string
	if _, ok := vars["TITLE"]; ok {
		title = vars["TITLE"].(string)
	} else {
		return errors.New(ParamEmpty)
	}
	var content string
	if _, ok := vars["CONTENT"]; ok {
		content = vars["CONTENT"].(string)
	} else {
		return errors.New(ParamEmpty)
	}
	var receivers string
	if _, ok := vars["RECEIVERS"]; ok {
		receivers = vars["RECEIVERS"].(string)
	} else {
		return errors.New(ParamEmpty)
	}

	reqBody := make(map[string]interface{})
	reqBody["msgtype"] = "markdown"
	at := make(map[string]string)
	at["atMobiles"] = receivers
	at["isAtAll"] = "False"
	reqBody["at"] = at
	markdown := make(map[string]string)
	markdown["title"] = title
	markdown["text"] = content
	reqBody["markdown"] = markdown
	data, _ := json.Marshal(reqBody)
	body := strings.NewReader(string(data))
	req, err := http.NewRequest(
		http.MethodPost,
		getUrl(d.Vars["DING_TALK_WEBHOOK"].(string), d.Vars["DING_TALK_SECRET"].(string)),
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
		if result["errcode"] != nil && result["errcode"].(float64) == 0 {
			return nil
		} else {
			return errors.New(result["errmsg"].(string))
		}
	}
}

func getUrl(webhook, secret string) string {
	step := "?"
	if strings.Contains(webhook, "?") {
		step = "&"
	}
	params := url.Values{}
	timestamp := time.Now().UnixNano() / 1e6
	b := &[]byte{}
	*b = append(*b, strconv.FormatInt(timestamp, 10)...)
	*b = append(*b, '\n')
	*b = append(*b, secret...)
	h := hmac.New(sha256.New, []byte(secret))
	if _, err := h.Write(*b); err != nil {
		fmt.Printf("getUrl err: %v\n", err)
	}
	sign := base64.StdEncoding.EncodeToString(h.Sum(nil))
	params.Add("timestamp", strconv.FormatInt(timestamp, 10))
	params.Add("sign", sign)
	return strings.Join([]string{webhook, params.Encode()}, step)
}
