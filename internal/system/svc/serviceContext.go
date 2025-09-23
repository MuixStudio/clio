package svc

import (
	"github.com/muixstudio/clio/internal/common/db"
	"github.com/muixstudio/clio/internal/system/config"
	"github.com/muixstudio/clio/internal/system/models/dao"
)

type ServiceContext struct {
	Config             config.Config
	LdapProviderModel         dao.LdapProviderModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// init db
	_ = db.MustNewDB(c.DB)

	return &ServiceContext{
		Config: c,
	}
}
