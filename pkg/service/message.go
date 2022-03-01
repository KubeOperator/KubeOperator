package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/repository"
)

type MessageService interface {
	SendMessage(mType string, result bool, content string, clusterName string, title string) error
}

type messageService struct {
	messageRepo repository.MessageRepository
}

func NewMessageService() MessageService {
	return &messageService{
		messageRepo: repository.NewMessageRepository(),
	}
}

func (m messageService) SendMessage(mType string, result bool, content string, clusterName string, title string) error {
	log.Debugf("title: %s, cluster name: %s, type: %s, content: %s, result: %v", title, clusterName, mType, content, result)
	return nil
}
