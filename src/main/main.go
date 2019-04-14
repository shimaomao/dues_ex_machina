package main

import (
	"go-contacts/src/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	router.POST("/api/user/register", controllers.CreateAccount)
	router.POST("/api/user/login", controllers.Authenticate)
	router.GET("/api/test", controllers.TestCORS)
	router.POST("/api/user/resetpassword", controllers.ResetPassword)
	router.POST("/api/user/deleteuser", controllers.DeleteAccount)
	router.POST("/api/user/getuser", controllers.GetAccount)
	router.POST("/api/user/addusersecret", controllers.AddClientSecrets)
	router.GET("/api/user/confirm", controllers.VerifyEmail)
	router.POST("/api/user/deleteusersecret", controllers.DeleteClientSecret)
	router.POST("/api/user/getallusersecret", controllers.GetAllClientSecret)
	MarketWatch()
	router.Run()

}
