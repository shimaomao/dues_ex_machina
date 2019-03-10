package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-contacts/src/controllers"
	"net/http"
	"os"
)

func main() {
	router := gin.Default()
	router.POST("/api/user/register", controllers.CreateAccount)
	router.POST("/api/user/login", controllers.Authenticate)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8001" //localhost
	}

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}

}
