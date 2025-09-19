package models

const (
	Users        = 0
	LdapDisabled = 1
)

type UserStatusCode uint16

var UserStatusCodeMap = map[uint16]string{
	Users:        "Normal User",
	LdapDisabled: "LDAP User Disabled",
}

func (u UserStatusCode) Desc() string {
	return UserStatusCodeMap[uint16(u)]
}
