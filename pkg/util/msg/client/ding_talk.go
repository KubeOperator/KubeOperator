package client

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type DingTalk struct {
	WebHook string `json:"webHook"`
	Secret  string `json:"secret"`
}

func (d DingTalk) Send(receivers []string, title string, content []byte, others ...string) error {

	reqBody := make(map[string]interface{})
	reqBody["msgtype"] = "markdown"
	at := make(map[string]string)
	toUser := ""
	for _, v := range receivers {
		toUser = v + "|"
	}
	at["atMobiles"] = toUser
	at["isAtAll"] = "False"
	reqBody["at"] = at
	markdown := make(map[string]string)
	markdown["title"] = title
	markdown["text"] = string(content)
	reqBody["markdown"] = markdown
	data, _ := json.Marshal(reqBody)
	body := strings.NewReader(string(data))
	req, err := http.NewRequest(
		http.MethodPost,
		getUrl(d.WebHook, d.Secret),
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
		logger.Log.Errorf("ding talk get url failed, error: %s", err.Error())
	}
	sign := base64.StdEncoding.EncodeToString(h.Sum(nil))
	params.Add("timestamp", strconv.FormatInt(timestamp, 10))
	params.Add("sign", sign)
	return strings.Join([]string{webhook, params.Encode()}, step)
}
