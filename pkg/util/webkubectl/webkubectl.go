package webkubectl

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type response struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

func GetConnectToken(name string, apiServer string, token string) (string, error) {
	req := map[string]string{
		"name":      name,
		"apiServer": apiServer,
		"token":     token,
	}

	j, _ := json.Marshal(&req)

	resp, err := http.Post("http://localhost:8082/api/kube-token", "application/json", bytes.NewReader(j))
	if err != nil {
		return "", err
	}
	var r response
	if resp.StatusCode == http.StatusOK {
		buf, _ := ioutil.ReadAll(resp.Body)
		_ = json.Unmarshal(buf, &r)
		if r.Success {
			return r.Token, nil
		} else {
			return "", errors.New(r.Message)
		}
	}
	return "", errors.New(string(resp.StatusCode))
}
