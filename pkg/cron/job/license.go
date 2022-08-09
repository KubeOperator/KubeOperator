package job

import (
	"fmt"
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
	subD := t.Sub(now)
	day := int(math.Floor(subD.Hours() / 24))
	if day > 0 && (day == 7 || day == 30 || day == 3 || day == 15) {
		message := map[string]string{}
		message["message"] = fmt.Sprintf("License还有%d天到期，请及时处理", day)
		err := l.msgService.SendMsg(constant.LicenseExpires, constant.System, map[string]string{"name": "license"}, true, message)
		if err != nil {
			logger.Log.Infof("send license msg error,%s", err.Error())
		}
	}
}
