package encrypt

import (
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const phaseName = "encrypt"

type InitEncryptPhase struct {
	Multilevel map[string]interface{}
}

func (c *InitEncryptPhase) Init() error {
	enable := c.Multilevel["enable"]
	if enable != nil && enable.(bool) {
		p, err := exec.LookPath("ko-encrypt")
		if err != nil {
			return err
		}
		secret := c.Multilevel["secret"].(string)
		parts := c.Multilevel["parts"].([]interface{})
		args := []string{"decrypt", "-t", secret}
		args = append(args, "-p")
		for i := range parts {
			args = append(args, parts[i].(string))
		}
		cmd := exec.Command(p, args...)
		cmd.Env = os.Environ()

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}
		defer stdout.Close()
		if err := cmd.Start(); err != nil {
			return err
		}
		bs, err := ioutil.ReadAll(stdout)
		if err != nil {
			return err
		}
		if err := cmd.Wait(); err != nil {
			return err
		}
		str := string(bs)
		str = strings.TrimPrefix(str, "\n")

		viper.Set("encrypt.key", string(bs))
	}
	return nil
}

func (c *InitEncryptPhase) PhaseName() string {
	return phaseName
}
