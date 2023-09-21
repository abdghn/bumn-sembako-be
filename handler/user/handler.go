/*
 * Created on 15/09/23 01.54
 *
 * Copyright (c) 2023 Abdul Ghani Abbasi
 */

package user

import (
	"bumn-sembako-be/helper"
	"bumn-sembako-be/model"
	"bumn-sembako-be/request"
	"bumn-sembako-be/usecase/user"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Handler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	ViewUsers(c *gin.Context)
	CreateUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	ViewOrganizations(c *gin.Context)
}

type handler struct {
	usecase user.Usecase
}

func NewHandler(usecase user.Usecase) Handler {

	return &handler{usecase: usecase}

}

func (h *handler) Register(c *gin.Context) {
	var user request.Register
	err := c.Bind(&user)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}

	if user.Username == "" {
		helper.HandleError(c, http.StatusBadRequest, "column cannot be empty")
		return
	}

	newUser, err := h.usecase.Register(user)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	newUser.Password = ""
	helper.HandleSuccess(c, newUser)

}

func (h *handler) Login(c *gin.Context) {
	var user = request.Login{}
	err := c.Bind(&user)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}

	dbUser, err := h.usecase.Login(user)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	isTrue := helper.ComparePassword(dbUser.Password, user.Password)
	if !isTrue {

		dbUser.RetryAttempts = dbUser.RetryAttempts + 1
		if dbUser.RetryAttempts > 3 {
			helper.HandleError(c, http.StatusTooManyRequests, "You have retry 3 times")
			return
		}
		_, err = h.usecase.Update(dbUser.ID, dbUser)
		helper.HandleError(c, http.StatusInternalServerError, "Password not matched")
		return
	}

	token := helper.GenerateToken(dbUser)
	result := map[string]interface{}{"token": token, "user": dbUser.Name, "userData": dbUser}
	helper.HandleSuccess(c, result)

}

func (h *handler) ViewUsers(c *gin.Context) {
	var req request.UserPaged
	var err error

	err = c.ShouldBindQuery(&req)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}

	users, err := h.usecase.ReadAllBy(req)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	countUsers := h.usecase.Count(req)

	helper.HandlePagedSuccess(c, users, req.Page, req.Size, countUsers)

}

func (h *handler) CreateUser(c *gin.Context) {
	var user request.User
	err := c.Bind(&user)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}

	if user.Username == "" {
		helper.HandleError(c, http.StatusBadRequest, "column cannot be empty")
		return
	}

	newUser, err := h.usecase.Create(user)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	newUser.Password = ""
	helper.HandleSuccess(c, newUser)

}

func (h *handler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusBadRequest, "id has be number")
		return
	}
	existUser, err := h.usecase.ReadById(id)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusNotFound, err.Error())
		return
	}
	var tempUser = model.User{}
	err = c.Bind(&tempUser)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}
	if tempUser.ID == 0 {
		helper.HandleError(c, http.StatusBadRequest, "input not permitted")
		return
	}

	if tempUser.Password == "" {
		tempUser.Password = existUser.Password
	}

	u, err := h.usecase.Update(id, &tempUser)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.HandleSuccess(c, u)
}

func (h *handler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusBadRequest, "id has be number")
		return
	}
	err = h.usecase.Delete(id)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusNotFound, err.Error())
		return
	}
	helper.HandleSuccess(c, "success delete data")

}

func (h *handler) ViewOrganizations(c *gin.Context) {
	organizations, err := h.usecase.ReadAllOrganization()
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.HandleSuccess(c, organizations)

}
