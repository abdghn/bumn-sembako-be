/*
 * Created on 15/09/23 16.41
 *
 * Copyright (c) 2023 Abdul Ghani Abbasi
 */

package region

import (
	"bumn-sembako-be/helper"
	"bumn-sembako-be/request"
	"bumn-sembako-be/usecase/region"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler interface {
	ViewProvincies(c *gin.Context)
	ViewRegenciesByProvinceId(c *gin.Context)
	ViewDistrictsByRegencyId(c *gin.Context)
	ViewVillagesByDistrictId(c *gin.Context)
}

type handler struct {
	usecase region.Usecase
}

func NewHandler(usecase region.Usecase) Handler {
	return &handler{usecase: usecase}
}

func (h *handler) ViewProvincies(c *gin.Context) {
	var req request.RegionPaged
	var err error

	err = c.ShouldBindQuery(&req)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}

	provincies, err := h.usecase.ReadAllProvinceBy(req)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.HandleSuccess(c, provincies)

}

func (h *handler) ViewRegenciesByProvinceId(c *gin.Context) {
	var req request.RegionPaged
	var err error

	err = c.ShouldBindQuery(&req)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}

	regencies, err := h.usecase.ReadAllRegencyBy(req)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.HandleSuccess(c, regencies)
}

func (h *handler) ViewDistrictsByRegencyId(c *gin.Context) {
	var req request.RegionPaged
	var err error

	err = c.ShouldBindQuery(&req)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}

	districts, err := h.usecase.ReadAllDistrictBy(req)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.HandleSuccess(c, districts)

}

func (h *handler) ViewVillagesByDistrictId(c *gin.Context) {
	var req request.RegionPaged
	var err error

	err = c.ShouldBindQuery(&req)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}

	villages, err := h.usecase.ReadAllVillageBy(req)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.HandleSuccess(c, villages)

}
