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
	if _, ok := vars["webhook"]; !ok {
		return nil, errors.New(ParamEmpty)
	}
	if _, ok := vars["secret"]; !ok {
		return nil, errors.New(ParamEmpty)
	}
	return &dingTalk{
		Vars: vars,
	}, nil
}

func (d dingTalk) SendMessage(vars map[string]interface{}) error {
	var title string
	if _, ok := vars["title"]; ok {
		title = vars["title"].(string)
	} else {
		return errors.New(ParamEmpty)
	}
	var content string
	if _, ok := vars["content"]; ok {
		content = vars["content"].(string)
	} else {
		return errors.New(ParamEmpty)
	}
	var receivers string
	if _, ok := vars["receivers"]; ok {
		receivers = vars["receivers"].(string)
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
		getUrl(d.Vars["webhook"].(string), d.Vars["secret"].(string)),
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
		fmt.Println(string(re))
		return nil
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
	h.Write(*b)
	sign := base64.StdEncoding.EncodeToString(h.Sum(nil))
	params.Add("timestamp", strconv.FormatInt(timestamp, 10))
	params.Add("sign", sign)
	return strings.Join([]string{webhook, params.Encode()}, step)
}
