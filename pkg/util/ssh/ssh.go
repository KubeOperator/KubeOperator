package ssh

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strings"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/util/hash"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"k8s.io/apimachinery/pkg/util/wait"
)

type SSH struct {
	User        string
	Host        string
	Port        int
	addr        string
	authMethods []ssh.AuthMethod
	dialer      sshDialer
	Retry       int
}

type Config struct {
	User        string `validate:"required"`
	Host        string `validate:"required"`
	Port        int    `validate:"required"`
	Password    string
	PrivateKey  []byte
	PassPhrase  []byte
	DialTimeOut time.Duration
	Retry       int
}

type Interface interface {
	Ping() error
	Exec(cmd ...string) (stdout string, stderr string, exit int, err error)
	Run(cmd ...string) error
	CombinedOutput(cmd ...string) ([]byte, error)
	CopyFile(src, dst string) error
	WriteFile(src io.Reader, dst string) error
	ReadFile(filename string) ([]byte, error)
	Stat(p string) (os.FileInfo, error)
	LookPath(file string) (string, error)
}

func New(c *Config) (*SSH, error) {
	validate := validator.New()
	err := validate.Struct(c)
	if err != nil {
		return nil, err
	}
	if c.Password == "" && c.PrivateKey == nil {
		return nil, errors.New("password or privateKey at least one")
	}

	authMethods := make([]ssh.AuthMethod, 0)

	keyboardInteractiveChallenge := func(
		user,
		instruction string,
		questions []string,
		echos []bool,
	) (answers []string, err error) {
		if len(questions) == 0 {
			return []string{}, nil
		}
		return []string{c.Password}, nil
	}

	if c.Password != "" {
		authMethods = append(authMethods, ssh.Password(c.Password), ssh.KeyboardInteractive(keyboardInteractiveChallenge))
	}
	if len(c.PrivateKey) != 0 {
		signer, err := MakePrivateKeySigner(c.PrivateKey, c.PassPhrase)
		if err != nil {
			return nil, err
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)

	if c.DialTimeOut == 0 {
		c.DialTimeOut = 5 * time.Second
	}

	return &SSH{
		User:        c.User,
		Host:        c.Host,
		Port:        c.Port,
		addr:        addr,
		authMethods: authMethods,
		dialer:      &timeoutDialer{&realSSHDialer{}, c.DialTimeOut},
		Retry:       c.Retry,
	}, nil
}

func (s *SSH) Ping() error {
	_, _, _, err := s.Exec("pwd")

	return err
}

func (s *SSH) CombinedOutput(cmd ...string) ([]byte, error) {
	stdout, stderr, exit, err := s.Exec(cmd...)
	if err != nil {
		return nil, err
	}
	if exit != 0 {
		return nil, fmt.Errorf("exit error %d:%s", exit, stderr)
	}
	return []byte(stdout), nil
}

func (s *SSH) Run(cmd ...string) error {
	_, stderr, exit, err := s.Exec(cmd...)
	if err != nil {
		return err
	}
	if exit != 0 {
		return fmt.Errorf("exit error %d:%s", exit, stderr)
	}
	return nil
}

func (s *SSH) Exec(cmd ...string) (stdout string, stderr string, exit int, err error) {
	config := &ssh.ClientConfig{
		User:            s.User,
		Auth:            s.authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := s.dialer.Dial("tcp", s.addr, config)
	if err != nil && s.Retry > 0 {
		err = wait.Poll(5*time.Second, time.Duration(s.Retry)*5*time.Second, func() (bool, error) {
			if client, err = s.dialer.Dial("tcp", s.addr, config); err != nil {
				return false, err
			}
			return true, nil
		})
	}
	if err != nil {
		return "", "", 0, fmt.Errorf("error getting SSH client to %s@%s: '%v'", s.User, s.addr, err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", "", 0, fmt.Errorf("error creating session to %s@%s: '%v'", s.User, s.addr, err)
	}
	defer session.Close()
	code := 0
	var bout, berr bytes.Buffer
	session.Stdout, session.Stderr = &bout, &berr
	c := strings.Join(cmd, " ")
	if err = session.Run(c); err != nil {
		if exiterr, ok := err.(*ssh.ExitError); ok {
			if code = exiterr.ExitStatus(); code != 0 {
				err = nil
			}
		} else {
			err = fmt.Errorf("failed running `%s` on %s@%s: '%v'", cmd, s.User, s.addr, err)
		}
	}
	return bout.String(), berr.String(), code, err
}

func (s *SSH) CopyFile(src, dst string) error {
	srcHash, err := hash.Sha256WithFile(src)
	if err != nil {
		return err
	}
	hashFile := "/tmp" + dst + ".sha256"
	buffer := new(bytes.Buffer)
	buffer.WriteString(fmt.Sprintf("%s %s", srcHash, dst))
	_ = s.WriteFile(buffer, hashFile)
	_, err = s.CombinedOutput(fmt.Sprintf("sha256sum --check --status %s", hashFile))
	if err == nil { // means dst exist and same as src
		logger.Log.Infof("skip copy `%s` because already existed", src)
		return nil
	}
	logger.Log.Infof("[%s] copy `%s` to %q", s.addr, src, dst)

	config := &ssh.ClientConfig{
		User:            s.User,
		Auth:            s.authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := s.dialer.Dial("tcp", s.addr, config)
	if err != nil {
		err = wait.Poll(5*time.Second, time.Duration(s.Retry)*5*time.Second, func() (bool, error) {
			if client, err = s.dialer.Dial("tcp", s.addr, config); err != nil {
				return false, err
			}
			return true, nil
		})
	}
	if err != nil {
		return err
	}
	defer client.Close()

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	srcFile, err := os.Open(src)

	if err != nil {
		return fmt.Errorf("open file error:%s:%s", src, err)
	}
	defer srcFile.Close()

	_ = sftpClient.MkdirAll(path.Dir(dst))
	dstFile, err := sftpClient.Create(dst)
	if err != nil {
		return fmt.Errorf("create file error:%s:%s", dst, err)
	}
	defer dstFile.Close()

	_, err = dstFile.ReadFrom(srcFile)
	return err
}

func (s *SSH) WriteFile(src io.Reader, dst string) error {
	logger.Log.Infof("[%s] Write data to %q", s.addr, dst)

	config := &ssh.ClientConfig{
		User:            s.User,
		Auth:            s.authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := s.dialer.Dial("tcp", s.addr, config)
	if err != nil {
		err = wait.Poll(5*time.Second, time.Duration(s.Retry)*5*time.Second, func() (bool, error) {
			if client, err = s.dialer.Dial("tcp", s.addr, config); err != nil {
				return false, err
			}
			return true, nil
		})
	}
	if err != nil {
		return err
	}
	defer client.Close()

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	err = sftpClient.MkdirAll(path.Dir(dst))
	if err != nil {
		return err
	}
	dstFile, err := sftpClient.Create(dst)
	if err != nil {
		return fmt.Errorf("create file error:%s:%s", dst, err)
	}
	defer dstFile.Close()

	_, err = dstFile.ReadFrom(src)
	return err
}

func (s *SSH) Stat(p string) (os.FileInfo, error) {
	config := &ssh.ClientConfig{
		User:            s.User,
		Auth:            s.authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := s.dialer.Dial("tcp", s.addr, config)
	if err != nil {
		err = wait.Poll(5*time.Second, time.Duration(s.Retry)*5*time.Second, func() (bool, error) {
			if client, err = s.dialer.Dial("tcp", s.addr, config); err != nil {
				return false, err
			}
			return true, nil
		})
	}
	if err != nil {
		return nil, err
	}
	defer client.Close()

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return nil, err
	}
	defer sftpClient.Close()

	return sftpClient.Stat(p)
}

func (s *SSH) ReadFile(filename string) ([]byte, error) {
	config := &ssh.ClientConfig{
		User:            s.User,
		Auth:            s.authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := s.dialer.Dial("tcp", s.addr, config)
	if err != nil {
		err = wait.Poll(5*time.Second, time.Duration(s.Retry)*5*time.Second, func() (bool, error) {
			if client, err = s.dialer.Dial("tcp", s.addr, config); err != nil {
				return false, err
			}
			return true, nil
		})
	}
	if err != nil {
		return nil, err
	}
	defer client.Close()

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return nil, err
	}
	defer sftpClient.Close()

	f, err := sftpClient.Open(filename)
	if err != nil {
		return nil, err
	}
	data := new(bytes.Buffer)
	_, err = f.WriteTo(data)
	return data.Bytes(), err
}

func (s *SSH) LookPath(file string) (string, error) {
	data, err := s.CombinedOutput(fmt.Sprintf("which %s", file))
	return string(data), err
}

type sshDialer interface {
	Dial(network, addr string, config *ssh.ClientConfig) (*ssh.Client, error)
}

type realSSHDialer struct{}

var _ sshDialer = &realSSHDialer{}

func (d *realSSHDialer) Dial(network, addr string, config *ssh.ClientConfig) (*ssh.Client, error) {
	conn, err := net.DialTimeout(network, addr, config.Timeout)
	if err != nil {
		return nil, err
	}
	_ = conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	c, chans, reqs, err := ssh.NewClientConn(conn, addr, config)
	if err != nil {
		return nil, err
	}
	_ = conn.SetReadDeadline(time.Time{})
	return ssh.NewClient(c, chans, reqs), nil
}

type timeoutDialer struct {
	dialer  sshDialer
	timeout time.Duration
}

func (d *timeoutDialer) Dial(network, addr string, config *ssh.ClientConfig) (*ssh.Client, error) {
	config.Timeout = d.timeout
	return d.dialer.Dial(network, addr, config)
}

func MakePrivateKeySignerFromFile(key string) (ssh.Signer, error) {
	buffer, err := ioutil.ReadFile(key)
	if err != nil {
		return nil, fmt.Errorf("error reading SSH key %s: '%v'", key, err)
	}
	return MakePrivateKeySigner(buffer, nil)
}

func MakePrivateKeySigner(privateKey []byte, passPhrase []byte) (ssh.Signer, error) {
	var signer ssh.Signer
	if passPhrase != nil {
		s, err := ssh.ParsePrivateKeyWithPassphrase(privateKey, passPhrase)
		if err != nil {
			return nil, fmt.Errorf("error parsing SSH key: '%v'", err)
		}
		signer = s
	} else {
		s, err := ssh.ParsePrivateKey(privateKey)
		if err != nil {
			return nil, fmt.Errorf("error parsing SSH key: '%v'", err)
		}
		signer = s
	}

	return signer, nil
}

func ParsePublicKeyFromFile(keyFile string) (*rsa.PublicKey, error) {
	buffer, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("error reading SSH key %s: '%v'", keyFile, err)
	}
	keyBlock, _ := pem.Decode(buffer)
	if keyBlock == nil {
		return nil, fmt.Errorf("error parsing SSH key %s: 'invalid PEM format'", keyFile)
	}
	key, err := x509.ParsePKIXPublicKey(keyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing SSH key %s: '%v'", keyFile, err)
	}
	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("SSH key could not be parsed as rsa public key")
	}
	return rsaKey, nil
}
