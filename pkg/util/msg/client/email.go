package client

import "gopkg.in/gomail.v2"

type Email struct {
	Address  string `json:"address"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (e Email) Send(receivers []string, title string, content []byte, others ...string) error {
	d := gomail.NewDialer(e.Address, e.Port, e.Username, e.Password)
	m := gomail.NewMessage()
	m.SetHeader("To", receivers...)
	m.SetHeader("From", e.Username)
	m.SetHeader("Subject", title)
	m.SetBody("text/html", string(content))
	return d.DialAndSend(m)
}
