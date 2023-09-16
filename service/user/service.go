/*
 * Created on 15/09/23 01.56
 *
 * Copyright (c) 2023 Abdul Ghani Abbasi
 */

package user

import (
	"bumn-sembako-be/helper"
	"bumn-sembako-be/model"
	"bumn-sembako-be/request"
	"fmt"
	"gorm.io/gorm"
)

type Service interface {
	Create(user *model.User) (*model.User, error)
	ReadById(id int) (*model.User, error)
	ReadByUsername(username string) (*model.User, error)
	Update(id int, user *request.User) (*model.User, error)
	Delete(id int) error
	Count(criteria map[string]interface{}) int64
	ReadAllBy(criteria map[string]interface{}, search string, page, size int) (*[]model.User, error)
}

type service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return &service{db: db}
}

func (e *service) ReadAllBy(criteria map[string]interface{}, search string, page, size int) (*[]model.User, error) {
	var users []model.User

	query := e.db.Where(criteria)

	if search != "" {
		query.Where("name LIKE ?", search+"%")
	}

	if page == 0 || size == 0 {
		page, size = -1, -1
	}

	limit, offset := helper.GetLimitOffset(page, size)
	err := query.Offset(offset).Order("created_at ASC").Limit(limit).Find(&users).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[user.service.ReadAllBy] error execute query %v \n", err)
		return nil, fmt.Errorf("failed view all data")
	}
	return &users, nil
}

func (e *service) ReadById(id int) (*model.User, error) {
	var user = model.User{}
	err := e.db.Table("users").Where("id = ?", id).First(&user).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[user.service.ReadById] error execute query %v \n", err)
		return nil, fmt.Errorf("id is not exists")
	}
	return &user, nil
}

func (e *service) Create(user *model.User) (*model.User, error) {
	tx := e.db.Begin()
	defer tx.Rollback()

	err := tx.Save(&user).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[user.service.Create] error execute query %v \n", err)
		return nil, fmt.Errorf("failed insert data")
	}

	helper.CommonLogger().Error(err)

	tx.Commit()

	return user, nil
}

func (e *service) Update(id int, user *request.User) (*model.User, error) {
	var upUser = model.User{}
	err := e.db.Table("users").Where("id = ?", id).First(&upUser).Updates(&user).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[user.service.Update] error execute query %v \n", err)
		return nil, fmt.Errorf("failed update data")
	}
	return &upUser, nil
}

func (e *service) ReadByUsername(username string) (*model.User, error) {
	var user = model.User{}
	err := e.db.Table("users").Where("username = ?", username).First(&user).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[user.service.ReadByUsername] error execute query %v \n", err)
		return nil, fmt.Errorf("username is not exists")
	}
	return &user, nil
}

func (e *service) Delete(id int) error {
	var user = model.User{}
	err := e.db.Table("users").Where("id = ?", id).First(&user).Delete(&user).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[user.service.Delete] error execute query %v \n", err)
		return fmt.Errorf("id is not exists")
	}
	return nil
}

func (e *service) Count(criteria map[string]interface{}) int64 {
	var result int64
	err := e.db.Table("users").Where(criteria).Count(&result).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		return 0
	}
	return result
}
