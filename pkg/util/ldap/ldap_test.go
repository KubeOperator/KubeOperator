package ldap

import (
	"fmt"
	"testing"
)

func TestLdap(t *testing.T) {
	vars := make(map[string]string)

	vars["ldap_address"] = "172.16.10.89"
	vars["ldap_port"] = "389"
	vars["ldap_username"] = "CN=zhengkun2,CN=Users"
	vars["ldap_password"] = "Calong@2015"
	vars["ldap_dn"] = "DC=ko,DC=com"
	vars["ldap_filter"] = "(&(objectClass=Person))"

	ldap, err := NewLdap(vars)
	if err != nil {
		fmt.Println(err)
	}
	err = ldap.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	//result, err := ldap.Search()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//} else {
	//	fmt.Println(result)
	//}
	err = ldap.Login("zhengkun2", "Calong@2015")
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("success")
	}
}
