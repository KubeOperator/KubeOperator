package client

import (
	"errors"
	"gopkg.in/gomail.v2"
	"strconv"
	"strings"
)

var (
	ParamEmpty = "PARAM_EMPTY"
)

type email struct {
	Vars   map[string]interface{}
	Dialer *gomail.Dialer
}

func NewEmailClient(vars map[string]interface{}) (*email, error) {
	var smtp string
	var fromMail string
	var password string
	var port string
	if _, ok := vars["smtp"]; ok {
		smtp = vars["smtp"].(string)
	} else {
		return nil, errors.New(ParamEmpty)
	}
	if _, ok := vars["fromMail"]; ok {
		fromMail = vars["fromMail"].(string)
	} else {
		return nil, errors.New(ParamEmpty)
	}
	if _, ok := vars["port"]; ok {
		port = vars["port"].(string)
	} else {
		return nil, errors.New(ParamEmpty)
	}
	if _, ok := vars["password"]; ok {
		password = vars["password"].(string)
	} else {
		return nil, errors.New(ParamEmpty)
	}
	intPort, err := strconv.Atoi(port)
	if err != nil {
		return nil, errors.New(ParamEmpty)
	}
	d := gomail.NewDialer(smtp, intPort, fromMail, password)
	return &email{
		Vars:   vars,
		Dialer: d,
	}, nil
}

func (e email) SendMessage(vars map[string]interface{}) error {
	toers := []string{}
	var toUsers string
	if _, ok := vars["toUsers"]; ok {
		toUsers = vars["toUsers"].(string)
	} else {
		return errors.New(ParamEmpty)
	}
	var subject string
	if _, ok := vars["subject"]; ok {
		subject = vars["subject"].(string)
	} else {
		return errors.New(ParamEmpty)
	}
	var body string
	if _, ok := vars["body"]; ok {
		body = vars["body"].(string)
	} else {
		return errors.New(ParamEmpty)
	}
	m := gomail.NewMessage()
	for _, tmp := range strings.Split(toUsers, ",") {
		toers = append(toers, strings.TrimSpace(tmp))
	}
	m.SetHeader("To", toers...)
	m.SetHeader("From", e.Vars["fromMail"].(string))
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	err := e.Dialer.DialAndSend(m)
	if err != nil {
		return err
	}
	return nil
}
