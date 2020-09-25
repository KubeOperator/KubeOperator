package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
)

type MessageService interface {
}

type messageService struct {
	messageRepo repository.MessageRepository
}

func NewMessageService() MessageService {
	return &messageService{
		messageRepo: repository.NewMessageRepository(),
	}
}

func (m messageService) SendMessage(messageType string, success bool, message model.Message, resourceId string) error {

	return nil
}
