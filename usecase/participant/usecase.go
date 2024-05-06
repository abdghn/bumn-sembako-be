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
	"bumn-sembako-be/service/region"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"gorm.io/gorm"
)

type Usecase interface {
	ReadAllBy(req request.ParticipantPaged) (*[]model.Participant, error)
	ReadAllLogBy(req request.ParticipantPaged) ([]model.ImportLog, error)
	Count(req request.ParticipantPaged) int64
	CountLogs(req request.ParticipantPaged) int64
	ReadById(id int) (*model.Participant, error)
	Update(id int, input request.UpdateParticipant) (*model.Participant, error)
	Edit(id int, input request.UpdateParticipant) (*model.Participant, error)
	GetTotalDashboard(req request.ParticipantFilter) (*model.TotalParticipantResponse, error)
	GetTotalDashboardV2(req request.ParticipantFilter) (*model.TotalParticipantResponse, error)
	BulkCreate(req request.ImportParticipant) (*model.ImportLog, error)
	ExportExcel(req request.ParticipantFilter) (string, error)
	ExportExcelData(req request.ParticipantFilter) (string, error)
	ExportCSVData(req request.ParticipantFilter) (string, error)
	Export(input request.Report) ([]*model.ReportPerFile, error)
	ConvertBase64(path string) template.URL
	ExportV2(input request.Report) ([]*model.ReportPerFile, error)
	Reset(id int) (*model.Participant, error)
	Delete(id int) error
	DeleteBy() error
}

type usecase struct {
	service       participant.Service
	regionService region.Service
}

func NewUsecase(service participant.Service, regionService region.Service) Usecase {
	return &usecase{service: service, regionService: regionService}

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

	if req.Type != "" {
		criteria["type"] = req.Type
	}

	return u.service.ReadAllBy(criteria, req.Search, req.Page, req.Size)
}

func (u *usecase) ReadAllLogBy(req request.ParticipantPaged) ([]model.ImportLog, error) {
	criteria := make(map[string]interface{})
	return u.service.ReadAllLogBy(criteria, req.Search, req.Page, req.Size)
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

	if req.Type != "" {
		criteria["type"] = req.Type
	}

	return u.service.Count(criteria, req.Search)
}

func (u *usecase) CountLogs(req request.ParticipantPaged) int64 {
	criteria := make(map[string]interface{})
	return u.service.CountLogs(criteria)
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

		participantRejected, err := u.service.UpdateStatus(id, req)
		if err != nil {
			return nil, err
		}

		err = u.service.Delete(participantRejected.ID)
		if err != nil {
			return nil, err
		}

		nik := u.service.Count(map[string]interface{}{"nik": input.NIK}, "")
		if nik > 0 {
			return nil, fmt.Errorf("NIK already exists")
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
			Type:               participantRejected.Type,
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

func (u *usecase) Edit(id int, input request.UpdateParticipant) (*model.Participant, error) {
	var err error
	model, err := u.ReadById(id)
	if err != nil {
		return nil, err
	}

	nik := u.service.CountNotInId(map[string]interface{}{"nik": input.NIK}, id)
	if nik > 0 {
		return nil, fmt.Errorf("NIK already exists")
	}

	req := &request.ParticipantEditInput{
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
		Status:             input.Status,
		UpdatedBy:          model.UpdatedBy,
		Type:               input.Type,
	}

	if input.Image != "" {
		req.Image = input.Image
	} else {
		req.Image = model.Image
	}

	if input.ImagePenerima != "" {
		req.ImagePenerima = input.ImagePenerima
	} else {
		req.ImagePenerima = model.ImagePenerima
	}

	updatedParticipant, err := u.service.Update(id, req)
	if err != nil {
		return nil, err
	}

	return updatedParticipant, err

}

func (u *usecase) GetTotalDashboard(req request.ParticipantFilter) (*model.TotalParticipantResponse, error) {
	var m model.TotalParticipantResponse
	//var date time.Time
	var startDate, endDate time.Time
	status := ""
	quota := int64(0)

	criteria := make(map[string]interface{})
	if req.Provinsi != "" {
		criteria["residence_provinsi"] = req.Provinsi
	}

	if req.Kota != "" {
		criteria["residence_kota"] = req.Kota
	}

	if req.Type != "" {
		criteria["type"] = req.Type
	}

	dataQuota, err := u.service.GetQuota(criteria)
	if err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {
		return nil, err
	}

	if dataQuota.ID > 0 {
		quota = dataQuota.Total
	}

	if req.Date != "" {
		//stringDate := req.Date + "T00:00:00.00Z"
		//date, err = time.Parse(time.RFC3339, stringDate)
		//if err != nil {
		//	return nil, err
		//}

		stringStartDate := req.Date + "T00:00:00.00Z"
		stringEndDate := req.Date + "T23:59:59.999Z"
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
	m.TotalPenerima = u.service.CountByRangeDate(criteria, startDate, endDate)

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

func (u *usecase) GetTotalDashboardV2(req request.ParticipantFilter) (*model.TotalParticipantResponse, error) {
	criteria := make(map[string]interface{})
	if req.Provinsi != "" {
		criteria["residence_provinsi"] = req.Provinsi
	}

	if req.Kota != "" {
		criteria["residence_kota"] = req.Kota
	}

	if req.Kecamatan != "" {
		criteria["residence_kecamatan"] = req.Kecamatan
	}

	if req.Kelurahan != "" {
		criteria["residence_kelurahan"] = req.Kelurahan
	}

	if req.Type != "" {
		criteria["type"] = req.Type
	}

	return u.service.CountAllStatus(criteria)
}

func (u *usecase) Export(input request.Report) ([]*model.ReportPerFile, error) {
	r := helper.NewRequestPdf("")

	var reportPerFile []*model.ReportPerFile
	//var date time.Time
	var startDate, endDate time.Time
	var totalData int
	const limit = 500
	criteria := make(map[string]interface{})
	criteria["status"] = "DONE"
	if input.Provinsi != "" {
		criteria["residence_provinsi"] = input.Provinsi
	}

	if input.Kota != "" {
		criteria["residence_kota"] = input.Kota
	}

	if input.HasPrinted != "" {
		if input.HasPrinted == "PRINTED" {
			criteria["has_printed"] = true
		} else if input.HasPrinted == "NOT PRINTED" {
			criteria["has_printed"] = false
		}
	}

	if input.Type != "" {
		criteria["type"] = input.Type
	}

	totalPage := int(math.Ceil(float64(input.TotalSudahMenerima) / float64(limit)))

	for i := 1; i <= totalPage; i++ {
		reports, err := u.service.ReadAllReportByRangeDateV2(criteria, startDate, endDate, i, limit)
		if err != nil {
			return nil, err
		}

		if len(reports) > 0 {
			totalData = reports[0].Total

		} else {
			return reportPerFile, nil
		}

		//for _, value := range reports {
		//	if value.Image != "" {
		//
		//		arr := strings.SplitAfter(value.Image, "/")
		//		//value.Image = "uploads/" + arr[1]
		//
		//		// Read the entire file into a byte slice
		//		bytes, errorReadFile := os.ReadFile("./uploads/" + arr[1])
		//		if errorReadFile != nil {
		//			continue
		//		}
		//
		//		var base64Encoding string
		//
		//		// Determine the content type of the image file
		//		mimeType := http.DetectContentType(bytes)
		//
		//		// Prepend the appropriate URI scheme header depending
		//		// on the MIME type
		//		switch mimeType {
		//		case "image/jpeg":
		//			base64Encoding += "data:image/jpeg;base64,"
		//		case "image/png":
		//			base64Encoding += "data:image/png;base64,"
		//		}
		//
		//		// Append the base64 encoded output
		//		base64Encoding += base64.StdEncoding.EncodeToString(bytes)
		//
		//		// Print the full base64 representation of the image
		//		value.ImageB64 = template.URL(base64Encoding)
		//	}
		//
		//}

		templateData := struct {
			Provinsi string
			Kota     string
			Date     string
			Jam      template.HTML
			Evaluasi template.HTML
			Solusi   template.HTML
			Reports  []*model.Report
			Total    int
		}{
			Provinsi: input.Provinsi,
			Kota:     input.Kota,
			Date:     input.Date,
			Jam:      input.Jam,
			Evaluasi: input.Evaluasi,
			Solusi:   input.Solusi,
			Reports:  reports,
			Total:    totalData,
		}

		//html template path
		templatePath := "templates/report.html"

		currentTime := time.Now()
		filename := currentTime.Format("20060102150405") + "-report.pdf"

		//path for download pdf
		outputPath := "uploads/" + filename

		if err := r.ParseTemplate(templatePath, templateData); err == nil {

			// Generate PDF with custom arguments
			args := []string{"low-quality"}

			// Generate PDF
			ok, errorGenerate := r.GeneratePDF(outputPath, args)
			if errorGenerate != nil {
				helper.CommonLogger().Error(errorGenerate)
				return nil, errorGenerate
			}
			fmt.Println(ok, "pdf generated successfully")
		} else {
			helper.CommonLogger().Error(err)
			fmt.Printf("error: %v", err)
		}

		r := &model.ReportPerFile{
			Name: filename,
			Path: "image/" + filename,
		}

		reportPerFile = append(reportPerFile, r)

		go func(list []*model.Report) {

			for _, report := range list {
				requestInput := &request.ParticipantEditInput{
					HasPrinted: true,
				}

				_, err = u.service.Update(report.ID, requestInput)
				if err != nil {
					return
				}

			}
		}(reports)

	}

	// if input.Date != "" {
	// 	//stringDate := req.Date + "T00:00:00.00Z"
	// 	//date, err = time.Parse(time.RFC3339, stringDate)
	// 	//if err != nil {
	// 	//	return nil, err
	// 	//}

	// 	stringStartDate := input.Date + "T00:00:00.00Z"
	// 	stringEndDate := input.Date + "T23:59:59.999Z"
	// 	startDate, err = time.Parse(time.RFC3339, stringStartDate)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	endDate, err = time.Parse(time.RFC3339, stringEndDate)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	//reports, err := u.service.ReadAllReport(criteria, date)

	return reportPerFile, nil
}

func (u *usecase) BulkCreate(req request.ImportParticipant) (*model.ImportLog, error) {
	//var m model.ImportLog
	var successRows, totalRows, failedRows int
	var status string
	newFile := excelize.NewFile()
	randString, err := helper.Randstring(20)
	if err != nil {
		return nil, err
	}

	xlsx, err := excelize.OpenFile(req.TmpPath)
	if err != nil {
		return nil, fmt.Errorf("error when open file: %v", err)
	}

	sheet1Name := "Sheet1"

	newFile.SetSheetName(newFile.GetSheetName(1), sheet1Name)
	newFile.SetCellValue(sheet1Name, "A1", "Nama")
	newFile.SetCellValue(sheet1Name, "B1", "NIK")
	newFile.SetCellValue(sheet1Name, "C1", "Jenis Kelamin")
	newFile.SetCellValue(sheet1Name, "D1", "No Handphone")
	newFile.SetCellValue(sheet1Name, "E1", "Alamat Sesuai KTP")
	newFile.SetCellValue(sheet1Name, "F1", "RT")
	newFile.SetCellValue(sheet1Name, "G1", "RW")
	newFile.SetCellValue(sheet1Name, "H1", "Provinsi")
	newFile.SetCellValue(sheet1Name, "I1", "Kota/Kabupaten")
	newFile.SetCellValue(sheet1Name, "J1", "Kecamatan")
	newFile.SetCellValue(sheet1Name, "K1", "Kelurahan")
	newFile.SetCellValue(sheet1Name, "L1", "Kode Pos")
	newFile.SetCellValue(sheet1Name, "M1", "Alamat Domisili")
	newFile.SetCellValue(sheet1Name, "N1", "RT Domisili")
	newFile.SetCellValue(sheet1Name, "O1", "RW Domisili")
	newFile.SetCellValue(sheet1Name, "P1", "Provinsi Domisili")
	newFile.SetCellValue(sheet1Name, "Q1", "Kota/Kabupaten Domisili")
	newFile.SetCellValue(sheet1Name, "R1", "Kecamatan Domisili")
	newFile.SetCellValue(sheet1Name, "S1", "Kelurahan Domisili")
	newFile.SetCellValue(sheet1Name, "T1", "Kode Pos Domisili")
	newFile.SetCellValue(sheet1Name, "U1", "Status")
	newFile.SetCellValue(sheet1Name, "V1", "Catatan")

	var rows []*request.ParticipantInput
	for i := 2; i < 60000; i++ {
		var note []string
		row := &request.ParticipantInput{
			Name:               strings.TrimSpace(xlsx.GetCellValue(sheet1Name, fmt.Sprintf("A%d", i))),
			NIK:                strings.TrimSpace(xlsx.GetCellValue(sheet1Name, fmt.Sprintf("B%d", i))),
			Gender:             xlsx.GetCellValue(sheet1Name, fmt.Sprintf("C%d", i)),
			Phone:              xlsx.GetCellValue(sheet1Name, fmt.Sprintf("D%d", i)),
			Address:            xlsx.GetCellValue(sheet1Name, fmt.Sprintf("E%d", i)),
			RT:                 xlsx.GetCellValue(sheet1Name, fmt.Sprintf("F%d", i)),
			RW:                 xlsx.GetCellValue(sheet1Name, fmt.Sprintf("G%d", i)),
			Provinsi:           xlsx.GetCellValue(sheet1Name, fmt.Sprintf("H%d", i)),
			Kota:               xlsx.GetCellValue(sheet1Name, fmt.Sprintf("I%d", i)),
			Kecamatan:          xlsx.GetCellValue(sheet1Name, fmt.Sprintf("J%d", i)),
			Kelurahan:          xlsx.GetCellValue(sheet1Name, fmt.Sprintf("K%d", i)),
			KodePOS:            xlsx.GetCellValue(sheet1Name, fmt.Sprintf("L%d", i)),
			ResidenceAddress:   xlsx.GetCellValue(sheet1Name, fmt.Sprintf("M%d", i)),
			ResidenceRT:        xlsx.GetCellValue(sheet1Name, fmt.Sprintf("N%d", i)),
			ResidenceRW:        xlsx.GetCellValue(sheet1Name, fmt.Sprintf("O%d", i)),
			ResidenceProvinsi:  xlsx.GetCellValue(sheet1Name, fmt.Sprintf("P%d", i)),
			ResidenceKota:      xlsx.GetCellValue(sheet1Name, fmt.Sprintf("Q%d", i)),
			ResidenceKecamatan: xlsx.GetCellValue(sheet1Name, fmt.Sprintf("R%d", i)),
			ResidenceKelurahan: xlsx.GetCellValue(sheet1Name, fmt.Sprintf("S%d", i)),
			ResidenceKodePOS:   xlsx.GetCellValue(sheet1Name, fmt.Sprintf("T%d", i)),
			Status:             xlsx.GetCellValue(sheet1Name, fmt.Sprintf("U%d", i)),
		}

		if row.Name == "" && row.NIK == "" && row.Gender == "" {
			break
		}

		if row.Name == "" {
			note = append(note, "Nama Kosong \n")
		} else {
			trimString := strings.ReplaceAll(row.Name, " ", "")
			validName := helper.ContainString("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ'.,-", trimString)
			if validName {
				note = append(note, "Nama Tidak Sesuai Format \n")
			}

		}

		if row.NIK == "" {
			note = append(note, "NIK Kosong \n")
		}

		if !regexp.MustCompile(`\d`).MatchString(row.NIK) {
			note = append(note, "NIK terdapat karakter atau simbol karakter \n")
		} else {
			if len(strings.TrimLeft(row.NIK, " ")) != 16 {
				note = append(note, "NIK tidak 16 digit \n")
			} else {
				validChar := helper.ContainString("1234567890", row.NIK)
				if validChar {
					note = append(note, "NIK Tidak Sesuai Format \n")
				}

				countBynik := u.service.CheckId(map[string]interface{}{"nik": row.NIK})
				if countBynik > 0 {
					note = append(note, "NIK Sudah Terdaftar \n")
				}
			}
		}

		if row.Gender == "" {
			note = append(note, "Jenis Kelamin Kosong \n")
		}

		if row.Phone == "" {
			note = append(note, "No Handphone Kosong \n")
		} else {
			if len(row.Phone) > 13 {
				note = append(note, "No Handphone lebih dari 13 digit \n")
			}

			if len(row.Phone) < 10 {
				note = append(note, "No Handphone kurang dari 10 digit \n")
			}

			if len(row.Phone) > 9 && len(row.Phone) < 14 {
				validPhone := helper.ContainString("1234567890+", row.Phone)
				if validPhone {
					note = append(note, "No Handphone Tidak Sesuai Format \n")
				}

			}

		}

		if row.Address == "" {
			note = append(note, "Alamat Kosong \n")
		}

		if row.RT == "" {
			note = append(note, "RT Kosong \n")
		}

		if len(row.RT) > 3 {
			note = append(note, "RT lebih dari 3 digit \n")
		}

		if row.RW == "" {
			note = append(note, "RW Kosong \n")
		}

		if len(row.RW) > 3 {
			note = append(note, "RW lebih dari 3 digit \n")
		}

		if row.Provinsi == "" {
			note = append(note, "Provinsi Kosong \n")
		} else {
			criteria := make(map[string]interface{})

			criteria["name"] = strings.ToUpper(row.Provinsi)

			_, err = u.regionService.ReadProvinceBy(criteria)
			if err != nil {
				note = append(note, "Provinsi tidak terdaftar \n")
			}

		}

		if row.Kota == "" {
			note = append(note, "Kota/Kabupaten Kosong \n")
		} else {

			criteria := make(map[string]interface{})

			criteria["name"] = strings.ToUpper(row.Kota)

			_, err = u.regionService.ReadRegencyBy(criteria)
			if err != nil {
				note = append(note, "Kota/Kabupaten tidak terdaftar \n")
			}

		}

		if row.Kecamatan == "" {
			note = append(note, "Kecamatan Kosong \n")
		}

		if row.Kelurahan == "" {
			note = append(note, "Kelurahan Kosong \n")
		}

		if row.KodePOS == "" {
			note = append(note, "Kode POS Kosong \n")
		}

		if row.ResidenceAddress == "" {
			note = append(note, "Alamat Domisili Kosong \n")
		}

		if row.ResidenceRT == "" {
			note = append(note, "RT Domisili Kosong \n")
		}

		if len(row.ResidenceRT) > 3 {
			note = append(note, "RT Domisili lebih dari 3 digit \n")
		}

		if row.ResidenceRW == "" {
			note = append(note, "RW Domisili Kosong \n")
		}

		if len(row.RW) > 3 {
			note = append(note, "RW Domisili lebih dari 3 digit \n")
		}

		if row.ResidenceProvinsi == "" {
			note = append(note, "Provinsi Domisili Kosong \n")
		} else {
			criteria := make(map[string]interface{})

			criteria["name"] = strings.ToUpper(row.ResidenceProvinsi)

			_, err = u.regionService.ReadProvinceBy(criteria)
			if err != nil {
				note = append(note, "Provinsi Domisili tidak terdaftar \n")
			}

		}

		if row.ResidenceKota == "" {
			note = append(note, "Kota/Kabupaten Domisili Kosong \n")
		} else {

			criteria := make(map[string]interface{})

			criteria["name"] = strings.ToUpper(row.ResidenceKota)

			_, err = u.regionService.ReadRegencyBy(criteria)
			if err != nil {
				note = append(note, "Kota/Kabupaten Domisili tidak terdaftar \n")
			}

		}

		if row.ResidenceKecamatan == "" {
			note = append(note, "Domisili Kecamatan Kosong \n")
		}

		if row.ResidenceKelurahan == "" {
			note = append(note, "Domisili Kelurahan \n")
		}

		if row.ResidenceKodePOS == "" {
			note = append(note, "Domisili Kode POS Kosong \n")
		}

		if row.Status == "" {
			note = append(note, "Status Kosong \n")
		}

		if len(note) == 0 {
			newParticipant := &model.Participant{
				Name:               strings.TrimSpace(row.Name),
				NIK:                strings.TrimSpace(row.NIK),
				Gender:             row.Gender,
				Phone:              row.Phone,
				Address:            row.Address,
				RT:                 row.RT,
				RW:                 row.RW,
				Provinsi:           strings.ToUpper(row.Provinsi),
				Kota:               strings.ToUpper(row.Kota),
				Kecamatan:          strings.ToUpper(row.Kecamatan),
				Kelurahan:          strings.ToUpper(row.Kelurahan),
				KodePOS:            row.KodePOS,
				ResidenceAddress:   row.ResidenceAddress,
				ResidenceRT:        row.ResidenceRT,
				ResidenceRW:        row.ResidenceRW,
				ResidenceProvinsi:  strings.ToUpper(row.ResidenceProvinsi),
				ResidenceKota:      strings.ToUpper(row.ResidenceKota),
				ResidenceKecamatan: strings.ToUpper(row.ResidenceKecamatan),
				ResidenceKelurahan: strings.ToUpper(row.ResidenceKelurahan),
				ResidenceKodePOS:   row.ResidenceKodePOS,
				Status:             row.Status,
				Reference:          randString,
				Type:               req.Type,
			}

			_, err = u.service.Create(newParticipant)
			if err != nil {
				return nil, err
			}
			successRows++
		} else {
			index := 2
			row.Note = strings.Join(note, ",")
			failedRows++
			rows = append(rows, row)
			index++

		}

		totalRows++
	}
	if failedRows > 0 {

		for i, row := range rows {
			newFile.SetCellValue(sheet1Name, fmt.Sprintf("A%d", i+2), row.Name)
			newFile.SetCellValue(sheet1Name, fmt.Sprintf("B%d", i+2), row.NIK)
			newFile.SetCellValue(sheet1Name, fmt.Sprintf("C%d", i+2), row.Gender)
			newFile.SetCellValue(sheet1Name, fmt.Sprintf("D%d", i+2), row.Phone)
			newFile.SetCellValue(sheet1Name, fmt.Sprintf("E%d", i+2), row.Address)
			newFile.SetCellValue(sheet1Name, fmt.Sprintf("F%d", i+2), row.RT)
			newFile.SetCellValue(sheet1Name, fmt.Sprintf("G%d", i+2), row.RW)
			newFile.SetCellValue(sheet1Name, fmt.Sprintf("H%d", i+2), row.Provinsi)
			newFile.SetCellValue(sheet1Name, fmt.Sprintf("I%d", i+2), row.Kota)
			newFile.SetCellValue(sheet1Name, fmt.Sprintf("J%d", i+2), row.Kecamatan)
			newFile.SetCellValue(sheet1Name, fmt.Sprintf("K%d", i+2), row.Kelurahan)
			newFile.SetCellValue(sheet1Name, fmt.Sprintf("L%d", i+2), row.KodePOS)
			newFile.SetCellValue(sheet1Name, fmt.Sprintf("M%d", i+2), row.ResidenceAddress)
			newFile.SetCellValue(sheet1Name, fmt.Sprintf("N%d", i+2), row.ResidenceRT)
			newFile.SetCellValue(sheet1Name, fmt.Sprintf("O%d", i+2), row.ResidenceRW)
			newFile.SetCellValue(sheet1Name, fmt.Sprintf("P%d", i+2), row.ResidenceProvinsi)
			newFile.SetCellValue(sheet1Name, fmt.Sprintf("Q%d", i+2), row.ResidenceKota)
			newFile.SetCellValue(sheet1Name, fmt.Sprintf("R%d", i+2), row.ResidenceKecamatan)
			newFile.SetCellValue(sheet1Name, fmt.Sprintf("S%d", i+2), row.ResidenceKelurahan)
			newFile.SetCellValue(sheet1Name, fmt.Sprintf("T%d", i+2), row.ResidenceKodePOS)
			newFile.SetCellValue(sheet1Name, fmt.Sprintf("U%d", i+2), row.Status)
			newFile.SetCellValue(sheet1Name, fmt.Sprintf("V%d", i+2), row.Note)
		}

		path := "./uploads"
		ext := ".xlsx"
		currentTime := time.Now()
		filename := currentTime.Format("20060102150405") + ext
		tmpFile := path + "/" + filename
		err = newFile.SaveAs(tmpFile)
		req.Path = "image/" + filename
		if err != nil {
			return nil, err
		}
	}

	if failedRows > 0 && successRows == 0 {
		status = "Error All"
	} else if successRows > 0 && failedRows > 0 {
		status = "Success With error"
	} else {
		status = "Success All"
	}

	m := &model.ImportLog{
		FileName:    req.Name,
		Status:      status,
		TotalRows:   totalRows,
		SuccessRows: successRows,
		FailedRows:  failedRows,
		Path:        req.Path,
		UploadedBy:  req.UploadedBy,
		Type:        req.Type,
	}

	if successRows > 0 {
		m.Reference = randString
	}

	newImportLog, err := u.service.CreateLog(m)
	if err != nil {
		return nil, err
	}

	return newImportLog, nil
}

func (u *usecase) ExportExcel(req request.ParticipantFilter) (string, error) {
	criteria := make(map[string]interface{})
	xlsx := excelize.NewFile()
	sheet1Name := "Sheet1"
	xlsx.SetSheetName(xlsx.GetSheetName(1), sheet1Name)
	xlsx.SetCellValue(sheet1Name, "A1", "Provinsi")
	xlsx.SetCellValue(sheet1Name, "B1", "Kota")
	xlsx.SetCellValue(sheet1Name, "C1", "Total Penerima Bantuan Kota Anda")
	xlsx.SetCellValue(sheet1Name, "D1", "Sudah Menerima Bantuan")
	xlsx.SetCellValue(sheet1Name, "E1", "Belum Unggah Foto Penerima Bantuan")
	xlsx.SetCellValue(sheet1Name, "F1", "Belum Menerima Bantuan")
	xlsx.SetCellValue(sheet1Name, "G1", "Data Tidak Sesuai")

	if req.Provinsi != "" {
		criteria["residence_provinsi"] = req.Provinsi
	}

	if req.Kota != "" {
		criteria["residence_kota"] = req.Kota
	}

	if req.Type != "" {
		criteria["type"] = req.Type
	}

	rows, err := u.service.CountAllStatusGroup(criteria)
	if err != nil {
		return "", err
	}

	fmt.Print(len(rows))

	for i, row := range rows {
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("A%d", i+2), row.ResidenceProvinsi)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("B%d", i+2), row.ResidenceKota)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("C%d", i+2), row.TotalPenerima)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("D%d", i+2), row.TotalSudahMenerima)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("E%d", i+2), row.TotalPartialDone)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("F%d", i+2), row.TotalBelumMenerima)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("G%d", i+2), row.TotalDataGugur)
	}

	path := "./uploads"
	ext := ".xlsx"
	currentTime := time.Now()
	filename := currentTime.Format("20060102150405") + "-dashboard" + ext
	tmpFile := path + "/" + filename
	err = xlsx.SaveAs(tmpFile)
	if err != nil {
		return "", nil
	}

	timer1 := time.NewTimer(30 * time.Second)
	go func(filename string) {
		<-timer1.C
		os.Remove(filename)
	}(tmpFile)

	return "image/" + filename, nil

}

func (u *usecase) ExportExcelData(req request.ParticipantFilter) (string, error) {
	criteria := make(map[string]interface{})
	xlsx := excelize.NewFile()
	sheet1Name := "Sheet1"
	xlsx.SetSheetName(xlsx.GetSheetName(1), sheet1Name)
	xlsx.SetCellValue(sheet1Name, "A1", "Nama")
	xlsx.SetCellValue(sheet1Name, "B1", "NIK")
	xlsx.SetCellValue(sheet1Name, "C1", "Jenis Kelamin")
	xlsx.SetCellValue(sheet1Name, "D1", "No Handphone")
	xlsx.SetCellValue(sheet1Name, "E1", "Alamat Sesuai KTP")
	xlsx.SetCellValue(sheet1Name, "F1", "RT")
	xlsx.SetCellValue(sheet1Name, "G1", "RW")
	xlsx.SetCellValue(sheet1Name, "H1", "Provinsi")
	xlsx.SetCellValue(sheet1Name, "I1", "Kota/Kabupaten")
	xlsx.SetCellValue(sheet1Name, "J1", "Kecamatan")
	xlsx.SetCellValue(sheet1Name, "K1", "Kelurahan")
	xlsx.SetCellValue(sheet1Name, "L1", "Kode Pos")
	xlsx.SetCellValue(sheet1Name, "M1", "Alamat Domisili")
	xlsx.SetCellValue(sheet1Name, "N1", "RT Domisili")
	xlsx.SetCellValue(sheet1Name, "O1", "RW Domisili")
	xlsx.SetCellValue(sheet1Name, "P1", "Provinsi Domisili")
	xlsx.SetCellValue(sheet1Name, "Q1", "Kota/Kabupaten Domisili")
	xlsx.SetCellValue(sheet1Name, "R1", "Kecamatan Domisili")
	xlsx.SetCellValue(sheet1Name, "S1", "Kelurahan Domisili")
	xlsx.SetCellValue(sheet1Name, "T1", "Kode Pos Domisili")
	xlsx.SetCellValue(sheet1Name, "U1", "Status")
	xlsx.SetCellValue(sheet1Name, "V1", "Tipe")

	if req.Provinsi != "" {
		criteria["residence_provinsi"] = req.Provinsi
	}

	if req.Kota != "" {
		criteria["residence_kota"] = req.Kota
	}

	if req.Kecamatan != "" {
		criteria["residence_kecamatan"] = req.Kecamatan
	}

	if req.Kelurahan != "" {
		criteria["residence_kelurahan"] = req.Kelurahan
	}

	if req.Type != "" {
		criteria["type"] = req.Type
	}

	rows, err := u.service.ReadAllWithoutPagination(criteria)
	if err != nil {
		return "", err
	}

	for i, row := range rows {
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("A%d", i+2), row.Name)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("B%d", i+2), row.NIK)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("C%d", i+2), row.Gender)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("D%d", i+2), row.Phone)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("E%d", i+2), row.Address)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("F%d", i+2), row.RT)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("G%d", i+2), row.RW)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("H%d", i+2), row.Provinsi)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("I%d", i+2), row.Kota)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("J%d", i+2), row.Kecamatan)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("K%d", i+2), row.Kelurahan)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("L%d", i+2), row.KodePOS)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("M%d", i+2), row.ResidenceAddress)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("N%d", i+2), row.ResidenceRT)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("O%d", i+2), row.ResidenceRW)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("P%d", i+2), row.ResidenceProvinsi)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("Q%d", i+2), row.ResidenceKota)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("R%d", i+2), row.ResidenceKecamatan)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("S%d", i+2), row.ResidenceKelurahan)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("T%d", i+2), row.ResidenceKodePOS)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("U%d", i+2), row.Status)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("V%d", i+2), row.Type)
	}

	path := "./uploads"
	ext := ".xlsx"
	currentTime := time.Now()
	filename := currentTime.Format("20060102150405") + "-export-data" + ext
	tmpFile := path + "/" + filename
	err = xlsx.SaveAs(tmpFile)
	if err != nil {
		return "", nil
	}

	return "image/" + filename, nil

}

func (u *usecase) ExportCSVData(req request.ParticipantFilter) (string, error) {
	criteria := make(map[string]interface{})
	path := "./uploads"
	ext := ".csv"
	currentTime := time.Now()
	filename := currentTime.Format("20060102150405") + "-export-data" + ext
	tmpFile := path + "/" + filename

	// WRITE TO CSV
	file, err := os.Create(tmpFile)
	defer file.Close()
	if err != nil {
		log.Fatalln("failed to open file", err)
	}

	w := csv.NewWriter(file)
	defer w.Flush()

	//xlsx := excelize.NewFile()
	//sheet1Name := "Sheet1"
	//xlsx.SetSheetName(xlsx.GetSheetName(1), sheet1Name)
	//xlsx.SetCellValue(sheet1Name, "A1", "Nama")
	//xlsx.SetCellValue(sheet1Name, "B1", "NIK")
	//xlsx.SetCellValue(sheet1Name, "C1", "Gender")
	//xlsx.SetCellValue(sheet1Name, "D1", "Phone")
	//xlsx.SetCellValue(sheet1Name, "E1", "Address")
	//xlsx.SetCellValue(sheet1Name, "F1", "RT")
	//xlsx.SetCellValue(sheet1Name, "G1", "RW")

	if req.Provinsi != "" {
		criteria["residence_provinsi"] = req.Provinsi
	}

	if req.Kota != "" {
		criteria["residence_kota"] = req.Kota
	}

	if req.Kecamatan != "" {
		criteria["residence_kecamatan"] = req.Kecamatan
	}

	if req.Kelurahan != "" {
		criteria["residence_kelurahan"] = req.Kelurahan
	}

	if req.Type != "" {
		criteria["type"] = req.Type
	}

	if req.Status != "" {
		criteria["status"] = req.Status
	}

	rows, err := u.service.ReadAllWithoutPagination(criteria)
	if err != nil {
		return "", err
	}

	for _, record := range rows {
		row := []string{
			record.Name,
			"'" + record.NIK,
			record.Gender,
			record.Phone,
			record.Address,
			record.RT,
			record.RW,
			record.Provinsi,
			record.Kota,
			record.Kecamatan,
			record.Kelurahan,
			record.KodePOS,
			record.ResidenceAddress,
			record.ResidenceRT,
			record.ResidenceRW,
			record.ResidenceProvinsi,
			record.ResidenceKota,
			record.ResidenceKecamatan,
			record.ResidenceKelurahan,
			record.ResidenceKodePOS,
			record.Status,
			record.Type,
		}
		if err := w.Write(row); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}
	var data [][]string
	w.WriteAll(data)

	timer1 := time.NewTimer(30 * time.Second)
	go func(filename string) {
		<-timer1.C
		os.Remove(filename)
	}(tmpFile)

	return "image/" + filename, nil

}

func (u *usecase) ExportV2(input request.Report) ([]*model.ReportPerFile, error) {
	r := helper.NewRequestPdf("")

	var reportPerFile []*model.ReportPerFile
	//var date time.Time
	var startDate, endDate time.Time
	var totalData int
	const limit = 100
	criteria := make(map[string]interface{})
	criteria["status"] = "DONE"
	if input.Provinsi != "" {
		criteria["residence_provinsi"] = input.Provinsi
	}

	if input.Kota != "" {
		criteria["residence_kota"] = input.Kota
	}

	if input.Type != "" {
		criteria["type"] = input.Type
	}

	// totalPage := int(math.Ceil(float64(input.TotalSudahMenerima) / float64(limit)))

	// for i := 1; i <= totalPage; i++ {
	// 	reports, err := u.service.ReadAllReportByRangeDate(criteria, startDate, endDate, i, limit)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	if len(reports) > 0 {
	// 		totalData = reports[0].Total

	// 	}

	// 	templateData := struct {
	// 		Provinsi string
	// 		Kota     string
	// 		Date     string
	// 		Jam      template.HTML
	// 		Evaluasi template.HTML
	// 		Solusi   template.HTML
	// 		Reports  []*model.Report
	// 		Total    int
	// 	}{
	// 		Provinsi: input.Provinsi,
	// 		Kota:     input.Kota,
	// 		Date:     input.Date,
	// 		Jam:      input.Jam,
	// 		Evaluasi: input.Evaluasi,
	// 		Solusi:   input.Solusi,
	// 		Reports:  reports,
	// 		Total:    totalData,
	// 	}

	// 	//html template path
	// 	templatePath := "templates/report-v2.html"

	// 	currentTime := time.Now()
	// 	filename := currentTime.Format("20060102150405") + "-report.pdf"

	// 	//path for download pdf
	// 	outputPath := "uploads/" + filename

	// 	if err := r.ParseTemplate(templatePath, templateData); err == nil {

	// 		// Generate PDF with custom arguments
	// 		args := []string{"low-quality"}

	// 		// Generate PDF
	// 		ok, errorGenerate := r.GeneratePDF(outputPath, args)
	// 		if errorGenerate != nil {
	// 			helper.CommonLogger().Error(errorGenerate)
	// 			return nil, errorGenerate
	// 		}
	// 		fmt.Println(ok, "pdf generated successfully")
	// 	} else {
	// 		helper.CommonLogger().Error(err)
	// 		fmt.Printf("error: %v", err)
	// 	}

	// 	r := &model.ReportPerFile{
	// 		Name: filename,
	// 		Path: "image/" + filename,
	// 	}

	// 	reportPerFile = append(reportPerFile, r)

	// }

	// if input.Date != "" {
	// 	//stringDate := req.Date + "T00:00:00.00Z"
	// 	//date, err = time.Parse(time.RFC3339, stringDate)
	// 	//if err != nil {
	// 	//	return nil, err
	// 	//}

	// 	stringStartDate := input.Date + "T00:00:00.00Z"
	// 	stringEndDate := input.Date + "T23:59:59.999Z"
	// 	startDate, err = time.Parse(time.RFC3339, stringStartDate)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	endDate, err = time.Parse(time.RFC3339, stringEndDate)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	//reports, err := u.service.ReadAllReport(criteria, date)

	reports, err := u.service.ReadAllReportByRangeDate(criteria, startDate, endDate)
	if err != nil {
		return nil, err
	}

	if len(reports) > 0 {
		totalData = reports[0].Total

	}

	templateData := struct {
		Provinsi string
		Kota     string
		Date     string
		Jam      template.HTML
		Evaluasi template.HTML
		Solusi   template.HTML
		Reports  []*model.Report
		Total    int
	}{
		Provinsi: input.Provinsi,
		Kota:     input.Kota,
		Date:     input.Date,
		Jam:      input.Jam,
		Evaluasi: input.Evaluasi,
		Solusi:   input.Solusi,
		Reports:  reports,
		Total:    totalData,
	}

	//html template path
	templatePath := "templates/report-v2.html"

	currentTime := time.Now()
	filename := currentTime.Format("20060102150405") + "-report.pdf"

	//path for download pdf
	outputPath := "uploads/" + filename

	if err := r.ParseTemplate(templatePath, templateData); err == nil {

		// Generate PDF with custom arguments
		args := []string{"low-quality"}

		// Generate PDF
		ok, errorGenerate := r.GeneratePDF(outputPath, args)
		if errorGenerate != nil {
			helper.CommonLogger().Error(errorGenerate)
			return nil, errorGenerate
		}
		fmt.Println(ok, "pdf generated successfully")
	} else {
		helper.CommonLogger().Error(err)
		fmt.Printf("error: %v", err)
	}

	reportPerFile = append(reportPerFile, &model.ReportPerFile{
		Name: filename,
		Path: "image/" + filename})

	return reportPerFile, nil
}

func (u *usecase) ConvertBase64(path string) template.URL {
	// arr := strings.SplitAfter(path, "/")
	//value.Image = "uploads/" + arr[1]

	// Read the entire file into a byte slice
	bytes, errorReadFile := os.ReadFile("./uploads/" + path)
	if errorReadFile != nil {
		return ""
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
	return template.URL(base64Encoding)
}

func (u *usecase) Reset(id int) (*model.Participant, error) {
	return u.service.Reset(id)
}

func (u *usecase) Delete(id int) error {
	return u.service.Delete(id)
}

func (u *usecase) DeleteBy() error {
	criteria := make(map[string]interface{})
	var err error
	participants, err := u.service.ReadAllDuplicates()
	if err != nil {
		return err
	}

	// TODO: Update logic many delete
	for _, m := range participants {
		total := u.service.Count(map[string]interface{}{
			"nik":    m.NIK,
			"status": "REJECTED",
		}, "")

		if total == 2 {
			err = u.service.DeleteBy(map[string]interface{}{
				"nik":    m.NIK,
				"status": "REJECTED",
			})
			if err != nil {
				return err
			}
		} else {
			totalNotDone := u.service.Count(map[string]interface{}{
				"nik":    m.NIK,
				"status": "NOT DONE",
			}, "")
			if totalNotDone > 1 {
				err = u.service.DeleteBy(map[string]interface{}{
					"id":     m.ID,
					"status": "NOT DONE",
				})
				if err != nil {
					return err
				}
			} else {

			}
		}

	}
	return u.service.DeleteBy(criteria)
}
