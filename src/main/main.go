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
	// port := os.Getenv("PORT")

	// if port == "" {
	// 	port = "8007" //localhost
	// }

	// fmt.Println(port)

	router.Run()

}
