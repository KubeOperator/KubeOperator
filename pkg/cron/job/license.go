package job

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"math"
	"time"
)

type LicenseExpire struct {
	licenseService service.LicenseService
	msgService     service.MsgService
}

func NewLicenseExpire() *LicenseExpire {
	return &LicenseExpire{
		licenseService: service.NewLicenseService(),
		msgService:     service.NewMsgService(),
	}
}

func (l *LicenseExpire) Run() {

	detail, err := l.licenseService.Get()
	if err != nil {
		return
	}
	if detail.Status == "invalid" {
		return
	}
	expiration := detail.LicenseInfo.Expired
	t, _ := time.Parse("2006-01-01", expiration)
	now := time.Now()
	subD := now.Sub(t)
	day := int(math.Floor(subD.Hours() / 24))
	if day < 7 {
		err := l.msgService.SendMsg(constant.LicenseExpires, constant.System, nil, true, nil)
		if err != nil {
			logger.Log.Infof("send license msg error,%s", err.Error())
		}
	}
}
