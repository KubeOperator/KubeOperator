package captcha

import (
	"errors"
	"strings"

	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/mojocn/base64Captcha"
)

var store = base64Captcha.DefaultMemStore
var verifyCodeFailed = errors.New("VERIFY_CODE_FAILED")

func VerifyCode(codeId string, code string) error {
	if code == "" {
		return verifyCodeFailed
	}
	vv := store.Get(codeId, true)
	vv = strings.TrimSpace(vv)
	code = strings.TrimSpace(code)

	if strings.EqualFold(vv, code) {
		return nil
	} else {
		return verifyCodeFailed
	}
}

func CreateCaptcha() (*dto.Captcha, error) {
	var driverString base64Captcha.DriverString
	driverString.Source = "1234567890QWERTYUIOPLKJHGFDSAZXCVBNMqwertyuioplkjhgfdsazxcvbnm"
	driverString.Width = 120
	driverString.Height = 50
	driverString.NoiseCount = 0
	driverString.Length = 4
	driverString.ShowLineOptions = 6
	driverString.Fonts = []string{"wqy-microhei.ttc"}
	driver := driverString.ConvertFonts()
	c := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err := c.Generate()
	if err != nil {
		return nil, err
	}
	return &dto.Captcha{
		Image:     b64s,
		CaptchaId: id,
	}, nil
}
