package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type SystemSetting struct {
	model.SystemSetting
}

type SystemSettingCreate struct {
	Vars map[string]string `json:"vars" validate:"required"`
	Tab  string            `json:"tab" validate:"required"`
}

type SystemSettingUpdate struct {
	Vars map[string]string `json:"vars" validate:"required"`
	Tab  string            `json:"tab" validate:"required"`
}

type SystemSettingResult struct {
	Vars map[string]string `json:"vars" validate:"required"`
	Tab  string            `json:"tab" validate:"required"`
}

type LdapResult struct {
	Data interface{} `json:"data"`
}

type LdapSetting struct {
	Endpoint  string `json:"ldap_address"`
	Port      string `json:"ldap_port"`
	UserName  string `json:"ldap_username"`
	UserDn    string `json:"ldap_dn"`
	Password  string `json:"ldap_password"`
	Filter    string `json:"ldap_filter"`
	Mapping   string `json:"ldap_mapping"`
	Status    string `json:"ldap_status"`
	TLS       string `json:"ldap_tls"`
	SizeLimit int
	TimeLimit int
}
