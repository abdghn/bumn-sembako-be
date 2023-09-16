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
}

type usecase struct {
	service user.Service
}

func NewUsecase(service user.Service) Usecase {

	return &usecase{service: service}

}

func (u *usecase) Create(user request.User) (*model.User, error) {

	helper.HashPassword(&user.Password)

	newUser := &model.User{
		Name:     user.Name,
		Username: user.Username,
		Password: user.Password,
		Role:     user.Role,
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

	helper.HashPassword(&user.Password)

	newUser := &model.User{
		Name:     user.Name,
		Username: user.Username,
		Password: user.Password,
		Role:     ROLE,
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
	return u.service.ReadByUsername(user.Username)
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
		Name:     user.Name,
		Username: user.Username,
		Password: user.Password,
		Role:     user.Role,
	}
	return u.service.Update(id, updateUser)
}

func (u *usecase) Delete(id int) error {
	return u.service.Delete(id)
}
