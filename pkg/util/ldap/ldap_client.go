package ldap

import (
	"errors"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/util/escape"
	"github.com/go-ldap/ldap"
)

var (
	ParamEmpty = "PARAM_EMPTY"
)

type LdapClient struct {
	Vars map[string]string
	Conn *ldap.Conn
}

func NewLdap(vars map[string]string) *LdapClient {
	return &LdapClient{
		Vars: vars,
	}
}

func (l *LdapClient) Connect() error {
	var endpoint string
	var port string
	var username string
	var password []byte
	if _, ok := l.Vars["ldap_address"]; ok {
		endpoint = l.Vars["ldap_address"]
	} else {
		return errors.New(ParamEmpty)
	}
	if _, ok := l.Vars["ldap_port"]; ok {
		port = l.Vars["ldap_port"]
	} else {
		return errors.New(ParamEmpty)
	}
	if _, ok := l.Vars["ldap_username"]; ok {
		username = l.Vars["ldap_username"]
	} else {
		return errors.New(ParamEmpty)
	}
	password = escape.GetByte(l.Vars["ldap_password"])
	conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%s", endpoint, port))
	if err != nil {
		return err
	}
	err = conn.Bind(username, string(password))
	if err != nil {
		return err
	}
	l.Conn = conn
	escape.Clean(string(password))
	return err
}

func (l *LdapClient) Search() ([]*ldap.Entry, error) {
	var dn string
	if _, ok := l.Vars["ldap_dn"]; ok {
		dn = l.Vars["ldap_dn"]
	} else {
		return nil, errors.New(ParamEmpty)
	}
	var userFilter string
	if _, ok := l.Vars["ldap_filter"]; ok {
		userFilter = l.Vars["ldap_filter"]
	} else {
		return nil, errors.New(ParamEmpty)
	}

	searchRequest := ldap.NewSearchRequest(dn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		userFilter,
		[]string{"cn", "mail"},
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

func (l *LdapClient) Login(userName string, password []byte) error {
	var dn string
	if _, ok := l.Vars["ldap_dn"]; ok {
		dn = l.Vars["ldap_dn"]
	} else {
		return errors.New(ParamEmpty)
	}
	searchRequest := ldap.NewSearchRequest(dn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(cn=%s))", userName),
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
	err = l.Conn.Bind(userdn, string(password))
	if err != nil {
		return errors.New("PASSWORD_NOT_MATCH")
	}
	defer l.Conn.Close()
	escape.Clean(string(password))
	return nil
}
