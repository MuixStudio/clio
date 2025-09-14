package db

import (
	"context"

	"gorm.io/gorm"
)

// TxManagerI 事务管理器接口
type TxManagerI interface {
	ExecTx(context.Context, func(ctx context.Context) error) error
}

// TxManager 基于GORM的事务管理器实现
type TxManager struct {
	db      *gorm.DB
	context context.Context
}

// NewTxManager 创建新的事务管理器
func NewTxManager(db *gorm.DB) TxManagerI {
	return &TxManager{
		db:      db,
		context: context.Background(),
	}
}

// ExecTx 执行事务
func (g *TxManager) ExecTx(fn func(tx *gorm.DB) error) error {
	return g.db.Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

// ExecTxWithContext 执行事务
func (g *TxManager) ExecTxWithContext(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 将事务对象设置到上下文中
		return fn(tx)
	})
}

// GetTx 从上下文中获取事务对象
func GetTx(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		return nil
	}
	return tx
}
