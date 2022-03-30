package kubepi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var client Interface
var once sync.Once

func GetClient(ops ...Option) Interface {
	once.Do(func() {
		client = NewKubePi()
	})

	client.SetOptions(ops...)
	return client
}

type Interface interface {
	Open(name, apiServer, token string) (Opener, error)
	Close(name, apiServer string) error
	SetOptions(...Option)
}

type KubePi struct {
	Username      string
	Password      string
	Host          string
	Port          int
	sessionCookie *http.Cookie
	mutex         sync.Mutex
}

func (k *KubePi) SetOptions(ops ...Option) {
	for _, o := range ops {
		o(k)
	}
}

type Option func(pi *KubePi)

func WithUsernameAndPassword(username string, password string) Option {
	return func(pi *KubePi) {
		// 如果修改了用户 重新登录
		if pi.Username != username {
			pi.sessionCookie = nil
		}
		pi.Username = username
		pi.Password = password
	}
}

const (
	DefaultKubePiUsername = "admin"
	DefaultKubePiPassword = "kubepi"
)

func NewKubePi() *KubePi {
	kp := &KubePi{
		Username: DefaultKubePiUsername,
		Password: DefaultKubePiPassword,
		Host:     viper.GetString("kubepi.host"),
		Port:     viper.GetInt("kubepi.port"),
	}
	return kp
}

type kubePiCluster struct {
	Name string `json:"name"`
	Spec spec   `json:"spec"`
}

type spec struct {
	Authentication Authentication `json:"authentication" storm:"inline"`
	Connect        connect        `json:"connect"`
}

type Authentication struct {
	Mode        string `json:"mode"`
	BearerToken string `json:"bearerToken"`
}

type connect struct {
	Direction string  `json:"direction"`
	Forward   forward `json:"forward"`
}

type forward struct {
	ApiServer string `json:"apiServer"`
}

type credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ImportConfig struct {
	Name      string `json:"name"`
	ApiServer string `json:"api_server"`
	Token     string `json:"token"`
}

type LoginResponse struct {
	//Data    string `json:"data"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type IsLoginResponse struct {
	Success bool `json:"success"`
	Data    bool `json:"data"`
}

func (k *KubePi) isLogin() (bool, error) {
	if k.sessionCookie == nil {
		return false, nil
	}
	url := fmt.Sprintf("http://%s:%d/kubepi/api/v1/sessions/status", k.Host, k.Port)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}
	req.AddCookie(k.sessionCookie)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	var ils IsLoginResponse
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	if err := json.Unmarshal(bs, &ils); err != nil {
		return false, err
	}
	return ils.Data, nil
}

func (k *KubePi) login() error {
	login, err := k.isLogin()
	if err != nil {
		return err
	}
	if login {
		return nil
	}
	url := fmt.Sprintf("http://%s:%d/kubepi/api/v1/sessions", k.Host, k.Port)
	cred := credential{
		Username: k.Username,
		Password: k.Password,
	}

	js, err := json.Marshal(cred)
	if err != nil {
		return err
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(js))
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var lr LoginResponse
	err = json.Unmarshal(bs, &lr)
	if err != nil {
		return err
	}

	if resp.StatusCode == 200 {
		for i := range resp.Cookies() {
			if resp.Cookies()[i].Name == "SESS_COOKIE_KUBEPI" {
				k.sessionCookie = resp.Cookies()[i]
			}
		}
	} else {
		return fmt.Errorf("login error: %s", lr.Message)
	}
	return nil
}

type ImportResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type ListClustersResponse struct {
	Data    []kubePiCluster `json:"data"`
	Success bool            `json:"success"`
	Message string          `json:"message"`
}

func (k *KubePi) isClusterExists(name, apiServer string) (bool, error) {
	if err := k.login(); err != nil {
		return false, err
	}
	url := fmt.Sprintf("http://%s:%d/kubepi/api/v1/clusters", k.Host, k.Port)
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(k.sessionCookie)
	if err != nil {
		return false, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	var cls ListClustersResponse
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	if err := json.Unmarshal(bs, &cls); err != nil {
		return false, err
	}
	for _, c := range cls.Data {
		if c.Name == name {
			if !strings.HasPrefix(apiServer, "https") {
				apiServer = fmt.Sprintf("https://%s", apiServer)
			}
			if c.Spec.Connect.Forward.ApiServer == apiServer {
				return true, nil
			} else {
				return false, fmt.Errorf("cluster %s already in kubepi, but apiserver is %s", c.Name, c.Spec.Connect.Forward.ApiServer)
			}
		}
	}
	return false, nil
}

func (k *KubePi) ensureImport(name, apiServer, token string) error {
	if err := k.login(); err != nil {
		return err
	}

	exists, err := k.isClusterExists(name, apiServer)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	config := kubePiCluster{Name: name, Spec: spec{
		Authentication: Authentication{
			Mode:        "bearer",
			BearerToken: token,
		},
		Connect: connect{
			Direction: "forward",
			Forward:   forward{ApiServer: apiServer},
		},
	}}

	js, err := json.Marshal(&config)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("http://%s:%d/kubepi/api/v1/clusters", k.Host, k.Port)
	client := http.Client{}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(js))
	request.AddCookie(k.sessionCookie)
	if err != nil {
		return err
	}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	rb, err := ioutil.ReadAll(resp.Body)
	var ir ImportResponse
	if err := json.Unmarshal(rb, &ir); err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return err
	}
	return nil
}

func (k *KubePi) ensureDelete(name, apiServer string) error {
	if err := k.login(); err != nil {
		return err
	}

	exists, _ := k.isClusterExists(name, apiServer)
	if !exists {
		return nil
	}

	url := fmt.Sprintf("http://%s:%d/kubepi/api/v1/clusters/%s", k.Host, k.Port, name)
	client := http.Client{}

	request, err := http.NewRequest(http.MethodDelete, url, nil)
	request.AddCookie(k.sessionCookie)
	if err != nil {
		return err
	}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	rb, err := ioutil.ReadAll(resp.Body)
	var ir ImportResponse
	if err := json.Unmarshal(rb, &ir); err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return err
	}
	return nil
}

type Opener struct {
	SessionCookie *http.Cookie
	Redirect      string
}

func (k *KubePi) Open(name, apiServer, token string) (Opener, error) {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	if err := k.ensureImport(name, apiServer, token); err != nil {
		return Opener{}, err
	}
	url := fmt.Sprintf("/kubepi/dashboard?cluster=%s", name)
	return Opener{SessionCookie: k.sessionCookie, Redirect: url}, nil
}

func (k *KubePi) Close(name, apiServer string) error {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	if err := k.ensureDelete(name, apiServer); err != nil {
		return err
	}
	return nil
}
