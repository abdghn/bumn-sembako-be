/*
 * Created on 15/09/23 02.32
 *
 * Copyright (c) 2023 Abdul Ghani Abbasi
 */

package participant

import (
	"bumn-sembako-be/helper"
	"bumn-sembako-be/request"
	"bumn-sembako-be/usecase/participant"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Handler interface {
	ViewParticipants(c *gin.Context)
	ViewParticipant(c *gin.Context)
}

type handler struct {
	usecase participant.Usecase
}

func NewHandler(usecase participant.Usecase) Handler {
	return &handler{usecase: usecase}

}

func (h *handler) ViewParticipants(c *gin.Context) {
	var req request.ParticipantPaged
	var err error

	err = c.ShouldBindQuery(&req)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}

	participants, err := h.usecase.ReadAllBy(req)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	countParticipants := h.usecase.Count(req)

	helper.HandlePagedSuccess(c, participants, req.Page, req.Size, countParticipants)

}

func (h *handler) ViewParticipant(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusBadRequest, "id has be number")
		return
	}

	u, err := h.usecase.ReadById(id)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusNotFound, err.Error())
		return
	}
	helper.HandleSuccess(c, u)

}
