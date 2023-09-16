/*
 * Created on 15/09/23 23.10
 *
 * Copyright (c) 2023 Abdul Ghani Abbasi
 */

package request

type RegionPaged struct {
	ProvinceID int    `form:"province_id"`
	RegencyID  int    `form:"regency_id"`
	DistrictID int    `form:"district_id"`
	Search     string `form:"search"`
	Page       int    `form:"page"`
	Size       int    `form:"size"`
}
