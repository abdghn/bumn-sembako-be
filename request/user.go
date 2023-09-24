/*
 * Created on 15/09/23 14.20
 *
 * Copyright (c) 2023 Abdul Ghani Abbasi
 */

package request

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Register struct {
	Name           string `json:"name"`
	Provinsi       string `json:"provinsi"`
	Kota           string `json:"kota"`
	OrganizationID uint    `json:"organization_id"`
	Username       string `json:"username"`
	Password       string `json:"password"`
}

type UpdateUser struct {
	Name        string `json:"name"`
	OldPassword string `json:"old_password"`
	Password    string `json:"password"`
}

type User struct {
	Name           string `json:"name"`
	Username       string `json:"username"`
	Role           string `json:"role"`
	Password       string `json:"password"`
	Provinsi       string `json:"provinsi"`
	Kota           string `json:"kota"`
	OrganizationID uint    `json:"organization_id"`
	RetryAttempts  int64  `json:"retry_attempts"`
}

type UserPaged struct {
	Search string `form:"search"`
	Page   int    `form:"page"`
	Size   int    `form:"size"`
}
