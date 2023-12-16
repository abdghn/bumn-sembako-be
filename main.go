/*
 * Created on 15/09/23 01.33
 *
 * Copyright (c) 2023 Abdul Ghani Abbasi
 */

package main

import (
	"bumn-sembako-be/config"
	participantHandler "bumn-sembako-be/handler/participant"
	regionHandler "bumn-sembako-be/handler/region"
	userHandler "bumn-sembako-be/handler/user"
	participantService "bumn-sembako-be/service/participant"
	regionService "bumn-sembako-be/service/region"
	userService "bumn-sembako-be/service/user"
	participantUsecase "bumn-sembako-be/usecase/participant"
	regionUsecase "bumn-sembako-be/usecase/region"
	userUsecase "bumn-sembako-be/usecase/user"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	viper.SetConfigFile(".env")

	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Errorf("error: %v", err)
		return
	}

	port := viper.Get("PORT").(string)
	dbUser := viper.Get("DB_USER").(string)
	dbPass := viper.Get("DB_PASSWORD").(string)
	dbHost := viper.Get("DB_HOST").(string)
	dbPort := viper.Get("DB_PORT").(string)
	dbName := viper.Get("DB_NAME").(string)

	db := config.DbConnect(dbUser, dbPass, dbHost, dbPort, dbName)

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3004", "http://localhost:8080", "*"}, AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	v1 := router.Group("/bumn-sembako/api/v1")

	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	us := userService.NewService(db)
	uu := userUsecase.NewUsecase(us)
	uh := userHandler.NewHandler(uu)

	rs := regionService.NewService(db)
	ru := regionUsecase.NewUsecase(rs)
	rh := regionHandler.NewHandler(ru)

	ps := participantService.NewService(db)
	pu := participantUsecase.NewUsecase(ps, rs)
	ph := participantHandler.NewHandler(pu)

	router.StaticFS("/bumn-sembako/api/template", http.Dir("templates"))
	v1.StaticFS("/image", http.Dir("uploads"))

	v1.POST("login", uh.Login)
	v1.POST("register", uh.Register)
	v1.POST("register-yayasan", uh.RegisterYayasan)
	v1.GET("dashboard", ph.ViewDashboard)
	v1.GET("excel", ph.ExportExcel)
	v1.POST("report/export", ph.ExportReport)
	v1.POST("report/export-new", ph.ExportReportV2)
	v1.GET("photo/:path", ph.ImageHandler)
	v1.GET("photobase64/:path", ph.ImageBase64Handler)
	//
	user := v1.Group("/user")
	{
		user.GET("", uh.ViewUsers)
		user.POST("", uh.CreateUser)
		user.PUT("/:id", uh.UpdateUser)
		user.DELETE("/:id", uh.DeleteUser)
		user.GET("organization", uh.ViewOrganizations)
		user.GET("organization/eo", uh.ViewEOOrganizations)
		user.GET("organization/yayasan", uh.ViewYayasanOrganizations)
	}

	participant := v1.Group("/participant")
	{
		participant.GET("", ph.ViewParticipants)
		participant.GET("/:id", ph.ViewParticipant)
		participant.PUT("/:id", ph.Update)
		participant.PUT("/edit/:id", ph.Edit)
		participant.POST("import", ph.BulkCreate)
		participant.GET("import", ph.ViewLogs)
		participant.PUT("/reset/:id", ph.Reset)
		participant.DELETE("/:id", ph.Delete)
	}

	region := v1.Group("/region")
	{
		region.GET("province", rh.ViewProvincies)
		region.GET("regency", rh.ViewRegenciesByProvinceId)
		region.GET("district", rh.ViewDistrictsByRegencyId)
		region.GET("village", rh.ViewVillagesByDistrictId)
	}

	err = router.Run(":" + port)
	if err != nil {
		fmt.Errorf("error: %v", err)
	}

}
