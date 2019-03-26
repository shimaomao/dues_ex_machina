package main

import (
	"fmt"
	"go-contacts/src/controllers"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	router.POST("/api/user/register", controllers.CreateAccount)
	router.POST("/api/user/login", controllers.Authenticate)

	router.POST("/api/user/resetpassword", controllers.ResetPassword)
	port := os.Getenv("PORT")

	if port == "" {
		port = "8004" //localhost
	}

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}

}
