/*
 * Created on 22/09/23 01.35
 *
 * Copyright (c) 2023 Abdul Ghani Abbasi
 */

package helper

import (
	"bytes"
	"html/template"
	"os"
	"strconv"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

// pdf requestpdf struct
type RequestPdf struct {
	body string
}

// new request to pdf function
func NewRequestPdf(body string) *RequestPdf {
	return &RequestPdf{
		body: body,
	}
}

// parsing template function
func (r *RequestPdf) ParseTemplate(templateFileName string, data interface{}) error {

	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}

// generate pdf function
func (r *RequestPdf) GeneratePDF(pdfPath string, args []string) (bool, error) {
	t := time.Now().Unix()
	// write whole the body

	if _, err := os.Stat("./cloneTemplate/"); os.IsNotExist(err) {
		errDir := os.Mkdir("./cloneTemplate/", 0777)
		if errDir != nil {
			return false, err
		}
	}
	err1 := os.WriteFile("./cloneTemplate/"+strconv.FormatInt(int64(t), 10)+".html", []byte(r.body), 0644)
	if err1 != nil {
		return false, err1
	}

	f, err := os.Open("./cloneTemplate/" + strconv.FormatInt(int64(t), 10) + ".html")
	if f != nil {
		defer f.Close()
	}
	if err != nil {
		return false, err
	}

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return false, err
	}

	// Use arguments to customize PDF generation process
	for _, arg := range args {
		switch arg {
		case "low-quality":
			pdfg.LowQuality.Set(true)
		case "no-pdf-compression":
			pdfg.NoPdfCompression.Set(true)
		case "grayscale":
			pdfg.Grayscale.Set(true)
			// Add other arguments as needed
		}
	}

	pageReader := wkhtmltopdf.NewPageReader(f)
	pageReader.NoBackground.Set(true)

	pageReader.PageOptions.EnableLocalFileAccess.Set(true)
	pageReader.PageOptions.DisableSmartShrinking.Set(true)
	pageReader.PageOptions.Encoding.Set("UTF-8")

	pdfg.AddPage(pageReader)

	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)

	pdfg.Dpi.Set(120)

	err = pdfg.Create()
	if err != nil {
		return false, err
	}

	err = pdfg.WriteFile(pdfPath)
	if err != nil {
		return false, err
	}

	dir, err := os.Getwd()
	if err != nil {
		return false, err
	}

	defer os.RemoveAll(dir + "./cloneTemplate")

	return true, nil
}
