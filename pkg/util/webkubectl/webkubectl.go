package webkubectl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
)

type response struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

func GetConnectToken(name string, apiServer string, token []byte) (string, error) {
	req := map[string]string{
		"name":      name,
		"apiServer": apiServer,
		"token":     string(token),
	}

	j, err := json.Marshal(&req)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("http://%s:%d/api/kube-token", viper.GetString("webkubectl.host"), viper.GetInt("webkubectl.port"))
	resp, err := http.Post(url, "application/json", bytes.NewReader(j))
	if err != nil {
		return "", err
	}
	var r response
	if resp.StatusCode == http.StatusOK {
		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		if err := json.Unmarshal(buf, &r); err != nil {
			return "", err
		}
		if r.Success {
			return r.Token, nil
		} else {
			return "", errors.New(r.Message)
		}
	}
	return "", errors.New(fmt.Sprint(resp.StatusCode))
}
