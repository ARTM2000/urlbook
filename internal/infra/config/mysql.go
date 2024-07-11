package config

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMysqlDBConn(username, password, host, port, dbName string) *gorm.DB {
	var db *gorm.DB
	var err error

	for {
		db, err = gorm.Open(
			mysql.New(mysql.Config{
				DSN: fmt.Sprintf(
					"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
					username,
					password,
					host,
					port,
					dbName,
				),
			}),
			&gorm.Config{},
		)
	
		if err != nil {
			slog.LogAttrs(
				context.Background(),
				slog.LevelError,
				"failed to connect to mysql database",
				slog.Any("error", err.Error()),
			)
			time.Sleep(time.Second * 5)
			continue
		}
		break
	}

	return db
}
