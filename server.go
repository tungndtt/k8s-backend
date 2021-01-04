package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	port           string = ":8080"
	serverCertPath string = "./serverCerts/server.crt"
	serverKeyPath  string = "./serverCerts/server.key"
	filePath       string = "./files/"
)

func InitRoute() {
	router := gin.Default()

	// define all server routes

	router.GET("/", test)

	// default memory limit is 32 MiB
	router.POST("/upload", upload)

	router.RunTLS(port, serverCertPath, serverKeyPath)
}

// define all handling functions

func test(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"info": "welcome to server",
	})
}

// handle upload file
func upload(c *gin.Context) {
	// single file
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println(err)
	}
	err = c.SaveUploadedFile(file, filePath)
	if err != nil {
		fmt.Println(err)
	}
}
