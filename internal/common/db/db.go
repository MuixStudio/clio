package db

import (
	"fmt"
	"time"

	"github.com/muixstudio/clio/internal/common/db/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func NewDB(conf *config.Config) (*DB, error) {
	if conf.MaxIdleConns == 0 {
		conf.MaxIdleConns = 10
	}
	if conf.MaxOpenConns == 0 {
		conf.MaxOpenConns = 100
	}
	if conf.MaxLifetime == 0 {
		conf.MaxLifetime = 3600
	}
	//timeout := config.DBTimeout
	//if conf.TimeOut > 0 {
	//	timeout = conf.TimeOut
	//}
	//
	//writeTimeout := config.DBWriteTimeout
	//if conf.WriteTimeOut > 0 {
	//	writeTimeout = conf.WriteTimeOut
	//}
	//
	//readTimeout := config.DBReadTimeout
	//if conf.ReadTimeOut > 0 {
	//	readTimeout = conf.ReadTimeOut
	//}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		conf.Host,
		conf.Username,
		conf.Password,
		conf.Database,
		conf.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: NewOrmLog(),
	})
	if err != nil {
		return nil, err
	}
	sdb, err := db.DB()
	if err != nil {
		return nil, err
	}
	sdb.SetMaxIdleConns(conf.MaxIdleConns)
	sdb.SetMaxOpenConns(conf.MaxOpenConns)
	sdb.SetConnMaxLifetime(time.Second * time.Duration(conf.MaxLifetime))

	err = db.Use(NewCustomePlugin())
	if err != nil {
		return nil, err
	}

	return &DB{DB: db}, nil
}

func MustNewDB(conf *config.Config) *DB {
	db, err := NewDB(conf)
	if err != nil {
		panic(err)
	}

	return db
}
