package main

import (
	"github.com/muixstudio/clio/internal/common/db"
	dbconfig "github.com/muixstudio/clio/internal/common/db/config"
	"github.com/muixstudio/clio/internal/user/internal/models/models"
)

func main() {
	//var configFile = flag.String("f", "../../etc/user.yaml", "the config file")
	//var c config.Config
	//conf.MustLoad(*configFile, &c)
	dbInstance := db.MustNewDB(&dbconfig.Config{
		Username: "postgres",
		Password: "clio2025",
		Host:     "127.0.0.1",
		Port:     5432,
		Database: "clio",
		//Charset:      c.DB.Charset,
		//MaxOpenConns: c.DB.MaxOpenConns,
		//MaxLifetime:  c.DB.MaxLifetime,
		//TimeOut:      c.DB.TimeOut,
		//WriteTimeOut: c.DB.WriteTimeOut,
		//ReadTimeOut:  c.DB.ReadTimeOut,
	})
	err := dbInstance.AutoMigrate(
		&models.User{},
		&models.Token{},
		&models.UserProfile{},
		&models.UserGroup{},
		&models.UserUserGroup{},
		&models.LdapProvider{},
	)
	if err != nil {
		panic(err)
	}
}
