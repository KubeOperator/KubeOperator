package service

import (
	"bytes"
	"github.com/KubeOperator/KubeOperator/bindata"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"html/template"
	"io"
)

type MsgService interface {
	GetMsgContent(msgType, sendType string, content map[string]interface{}) (string, error)
}

type msgService struct {
}

func NewMsgService() MsgService {
	return &msgService{}
}

func (m msgService) SendMsg(msgType, resourceName, resourceType string, detail map[string]string, success bool) error {

	return nil
}

func (m msgService) GetMsgContent(msgType, sendType string, content map[string]interface{}) (string, error) {
	tempUrl := dto.Templates[msgType][sendType]
	data, err := bindata.Asset(tempUrl)
	if err != nil {
		return "", err
	}
	newTm := template.New(sendType)
	tm, err := newTm.Parse(string(data))
	if err != nil {
		return "", err
	}
	reader, outStream := io.Pipe()
	go func() {
		err = tm.Execute(outStream, content)
		if err != nil {
			panic(err)
		}
		outStream.Close()
	}()

	buffer := new(bytes.Buffer)
	_, err = buffer.ReadFrom(reader)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}
