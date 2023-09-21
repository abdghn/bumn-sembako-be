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
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Handler interface {
	ViewParticipants(c *gin.Context)
	ViewParticipant(c *gin.Context)
	Update(c *gin.Context)
	ViewDashboard(c *gin.Context)
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

func (h *handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusBadRequest, "id has be number")
		return
	}

	var tempParticipant request.UpdateParticipant

	err = c.ShouldBind(&tempParticipant)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}

	if tempParticipant.Status != "PARTIAL_DONE" && tempParticipant.Status != "" && tempParticipant.File != nil {

		path := "./uploads"
		if _, err := os.Stat(path); os.IsNotExist(err) {
			_ = os.Mkdir(path, os.ModePerm)
		}

		file := tempParticipant.File

		// generate new file name
		ext := filepath.Ext(file.Filename)
		currentTime := time.Now()
		filename := currentTime.Format("20060102150405") + ext

		tmpFile := path + "/" + filename
		if err = c.SaveUploadedFile(file, tmpFile); err != nil {
			helper.HandleError(c, http.StatusBadRequest, "failed to saving image")
			return
		}

		tempParticipant.Image = "image/" + filename

		fmt.Print(tempParticipant.Image)

	}

	updatedParticipant, err := h.usecase.Update(id, tempParticipant)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.HandleSuccess(c, updatedParticipant)
}

func (h *handler) ViewDashboard(c *gin.Context) {
	var req request.ParticipantFilter
	var err error

	err = c.ShouldBindQuery(&req)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}

	result, err := h.usecase.GetTotalDashboard(req)
	helper.HandleSuccess(c, result)
}
