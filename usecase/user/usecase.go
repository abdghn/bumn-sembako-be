/*
 * Created on 15/09/23 01.57
 *
 * Copyright (c) 2023 Abdul Ghani Abbasi
 */

package user

import (
	"bumn-sembako-be/helper"
	"bumn-sembako-be/model"
	"bumn-sembako-be/request"
	"bumn-sembako-be/service/user"
	"fmt"
	"gorm.io/gorm"
)

const ROLE = "STAFF-LAPANGAN"

type Usecase interface {
	Register(user request.Register) (*model.User, error)
	Create(user request.User) (*model.User, error)
	Login(user request.Login) (*model.User, error)
	ReadAllBy(req request.UserPaged) (*[]model.User, error)
	Count(req request.UserPaged) int64
	ReadById(id int) (*model.User, error)
	Update(id int, user *model.User) (*model.User, error)
	Delete(id int) error
	ReadAllOrganization() (*[]model.Organization, error)
}

type usecase struct {
	service user.Service
}

func NewUsecase(service user.Service) Usecase {

	return &usecase{service: service}

}

func (u *usecase) Create(user request.User) (*model.User, error) {
	getUser, err := u.service.CheckByUsername(user.Username)
	if err != nil && err != gorm.ErrRecordNotFound {
		helper.CommonLogger().Error(err)
		return nil, err
	}

	if getUser.ID > 0 {
		return nil, fmt.Errorf("username sudah terdaftar, silahkan menggunakan username lain")
	}

	helper.HashPassword(&user.Password)

	newUser := &model.User{
		Name:           user.Name,
		Username:       user.Username,
		Password:       user.Password,
		Role:           "ADMIN-EO",
		OrganizationID: user.OrganizationID,
		Kota:           user.Kota,
		Provinsi:       user.Provinsi,
		RetryAttempts:  0,
	}

	m, err := u.service.Create(newUser)
	if err != nil {
		helper.CommonLogger().Error(err)
		return nil, err
	}

	m.Password = ""

	return m, nil

}

func (u *usecase) Register(user request.Register) (*model.User, error) {
	getUser, err := u.service.CheckByUsername(user.Username)
	if err != nil && err != gorm.ErrRecordNotFound {
		helper.CommonLogger().Error(err)
		return nil, err
	}

	if getUser.ID > 0 {
		return nil, fmt.Errorf("username sudah terdaftar, silahkan menggunakan username lain")
	}

	helper.HashPassword(&user.Password)

	newUser := &model.User{
		Name:           user.Name,
		Username:       user.Username,
		Password:       user.Password,
		OrganizationID: user.OrganizationID,
		Kota:           user.Kota,
		Provinsi:       user.Provinsi,
		Role:           ROLE,
	}

	m, err := u.service.Create(newUser)
	if err != nil {
		helper.CommonLogger().Error(err)
		return nil, err
	}

	m.Password = ""

	return m, nil

}

func (u *usecase) Login(user request.Login) (*model.User, error) {

	getUser, err := u.service.ReadByUsername(user.Username)
	if err != nil {
		helper.CommonLogger().Error(err)
		return nil, err
	}
	return getUser, nil
}

func (u *usecase) ReadAllBy(req request.UserPaged) (*[]model.User, error) {
	criteria := make(map[string]interface{})

	return u.service.ReadAllBy(criteria, req.Search, req.Page, req.Size)
}

func (u *usecase) Count(req request.UserPaged) int64 {
	criteria := make(map[string]interface{})

	return u.service.Count(criteria)
}

func (u *usecase) ReadById(id int) (*model.User, error) {
	return u.service.ReadById(id)
}

func (u *usecase) Update(id int, user *model.User) (*model.User, error) {
	updateUser := &request.User{
		Name:          user.Name,
		Username:      user.Username,
		Password:      user.Password,
		Role:          user.Role,
		RetryAttempts: user.RetryAttempts,
	}
	return u.service.Update(id, updateUser)
}

func (u *usecase) Delete(id int) error {
	return u.service.Delete(id)
}

func (u *usecase) ReadAllOrganization() (*[]model.Organization, error) {
	return u.service.ReadAllOrganization()
}
