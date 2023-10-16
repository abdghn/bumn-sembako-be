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
	"fmt"
	"html/template"
	"net/http"
	"os"
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
	GetTotalDashboard(req request.ParticipantFilter) (*model.TotalParticipantResponse, error)
	BulkCreate(req request.ImportParticipant) (*model.ImportLog, error)
	Export(input request.Report) (string, error)
}

type usecase struct {
	service participant.Service
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

	return u.service.Count(criteria)
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
		stringEndDate := input.Date + "T23:59:59.999Z"
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

func (u *usecase) BulkCreate(req request.ImportParticipant) (*model.ImportLog, error) {
	//var m model.ImportLog
	var successRows, totalRows, failedRows int
	var status string
	newFile := excelize.NewFile()

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
	for i := 2; i < 20000; i++ {
		var note []string
		row := &request.ParticipantInput{
			Name:               xlsx.GetCellValue(sheet1Name, fmt.Sprintf("A%d", i)),
			NIK:                xlsx.GetCellValue(sheet1Name, fmt.Sprintf("B%d", i)),
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
			break;
		}

		if row.Name == "" {
			note = append(note, "Nama Kosong \n")
		}

		if row.NIK == "" {
			note = append(note, "NIK Kosong \n")
		}

		if len(row.NIK) > 16 {
			note = append(note, "NIK lebih dari 16 karakter \n")
		}

		if row.Gender == "" {
			note = append(note, "Jenis Kelamin Kosong \n")
		}

		if row.Phone == "" {
			note = append(note, "No Handphone Kosong \n")
		}

		if len(row.Phone) > 15 {
			note = append(note, "No Handphone dari 15 karakter \n")
		}

		if row.Address == "" {
			note = append(note, "Alamat Kosong \n")
		}

		if row.RT == "" {
			note = append(note, "RT Kosong \n")
		}

		if row.RW == "" {
			note = append(note, "RW Kosong \n")
		}

		if row.Provinsi == "" {
			note = append(note, "Provinsi Kosong \n")
		}

		if row.Kota == "" {
			note = append(note, "Kota Kosong \n")
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

		if row.ResidenceRW == "" {
			note = append(note, "RW Domisili Kosong \n")
		}

		if row.ResidenceProvinsi == "" {
			note = append(note, "Provinsi Domisili Kosong \n")
		}

		if row.ResidenceKota == "" {
			note = append(note, "Kota/Kabupaten Domisili Kosong \n")
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
				Name:               row.Name,
				NIK:                row.NIK,
				Gender:             row.Gender,
				Phone:              row.Phone,
				Address:            row.Address,
				RT:                 row.RT,
				RW:                 row.RW,
				Provinsi:           row.Provinsi,
				Kota:               row.Kota,
				Kecamatan:          row.Kecamatan,
				Kelurahan:          row.Kelurahan,
				KodePOS:            row.KodePOS,
				ResidenceAddress:   row.ResidenceAddress,
				ResidenceRT:        row.ResidenceRT,
				ResidenceRW:        row.ResidenceRW,
				ResidenceProvinsi:  row.ResidenceProvinsi,
				ResidenceKota:      row.ResidenceKota,
				ResidenceKecamatan: row.ResidenceKecamatan,
				ResidenceKelurahan: row.ResidenceKelurahan,
				ResidenceKodePOS:   row.ResidenceKodePOS,
				Status:             row.Status,
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
	}

	newImportLog, err := u.service.CreateLog(m)
	if err != nil {
		return nil, err
	}

	return newImportLog, nil
}
