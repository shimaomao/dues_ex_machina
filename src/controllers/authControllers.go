package controllers

import (
	hermes "go-contacts/src/email/examples/send"
	"go-contacts/src/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResetPassword(c *gin.Context) {
	go hermes.SendEmailResetPassword()

}
func CreateAccount(c *gin.Context) {
	cCp := c.Copy()
	result := make(chan gin.H)
	go func() {
		account := &models.UserInfo{}
		if err := cCp.ShouldBindJSON(&account); err != nil {
			cCp.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resp := account.Create()
		if resp["status"] == false {
			result <- gin.H{"status": false, "data": resp}
			return
		}

		result <- gin.H{"status": true, "data": resp}
	}()
	resultedData := <-result
	if resultedData["status"] == false {
		c.JSON(http.StatusBadRequest, resultedData)
		return
	} else {
		go hermes.SendEmailVerification()

	}
	c.JSON(http.StatusOK, resultedData)

}

var Authenticate = func(c *gin.Context) {
	cCp := c.Copy()
	result := make(chan gin.H)
	go func() {
		account := &models.UserInfo{}

		if err := cCp.ShouldBindJSON(&account); err != nil {
			cCp.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp := models.Login(account.Email, account.Password)
		result <- gin.H{"status": "you are logged in", "data": resp}

	}()
	c.JSON(http.StatusOK, <-result)

}
