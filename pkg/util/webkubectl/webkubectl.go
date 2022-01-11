package webkubectl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
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
	url := fmt.Sprintf("http://%s:%d/api/kube-token", viper.GetString("webkubectl.host"), viper.GetInt("webkubectl.port"))
	resp, err := http.Post(url, "application/json", bytes.NewReader(j))
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
	return "", errors.New(fmt.Sprint(resp.StatusCode))
}
