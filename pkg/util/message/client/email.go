package client

import (
	"errors"
	"strconv"
	"strings"

	"gopkg.in/gomail.v2"
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
	if _, ok := vars["SMTP_ADDRESS"]; ok {
		smtp = vars["SMTP_ADDRESS"].(string)
	} else {
		return nil, errors.New(ParamEmpty)
	}
	if _, ok := vars["SMTP_USERNAME"]; ok {
		fromMail = vars["SMTP_USERNAME"].(string)
	} else {
		return nil, errors.New(ParamEmpty)
	}
	if _, ok := vars["SMTP_PORT"]; ok {
		port = vars["SMTP_PORT"].(string)
	} else {
		return nil, errors.New(ParamEmpty)
	}
	if _, ok := vars["SMTP_PASSWORD"]; ok {
		password = vars["SMTP_PASSWORD"].(string)
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
	if _, ok := vars["RECEIVERS"]; ok {
		toUsers = vars["RECEIVERS"].(string)
	} else {
		return errors.New(ParamEmpty)
	}
	var subject string
	if _, ok := vars["TITLE"]; ok {
		subject = vars["TITLE"].(string)
	} else {
		return errors.New(ParamEmpty)
	}
	var body string
	if _, ok := vars["CONTENT"]; ok {
		body = vars["CONTENT"].(string)
	} else {
		return errors.New(ParamEmpty)
	}
	m := gomail.NewMessage()
	for _, tmp := range strings.Split(toUsers, ",") {
		toers = append(toers, strings.TrimSpace(tmp))
	}
	m.SetHeader("To", toers...)
	m.SetHeader("From", e.Vars["SMTP_USERNAME"].(string))
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	err := e.Dialer.DialAndSend(m)
	if err != nil {
		return err
	}
	return nil
}
