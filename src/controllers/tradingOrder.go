package controllers

import (
	"fmt"

	"go-contacts/src/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateMarketOrder(c *gin.Context) {
	cCp := c.Copy()
	result := make(chan gin.H)
	go func() {
		marketOrder := &models.MarketOrder{}
		if err := cCp.ShouldBindJSON(&marketOrder); err != nil {
			fmt.Println(cCp)

			result <- gin.H{"status": false, "error": err.Error()}
			return
		}
		resp := marketOrder.AddMarketOrder()
		if resp["status"] == false {
			result <- gin.H{"status": false, "data": resp}
			return
		}
		fmt.Println(resp)
		// go hermes.SendEmailVerification(account.Email)

		result <- gin.H{"status": true, "data": resp}
	}()
	resultedData := <-result
	if resultedData["status"] == false {
		c.JSON(http.StatusBadRequest, resultedData)
		return
	} else {

		c.JSON(http.StatusOK, resultedData)

	}

}

func GetAllMarketOrder(c *gin.Context) {
	cCp := c.Copy()
	result := make(chan gin.H)
	go func() {
		userInfo := &models.UserInfo{}
		if err := cCp.ShouldBindJSON(&userInfo); err != nil {
			fmt.Println(cCp)

			result <- gin.H{"status": false, "error": err.Error()}
			return
		}
		orderFound, resp := userInfo.GetAllMarketOrder()
		if orderFound == false {
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

		c.JSON(http.StatusOK, resultedData)

	}

}
