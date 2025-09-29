package svc

import (
	"github.com/muixstudio/clio/internal/common/db"
	dbconfig "github.com/muixstudio/clio/internal/common/db/config"
	"github.com/muixstudio/clio/internal/user/config"
	"github.com/muixstudio/clio/internal/user/models/dao"
)

type ServiceContext struct {
	Config             config.Config
	UserGroupModel     dao.UserGroupModel
	UserModel          dao.UserModel
	UserProfileModel   dao.ProfileModel
	UserUserGroupModel dao.UserUserGroupModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// init db
	dbInstance := db.MustNewDB(dbconfig.Config{
		Username: "postgres",
		Password: "clio2025",
		Host:     "127.0.0.1",
		Port:     5432,
		Database: "clio",
	})
	
	return &ServiceContext{
		Config:    c,
		UserModel: dao.NewUserModel(dbInstance.DB),
	}
}
