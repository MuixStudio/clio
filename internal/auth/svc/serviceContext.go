package svc

import (
	"github.com/muixstudio/clio/internal/auth/config"
	"github.com/muixstudio/clio/internal/auth/models/dao"
	"github.com/muixstudio/clio/internal/common/db"
)

type ServiceContext struct {
	Config       config.Config
	SessionModel dao.SessionModel
	TokenModel   dao.TokenModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// init db
	_ = db.MustNewDB(c.DB)

	return &ServiceContext{
		Config: c,
	}
}
