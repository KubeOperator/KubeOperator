package msg

import (
	"encoding/json"
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/util/msg/client"
)

type MsgClient interface {
	Send(receivers []string, title string, content []byte, others ...string) error
}

func NewMsgClient(name string, config interface{}) (MsgClient, error) {

	configJson, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	switch name {
	case constant.Email:
		var cli client.Email
		if err := json.Unmarshal(configJson, &cli); err != nil {
			return nil, err
		}
		return cli, nil
	case constant.DingTalk:
		var cli client.DingTalk
		if err := json.Unmarshal(configJson, &cli); err != nil {
			return nil, err
		}
		return cli, nil
	case constant.WorkWeiXin:
		var cli client.WorkWeiXin
		if err := json.Unmarshal(configJson, &cli); err != nil {
			return nil, err
		}
		return cli, nil
	default:
		return nil, errors.New("client not support")
	}
}
