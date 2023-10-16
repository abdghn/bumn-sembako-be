/*
 * Created on 15/09/23 16.41
 *
 * Copyright (c) 2023 Abdul Ghani Abbasi
 */

package region

import (
	"bumn-sembako-be/model"
	"bumn-sembako-be/request"
	"bumn-sembako-be/service/region"
)

type Usecase interface {
	ReadAllProvinceBy(req request.RegionPaged) ([]*model.Province, error)
	ReadAllRegencyBy(req request.RegionPaged) ([]*model.Regency, error)
	ReadAllDistrictBy(req request.RegionPaged) ([]*model.District, error)
	ReadAllVillageBy(req request.RegionPaged) ([]*model.Village, error)
}

type usecase struct {
	service region.Service
}

func NewUsecase(service region.Service) Usecase {
	return &usecase{service: service}

}

func (u *usecase) ReadAllProvinceBy(req request.RegionPaged) ([]*model.Province, error) {
	criteria := make(map[string]interface{})
	return u.service.ReadAllProvinceBy(criteria, req.Search)

}

func (u *usecase) ReadAllRegencyBy(req request.RegionPaged) ([]*model.Regency, error) {
	criteria := make(map[string]interface{})
	criteria["province_id"] = req.ProvinceID

	return u.service.ReadAllRegencyBy(criteria, req.Search)

}

func (u *usecase) ReadAllDistrictBy(req request.RegionPaged) ([]*model.District, error) {
	criteria := make(map[string]interface{})
	criteria["regency_id"] = req.RegencyID

	return u.service.ReadAllDistrictBy(criteria, req.Search)
}

func (u *usecase) ReadAllVillageBy(req request.RegionPaged) ([]*model.Village, error) {
	criteria := make(map[string]interface{})
	criteria["district_id"] = req.DistrictID

	return u.service.ReadAllVillageBy(criteria, req.Search)

}
