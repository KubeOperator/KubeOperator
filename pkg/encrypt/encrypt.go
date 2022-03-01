package encrypt

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/util/ssh"
	"github.com/spf13/viper"
)

const phaseName = "encrypt"

var log = logger.Default

type InitEncryptPhase struct {
	Multilevel map[string]interface{}
}

func (c *InitEncryptPhase) Init() error {
	enable := c.Multilevel["enable"]
	if enable != nil && enable.(bool) {
		p, err := exec.LookPath(fmt.Sprintf("encrypt_%s_%s", runtime.GOOS, runtime.GOARCH))
		if err != nil {
			return err
		}
		secret, ok := c.Multilevel["secret"].(string)
		if !ok {
			log.Errorf("type aassertion failed")
		}
		parts, ok := c.Multilevel["parts"].([]interface{})
		if !ok {
			log.Errorf("type aassertion failed")
		}
		args := []string{"decrypt", "-t", secret}
		args = append(args, "-p")
		argCommonds := ""
		for i := range parts {
			arg, ok := parts[i].(string)
			if !ok {
				log.Errorf("type aassertion failed")
				continue
			}
			args = append(args, arg)
			argCommonds += (arg + " ")
		}
		if ssh.CheckIllegal(argCommonds) {
			return errors.New("args contains invalid characters!")
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
		_ = strings.TrimPrefix(str, "\n")

		viper.Set("encrypt.key", string(bs))
	}
	return nil
}

func (c *InitEncryptPhase) PhaseName() string {
	return phaseName
}
