package main

import (
	"github.com/muixstudio/clio/internal/common/db"
	dbconfig "github.com/muixstudio/clio/internal/common/db/config"
	"github.com/muixstudio/clio/internal/user/models/dao"
)

func main() {
	dbInstance := db.MustNewDB(dbconfig.Config{
		Username: "postgres",
		Password: "clio2025",
		Host:     "127.0.0.1",
		Port:     5432,
		Database: "clio",
	})
	err := dbInstance.AutoMigrate(
		&dao.User{},
		&dao.Profile{},
		&dao.UserGroup{},
		&dao.UserUserGroup{},
	)
	if err != nil {
		panic(err)
	}
}
