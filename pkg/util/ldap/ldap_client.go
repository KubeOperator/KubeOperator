package ldap

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-ldap/ldap"
)

type Config struct {
	Endpoint string `json:"ldap_address"`
	Port     string `json:"ldap_port"`
	UserName string `json:"ldap_username"`
	UserDn   string `json:"ldap_dn"`
	Password string `json:"ldap_password"`
	Filter   string `json:"ldap_filter"`
	Mapping  string `json:"ldap_mapping"`
	TLS      string `json:"ldap_tls"`
}

func (c *Config) GetAttributes() ([]string, error) {
	m := make(map[string]string)
	err := json.Unmarshal([]byte(c.Mapping), &m)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, v := range m {
		result = append(result, v)
	}
	return result, nil
}

func (c *Config) GetMappings() (map[string]string, error) {
	m := make(map[string]string)
	err := json.Unmarshal([]byte(c.Mapping), &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

type LdapClient struct {
	Vars   map[string]string
	Config Config
	Conn   *ldap.Conn
}

func NewLdap(vars map[string]string) (*LdapClient, error) {
	con, err := json.Marshal(vars)
	if err != nil {
		return nil, err
	}
	config := Config{}
	err = json.Unmarshal(con, &config)
	if err != nil {
		return nil, err
	}
	return &LdapClient{
		Config: config,
	}, nil
}

func (l *LdapClient) Connect() error {
	var err error
	if l.Config.TLS == "ENABLE" {
		l.Conn, err = ldap.DialTLS("tcp", fmt.Sprintf("%s:%s", l.Config.Endpoint, l.Config.Port), &tls.Config{
			InsecureSkipVerify: true,
		})
	} else {
		l.Conn, err = ldap.Dial("tcp", fmt.Sprintf("%s:%s", l.Config.Endpoint, l.Config.Port))
	}
	if err != nil {
		return err
	}
	err = l.Conn.Bind(l.Config.UserName, l.Config.Password)
	if err != nil {
		return err
	}
	return err
}

func (l *LdapClient) Search(dn, filter string, sizeLimit, timeLimit int, attributes []string) ([]*ldap.Entry, error) {

	searchRequest := ldap.NewSearchRequest(dn,
		ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, timeLimit, false,
		filter,
		attributes,
		nil)
	sr, err := l.Conn.SearchWithPaging(searchRequest, uint32(sizeLimit))
	if err != nil {
		return nil, err
	}
	if len(sr.Entries) == 0 {
		return nil, errors.New("LDAP_USER_IS_NULL")
	}
	defer l.Conn.Close()
	return sr.Entries, err
}

func (l *LdapClient) Login(username, password string) error {

	mappings, err := l.Config.GetMappings()
	if err != nil {
		return err
	}
	var userFilter string
	for k, v := range mappings {
		if k == "Name" {
			userFilter = "(" + v + "=" + username + ")"
		}
	}
	searchRequest := ldap.NewSearchRequest(l.Config.UserDn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		userFilter,
		[]string{"dn", "cn", "uid"},
		nil)
	sr, err := l.Conn.Search(searchRequest)
	if err != nil {
		return err
	}
	if len(sr.Entries) != 1 {
		return errors.New("LDAP_LOGIN_USER_IS_NULL")
	}
	userdn := sr.Entries[0].DN
	err = l.Conn.Bind(userdn, password)
	if err != nil {
		return errors.New("PASSWORD_NOT_MATCH")
	}
	defer l.Conn.Close()
	return nil
}
