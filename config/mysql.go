/*
 * Created on 01/04/22 14.58
 *
 * Copyright (c) 2022 Abdul Ghani Abbasi
 */

package config

import (
	"bumn-sembako-be/helper"
	"bumn-sembako-be/model"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DbConnect(DBUSER, DBPASSWORD, DBHOST, DBPORT, DBNAME string) *gorm.DB {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DBUSER, DBPASSWORD, DBHOST, DBPORT, DBNAME)
	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		helper.CommonLogger().Error("Cannot Load Database: ", err.Error()+"\n")
	}

	db.Debug().AutoMigrate(
		model.User{},
		model.Participant{},
		model.Province{},
		model.Regency{},
		model.District{},
		model.Village{},
	)

	sqlDB, err := db.DB()
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}
