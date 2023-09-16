/*
 * Created on 01/04/22 17.20
 *
 * Copyright (c) 2022 Abdul Ghani Abbasi
 */

package helper

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResponsePaged struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Page    int         `json:"page"`
	Size    int         `json:"size"`
	Total   int64       `json:"total"`
}

func HandleSuccess(c *gin.Context, data interface{}) {
	responseData := Response{
		Status:  "200",
		Message: "Success",
		Data:    data,
	}
	c.JSON(http.StatusOK, responseData)
}

func HandlePagedSuccess(c *gin.Context, data interface{}, page, size int, total int64) {
	responseData := ResponsePaged{
		Status:  "200",
		Message: "Success",
		Data:    data,
		Page:    page,
		Size:    size,
		Total:   total,
	}
	c.JSON(http.StatusOK, responseData)
}

func HandleError(c *gin.Context, status int, message string) {
	responseData := Response{
		Status:  strconv.Itoa(status),
		Message: message,
	}
	c.JSON(status, responseData)
}
