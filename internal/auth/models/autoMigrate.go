package main

import (
	"github.com/muixstudio/clio/internal/common/db"
	dbconfig "github.com/muixstudio/clio/internal/common/db/config"
	dao2 "github.com/muixstudio/clio/internal/user/models/dao"
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
		&dao2.User{},
		&dao2.Token{},
		&dao2.UserProfile{},
		&dao2.UserGroup{},
		&dao2.UserUserGroup{},
		&dao2.LdapProvider{},
	)
	if err != nil {
		panic(err)
	}
}
