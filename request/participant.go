/*
 * Created on 15/09/23 14.32
 *
 * Copyright (c) 2023 Abdul Ghani Abbasi
 */

package request

type UpdateParticipant struct {
}

type ParticipantPaged struct {
	Search   string `form:"search"`
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	Provinsi string `form:"provinsi"`
	Kota     string `form:"kota"`
}
