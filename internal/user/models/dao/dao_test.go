package dao

import (
	"context"
	"fmt"
	"testing"

	"github.com/muixstudio/clio/internal/common/db"
	dbconfig "github.com/muixstudio/clio/internal/common/db/config"
)

var dbInstance = db.MustNewDB(dbconfig.Config{
	Username: "postgres",
	Password: "clio2025",
	Host:     "127.0.0.1",
	Port:     5432,
	Database: "clio",
})

func TestUserGroupCreate(t *testing.T) {
	usergroup := NewUserModel(dbInstance.DB)
	name := ""
	username := ""
	err := usergroup.Create(context.Background(), &User{
		Name:           &name,
		UserName:       &username,
		AuthProvider:   1,
		AuthProviderID: 2,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserGroupCreateInBatches(t *testing.T) {
	usergroup := NewUserModel(dbInstance.DB)
	err, _ := usergroup.CreateInBatches(context.Background(), []*User{}, 1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserGroupUpdate(t *testing.T) {
	usergroup := NewUserModel(dbInstance.DB)
	err := usergroup.Update(context.Background(), &User{
		BaseModel: BaseModel{
			ID: 711,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserGroupUpdateBatches(t *testing.T) {
	usergroup := NewUserModel(dbInstance.DB)
	err := usergroup.UpdateInBatches(context.Background(), []*User{
		{
			BaseModel: BaseModel{
				ID: 7,
			},
		},
		{
			BaseModel: BaseModel{
				ID: 1,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserGroupDelete(t *testing.T) {
	usergroup := NewUserModel(dbInstance.DB)
	err := usergroup.Delete(context.Background(), &User{
		BaseModel: BaseModel{
			ID: 71,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserGroupDeleteInBatches(t *testing.T) {
	usergroup := NewUserModel(dbInstance.DB)
	err := usergroup.DeleteInBatches(context.Background(), []*User{
		{
			BaseModel: BaseModel{
				ID: 17,
			},
		},
		{
			BaseModel: BaseModel{
				ID: 1,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserGroupFind(t *testing.T) {
	usergroup := NewUserModel(dbInstance.DB)
	re, err := usergroup.Find(context.Background(), &User{}, 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	for _, re := range re {
		fmt.Println(re.ID)
	}
}

func TestUserCount(t *testing.T) {
	usergroup := NewUserModel(dbInstance.DB)
	count, err := usergroup.Count(context.Background(), &User{
		BaseModel: BaseModel{
			ID: 2,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(count)
}
