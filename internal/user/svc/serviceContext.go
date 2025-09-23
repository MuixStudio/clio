package svc

import (
	"github.com/muixstudio/clio/internal/common/db"
	"github.com/muixstudio/clio/internal/user/config"
	"github.com/muixstudio/clio/internal/user/models/dao"
)

type ServiceContext struct {
	Config             config.Config
	UserGroupModel     dao.UserGroupModel
	UserModel          dao.UserModel
	UserProfileModel   dao.UserProfileModel
	UserUserGroupModel dao.UserUserGroupModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// init db
	_ = db.MustNewDB(c.DB)

	return &ServiceContext{
		Config: c,
	}
}
