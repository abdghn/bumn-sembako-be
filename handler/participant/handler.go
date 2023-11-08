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
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	ViewParticipants(c *gin.Context)
	ViewLogs(c *gin.Context)
	ViewParticipant(c *gin.Context)
	Update(c *gin.Context)
	ViewDashboard(c *gin.Context)
	ExportReport(c *gin.Context)
	BulkCreate(c *gin.Context)
	ExportExcel(c *gin.Context)
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

func (h *handler) ViewLogs(c *gin.Context) {
	var req request.ParticipantPaged
	var err error

	err = c.ShouldBindQuery(&req)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}

	participants, err := h.usecase.ReadAllLogBy(req)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	countParticipants := h.usecase.CountLogs(req)

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
		filename := currentTime.Format("20060102150405") + "-image" + ext

		tmpFile := path + "/" + filename
		if err = c.SaveUploadedFile(file, tmpFile); err != nil {
			helper.HandleError(c, http.StatusBadRequest, "failed to saving image")
			return
		}

		tempParticipant.Image = "image/" + filename

	}

	if tempParticipant.Status != "PARTIAL_DONE" && tempParticipant.Status != "" && tempParticipant.FilePenerima != nil {

		path := "./uploads"
		if _, err := os.Stat(path); os.IsNotExist(err) {
			_ = os.Mkdir(path, os.ModePerm)
		}

		file := tempParticipant.FilePenerima

		// generate new file name
		ext := filepath.Ext(file.Filename)
		currentTime := time.Now()
		filename := currentTime.Format("20060102150405") + "-image-penerima" + ext

		tmpFile := path + "/" + filename
		if err = c.SaveUploadedFile(file, tmpFile); err != nil {
			helper.HandleError(c, http.StatusBadRequest, "failed to saving image")
			return
		}

		tempParticipant.ImagePenerima = "image/" + filename

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

	result, err := h.usecase.GetTotalDashboardV2(req)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusNotFound, err.Error())
		return
	}
	helper.HandleSuccess(c, result)
}

func (h *handler) ExportReport(c *gin.Context) {
	var req request.Report
	var err error

	err = c.ShouldBindJSON(&req)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}

	req.Url = fmt.Sprintf("http://%s", c.Request.Host)

	path, err := h.usecase.Export(req)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	//pdfg, err := wkhtmltopdf.NewPDFGenerator()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//f, err := os.Open("./report.html")
	//if f != nil {
	//	defer f.Close()
	//}
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//pdfg.AddPage(wkhtmltopdf.NewPageReader(f))
	//
	//pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	//pdfg.Dpi.Set(300)
	//
	//err = pdfg.Create()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//err = pdfg.WriteFile("./output.pdf")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//log.Println("Done")

	helper.HandleSuccess(c, path)

}

func (h *handler) BulkCreate(c *gin.Context) {
	var req request.ImportParticipant
	var err error
	err = c.ShouldBind(&req)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}

	path := "./uploads"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.Mkdir(path, os.ModePerm)
	}

	file := req.File

	// generate new file name
	ext := filepath.Ext(file.Filename)
	currentTime := time.Now()
	filename := currentTime.Format("20060102150405") + ext

	tmpFile := path + "/" + filename
	if err = c.SaveUploadedFile(file, tmpFile); err != nil {
		helper.HandleError(c, http.StatusBadRequest, "failed to saving image")
		return
	}

	req.Name = file.Filename
	req.TmpPath = tmpFile

	result, err := h.usecase.BulkCreate(req)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.HandleSuccess(c, result)

}

func (h *handler) ExportExcel(c *gin.Context) {
	var req request.ParticipantFilter
	var err error

	err = c.ShouldBindQuery(&req)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}

	result, err := h.usecase.ExportExcel(req)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusNotFound, err.Error())
		return
	}
	helper.HandleSuccess(c, result)

}
