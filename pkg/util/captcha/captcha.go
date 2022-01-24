package captcha

import (
	"errors"

	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/mojocn/base64Captcha"
)

var store = base64Captcha.DefaultMemStore
var verifyCodeFailed = errors.New("VERIFY_CODE_FAILED")

func VerifyCode(codeId string, code string) error {
	if code == "" {
		return verifyCodeFailed
	}
	if store.Verify(codeId, code, true) {
		return nil
	} else {
		return verifyCodeFailed
	}
}

func CreateCaptcha() (*dto.Captcha, error) {
	var driverString base64Captcha.DriverString
	driverString.Source = "1234567890qwertyuioplkjhgfdsazxcvbnmQWERTYUIOPLKJHGFDSAZXCVBNM"
	driverString.Width = 120
	driverString.Height = 50
	driverString.NoiseCount = 5
	driverString.Length = 4
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
