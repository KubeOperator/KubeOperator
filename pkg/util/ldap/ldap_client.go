package ldap

import (
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
	conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%s", l.Config.Endpoint, l.Config.Port))
	if err != nil {
		return err
	}
	err = conn.Bind(l.Config.UserName, l.Config.Password)
	if err != nil {
		return err
	}
	l.Conn = conn
	return err
}

func (l *LdapClient) Search(attributes []string) ([]*ldap.Entry, error) {

	searchRequest := ldap.NewSearchRequest(l.Config.UserDn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		l.Config.Filter,
		attributes,
		nil)
	sr, err := l.Conn.Search(searchRequest)
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
