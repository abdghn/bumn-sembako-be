/*
 * Created on 15/09/23 02.31
 *
 * Copyright (c) 2023 Abdul Ghani Abbasi
 */

package participant

import (
	"bumn-sembako-be/model"
	"bumn-sembako-be/request"
	"bumn-sembako-be/service/participant"
)

type Usecase interface {
	ReadAllBy(req request.ParticipantPaged) (*[]model.Participant, error)
	Count(req request.ParticipantPaged) int64
	ReadById(id int) (*model.Participant, error)
	Update(id int, input request.UpdateParticipant) (*model.Participant, error)
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

	return u.service.Count(criteria)
}

func (u *usecase) Update(id int, input request.UpdateParticipant) (*model.Participant, error) {
	var err error
	participant, err := u.ReadById(id)
	if err != nil {
		return nil, err
	}

	if input.Status == "PARTIAL_DONE" {
		req := &request.PartialDone{Status: "PARTIAL_DONE"}

		updatedParticipant, err := u.service.UpdateStatus(participant.ID, req)
		if err != nil {
			return nil, err
		}

		return updatedParticipant, nil

	} else if input.Status == "REJECTED" {
		req := &request.PartialDone{Status: "REJECTED"}

		_, err = u.service.UpdateStatus(id, req)
		if err != nil {
			return nil, err
		}

		m := &model.Participant{
			Name:      input.Name,
			NIK:       input.NIK,
			Gender:    input.Gender,
			Phone:     input.Phone,
			Address:   input.Address,
			RT:        input.RT,
			RW:        input.RW,
			Provinsi:  input.Provinsi,
			Kota:      input.Kota,
			Kecamatan: input.Kecamatan,
			Kelurahan: input.Kelurahan,
			KodePOS:   input.KodePOS,
			Image:     input.Image,
			Status:    "DONE",
		}

		newParticipant, err := u.service.Create(m)
		if err != nil {
			return nil, err
		}

		return newParticipant, nil

	} else if input.Status == "DONE" {
		req := &request.PartialDone{Status: "DONE", Image: input.Image}

		updateParticipant, err := u.service.UpdateStatus(id, req)
		if err != nil {
			return nil, err
		}

		return updateParticipant, nil

	}

	return participant, nil

}
