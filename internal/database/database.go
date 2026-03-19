package database

import (
	"peer-link-server/config"
	"peer-link-server/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(cfg *config.DatabaseConfig, isProd bool) (*gorm.DB, error) {
	logLevel := logger.Info
	if isProd {
		logLevel = logger.Warn // 生产环境只打印慢查询和错误
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// 自动迁移建表（生产环境推荐用 migrate 工具替代）
	if err := db.AutoMigrate(&model.User{}, &model.Room{}); err != nil {
		return nil, err
	}

	return db, nil
}
