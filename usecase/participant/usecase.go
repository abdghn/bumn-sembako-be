/*
 * Created on 15/09/23 02.31
 *
 * Copyright (c) 2023 Abdul Ghani Abbasi
 */

package participant

import (
	"bumn-sembako-be/helper"
	"bumn-sembako-be/model"
	"bumn-sembako-be/request"
	"bumn-sembako-be/service/participant"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Usecase interface {
	ReadAllBy(req request.ParticipantPaged) (*[]model.Participant, error)
	Count(req request.ParticipantPaged) int64
	ReadById(id int) (*model.Participant, error)
	Update(id int, input request.UpdateParticipant) (*model.Participant, error)
	GetTotalDashboard(req request.ParticipantFilter) (*model.TotalParticipantResponse, error)
	Export(input request.Report) (string, error)
}

type usecase struct {
	service participant.Service
}

func NewUsecase(service participant.Service) Usecase {
	return &usecase{service: service}

}

func (u *usecase) ReadAllBy(req request.ParticipantPaged) (*[]model.Participant, error) {
	criteria := make(map[string]interface{})
	if req.Provinsi != "" {
		criteria["provinsi"] = req.Provinsi
	}

	if req.Kota != "" {
		criteria["kota"] = req.Kota
	}

	if req.Kecamatan != "" {
		criteria["kecamatan"] = req.Kecamatan
	}

	if req.Kelurahan != "" {
		criteria["kelurahan"] = req.Kelurahan
	}

	if req.Status != "" {
		criteria["status"] = req.Status
	}

	return u.service.ReadAllBy(criteria, req.Search, req.Page, req.Size)
}

func (u *usecase) ReadById(id int) (*model.Participant, error) {
	return u.service.ReadById(id)
}

func (u *usecase) Count(req request.ParticipantPaged) int64 {
	criteria := make(map[string]interface{})
	if req.Provinsi != "" {
		criteria["provinsi"] = req.Provinsi
	}

	if req.Kota != "" {
		criteria["kota"] = req.Kota
	}

	if req.Kecamatan != "" {
		criteria["kecamatan"] = req.Kecamatan
	}

	if req.Kelurahan != "" {
		criteria["kelurahan"] = req.Kelurahan
	}

	if req.Status != "" {
		criteria["status"] = req.Status
	}

	return u.service.Count(criteria)
}

func (u *usecase) Update(id int, input request.UpdateParticipant) (*model.Participant, error) {
	var err error
	participant, err := u.ReadById(id)
	if err != nil {
		return nil, err
	}

	if input.Status == "PARTIAL_DONE" {
		req := &request.PartialDone{Status: "PARTIAL_DONE", UpdatedBy: input.UpdatedBy}

		updatedParticipant, err := u.service.UpdateStatus(participant.ID, req)
		if err != nil {
			return nil, err
		}

		return updatedParticipant, nil

	} else if input.Status == "REJECTED" {
		req := &request.PartialDone{Status: "REJECTED", UpdatedBy: input.UpdatedBy}

		_, err = u.service.UpdateStatus(id, req)
		if err != nil {
			return nil, err
		}

		m := &model.Participant{
			Name:               input.Name,
			NIK:                input.NIK,
			Gender:             input.Gender,
			Phone:              input.Phone,
			Address:            input.Address,
			RT:                 input.RT,
			RW:                 input.RW,
			Provinsi:           input.Provinsi,
			Kota:               input.Kota,
			Kecamatan:          input.Kecamatan,
			Kelurahan:          input.Kelurahan,
			KodePOS:            input.KodePOS,
			ResidenceAddress:   input.ResidenceAddress,
			ResidenceRT:        input.ResidenceRT,
			ResidenceRW:        input.ResidenceRW,
			ResidenceProvinsi:  input.ResidenceProvinsi,
			ResidenceKota:      input.ResidenceKota,
			ResidenceKecamatan: input.ResidenceKecamatan,
			ResidenceKelurahan: input.ResidenceKelurahan,
			ResidenceKodePOS:   input.ResidenceKodePOS,
			Image:              input.Image,
			ImagePenerima:      input.ImagePenerima,
			Status:             "DONE",
			UpdatedBy:          input.UpdatedBy,
		}

		newParticipant, err := u.service.Create(m)
		if err != nil {
			return nil, err
		}

		return newParticipant, nil

	} else if input.Status == "DONE" {
		req := &request.PartialDone{Status: "DONE", Image: input.Image, ImagePenerima: input.ImagePenerima, UpdatedBy: input.UpdatedBy}

		updateParticipant, err := u.service.UpdateStatus(id, req)
		if err != nil {
			return nil, err
		}

		return updateParticipant, nil

	}

	return participant, nil

}

func (u *usecase) GetTotalDashboard(req request.ParticipantFilter) (*model.TotalParticipantResponse, error) {
	var m model.TotalParticipantResponse
	//var date time.Time
	var startDate, endDate time.Time
	status := ""
	quota := int64(0)

	criteria := make(map[string]interface{})
	if req.Provinsi != "" {
		criteria["provinsi"] = req.Provinsi
	}

	if req.Kota != "" {
		criteria["kota"] = req.Kota
	}

	dataQuota, err := u.service.GetQuota(criteria)
	if err != nil {
		log.Println(err)
		//return nil, err
	} else {
		quota = dataQuota.Total
	}

	if req.Date != "" {
		//stringDate := req.Date + "T00:00:00.00Z"
		//date, err = time.Parse(time.RFC3339, stringDate)
		//if err != nil {
		//	return nil, err
		//}

		stringStartDate := req.Date + "T00:00:00.00Z"
		stringEndDate := req.Date + "T23:59:59.999:Z"
		startDate, err = time.Parse(time.RFC3339, stringStartDate)
		if err != nil {
			return nil, err
		}

		endDate, err = time.Parse(time.RFC3339, stringEndDate)
		if err != nil {
			return nil, err
		}
	}

	//m.TotaPenerima = u.service.CountByDate(criteria, date)
	m.TotaPenerima = u.service.CountByRangeDate(criteria, startDate, endDate)

	status = "PARTIAL_DONE"
	criteria["status"] = status
	//m.TotalPartialDone = u.service.CountByDate(criteria, date)
	m.TotalPartialDone = u.service.CountByRangeDate(criteria, startDate, endDate)

	status = "REJECTED"
	criteria["status"] = status
	//m.TotalDataGugur = u.service.CountByDate(criteria, date)
	m.TotalDataGugur = u.service.CountByRangeDate(criteria, startDate, endDate)

	status = "DONE"
	criteria["status"] = status
	//m.TotalSudahMenerima = u.service.CountByDate(criteria, date)
	m.TotalSudahMenerima = u.service.CountByRangeDate(criteria, startDate, endDate)

	status = "NOT DONE"
	criteria["status"] = status

	m.TotalQuota = quota

	if m.TotalQuota > 0 {
		m.TotalQuota = m.TotalQuota - m.TotalSudahMenerima
	}

	//m.TotalBelumMenerima = u.service.CountByDate(criteria, date)
	m.TotalBelumMenerima = u.service.CountByRangeDate(criteria, startDate, endDate)

	return &m, nil
}

func (u *usecase) Export(input request.Report) (string, error) {
	r := helper.NewRequestPdf("")
	//var date time.Time
	var startDate, endDate time.Time
	var err error
	criteria := make(map[string]interface{})
	criteria["status"] = "DONE"
	if input.Provinsi != "" {
		criteria["provinsi"] = input.Provinsi
	}

	if input.Date != "" {
		//stringDate := req.Date + "T00:00:00.00Z"
		//date, err = time.Parse(time.RFC3339, stringDate)
		//if err != nil {
		//	return nil, err
		//}

		stringStartDate := input.Date + "T00:00:00.00Z"
		stringEndDate := input.Date + "T23:59:59.999:Z"
		startDate, err = time.Parse(time.RFC3339, stringStartDate)
		if err != nil {
			return "", err
		}

		endDate, err = time.Parse(time.RFC3339, stringEndDate)
		if err != nil {
			return "", err
		}
	}

	//reports, err := u.service.ReadAllReport(criteria, date)
	reports, err := u.service.ReadAllReportByRangeDate(criteria, startDate, endDate)
	if err != nil {
		return "", err
	}

	for _, value := range reports {
		if value.Image != "" {

			arr := strings.SplitAfter(value.Image, "/")
			//value.Image = "uploads/" + arr[1]

			// Read the entire file into a byte slice
			bytes, errorReadFile := os.ReadFile("./uploads/" + arr[1])
			if errorReadFile != nil {
				return "", errorReadFile
			}

			var base64Encoding string

			// Determine the content type of the image file
			mimeType := http.DetectContentType(bytes)

			// Prepend the appropriate URI scheme header depending
			// on the MIME type
			switch mimeType {
			case "image/jpeg":
				base64Encoding += "data:image/jpeg;base64,"
			case "image/png":
				base64Encoding += "data:image/png;base64,"
			}

			// Append the base64 encoded output
			base64Encoding += base64.StdEncoding.EncodeToString(bytes)

			// Print the full base64 representation of the image
			value.ImageB64 = template.URL(base64Encoding)
		}

	}

	templateData := struct {
		Provinsi string
		Kota     string
		Date     string
		Jam      template.HTML
		Evaluasi template.HTML
		Solusi   template.HTML
		Reports  []*model.Report
	}{
		Provinsi: input.Provinsi,
		Kota:     input.Kota,
		Date:     input.Date,
		Jam:      input.Jam,
		Evaluasi: input.Evaluasi,
		Solusi:   input.Solusi,
		Reports:  reports,
	}

	//html template path
	templatePath := "templates/report.html"

	currentTime := time.Now()
	filename := currentTime.Format("20060102150405") + "-report.pdf"

	//path for download pdf
	outputPath := "uploads/" + filename

	if err := r.ParseTemplate(templatePath, templateData); err == nil {

		// Generate PDF with custom arguments
		args := []string{"no-pdf-compression"}

		// Generate PDF
		ok, errorGenerate := r.GeneratePDF(outputPath, args)
		if errorGenerate != nil {
			return "", errorGenerate
		}
		fmt.Println(ok, "pdf generated successfully")
	} else {
		fmt.Printf("error: %v", err)
	}
	return "image/" + filename, nil
}
