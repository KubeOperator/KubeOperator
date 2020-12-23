package ldap

import (
	"fmt"
	"testing"
)

func TestLdap(t *testing.T) {
	vars := make(map[string]string)

	vars["endpoint"] = "172.16.10.141"
	vars["port"] = "389"
	vars["username"] = "cn=Manager,dc=ko,dc=com"
	vars["password"] = ""
	vars["dn"] = "dc=ko,dc=com"
	vars["userFilter"] = "(&(objectClass=organizationalPerson))"

	ldap := NewLdap(vars)
	err := ldap.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	result, err := ldap.Search()
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println(result)
		return
	}
	//err = ldap.Login("zwang","")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}else {
	//	fmt.Println("success")
	//}
}
