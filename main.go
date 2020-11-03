package main

import (
	"DTS_IT_Perbankan_Back_End/Digitalent-Kominfo_Implementation-MVC-Golang/app/config"
	"DTS_IT_Perbankan_Back_End/Digitalent-Kominfo_Implementation-MVC-Golang/app/controller"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main(){
	db := config.DBInit()
	inDB := &controller.InDB{DB: db}

	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/",inDB.CreateAccount)
	router.Run(":8080")
}

