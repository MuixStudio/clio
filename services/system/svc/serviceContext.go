package svc

import (
	"github.com/muixstudio/clio/services/common/db"
	"github.com/muixstudio/clio/services/system/config"
	"github.com/muixstudio/clio/services/system/models/dao"
)

type ServiceContext struct {
	Config            config.Config
	LdapProviderModel dao.LdapProviderModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// init db
	_ = db.MustNewDB(c.DB)

	return &ServiceContext{
		Config: c,
	}
}
