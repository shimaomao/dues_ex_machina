package controllers

import (
	"fmt"
	hermes "go-contacts/src/email/examples/send"
	"go-contacts/src/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResetPassword(c *gin.Context) {
	go hermes.SendEmailResetPassword()

}
func TestCORS(c *gin.Context) {
	c.JSON(http.StatusOK, "Done")

}

func CreateAccount(c *gin.Context) {
	cCp := c.Copy()
	fmt.Println(cCp)

	result := make(chan gin.H)
	go func() {
		account := &models.UserInfo{}
		if err := cCp.ShouldBindJSON(&account); err != nil {
			fmt.Println(cCp)

			result <- gin.H{"status": false, "error": err.Error()}
			return
		}
		resp := account.Create()
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

var DeleteAccount = func(c *gin.Context) {
	cCp := c.Copy()
	result := make(chan gin.H)
	go func() {
		account := &models.UserInfo{}
		if err := cCp.ShouldBindJSON(&account); err != nil {
			result <- gin.H{"status": false, "error": err.Error()}
			return
		}
		resp := account.DeleteUser()
		if resp["status"] == false {
			result <- gin.H{"status": false, "data": resp}

		} else {
			result <- gin.H{"status": true, "data": resp}

		}

	}()
	resultedData := <-result
	if resultedData["status"] == false {
		c.JSON(http.StatusBadRequest, resultedData)
		return
	} else {

		c.JSON(http.StatusOK, resultedData)

	}

}
var AddClientSecrets = func(c *gin.Context) {
	cCp := c.Copy()
	result := make(chan gin.H)
	go func() {
		secrets := &models.ClientSecrets{}
		if err := cCp.ShouldBindJSON(&secrets); err != nil {
			result <- gin.H{"status": false, "error": err.Error()}

			// cCp.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp := secrets.AddClientSecrets()
		if resp["status"] == false {
			result <- gin.H{"status": false, "data": resp}

		} else {
			result <- gin.H{"status": true, "data": resp}

		}
	}()
	resultedData := <-result
	if resultedData["status"] == false {
		c.JSON(http.StatusBadRequest, resultedData)
		return
	} else {

		c.JSON(http.StatusOK, resultedData)

	}
}

var GetAccount = func(c *gin.Context) {
	cCp := c.Copy()
	result := make(chan gin.H)
	go func() {
		account := &models.UserInfo{}

		if err := cCp.ShouldBindJSON(&account); err != nil {
			result <- gin.H{"status": false, "error": err.Error()}

			// cCp.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		isFound, resp := models.GetUser(account.Email)
		result <- gin.H{"status": isFound, "data": resp}

	}()
	c.JSON(http.StatusOK, <-result)

}

var VerifyEmail = func(c *gin.Context) {
	cCp := c.Copy()
	result := make(chan gin.H)
	go func() {
		//account := &models.EmailVerification{}

		// if err := cCp.ShouldBind(&account); err != nil {
		// 	result <- gin.H{"status": false, "error": err.Error()}

		// 	// cCp.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 	return
		// }
		queryStringParam := cCp.Request.URL.Query()
		email := queryStringParam["email"]
		token := queryStringParam["token"]
		if email == nil || token == nil {
			resp := make(map[string]interface{})
			resp["status"] = false
			result <- gin.H{"data": resp}
			return
		}

		if len(email) > 0 && len(token) > 0 {
			resp := models.CheckOperationSecretKey(email[0], token[0])
			result <- gin.H{"data": resp}

		} else {
			var resp map[string]interface{}
			resp["status"] = false
			result <- gin.H{"data": resp}
		}

	}()

	resultedData := <-result
	if resultedData["status"] == false {
		c.JSON(http.StatusBadRequest, resultedData)
		return
	} else {

		c.JSON(http.StatusOK, resultedData)

	}
}

var GetAllClientSecret = func(c *gin.Context) {
	cCp := c.Copy()
	result := make(chan gin.H)
	go func() {
		secret := &models.ClientSecrets{}
		if err := cCp.ShouldBindJSON(&secret); err != nil {
			result <- gin.H{"status": false, "error": err.Error()}

			// cCp.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		isSuccess, resp := secret.GetAllClientSecrets()
		if isSuccess == false {
			result <- gin.H{"status": false, "error": "No secrets found under this email"}
			return
		}
		result <- gin.H{"status": true, "error": resp}

	}()
	resultedData := <-result
	if resultedData["status"] == false {
		c.JSON(http.StatusBadRequest, resultedData)
		return
	} else {

		c.JSON(http.StatusOK, resultedData)

	}
}

var DeleteClientSecret = func(c *gin.Context) {
	cCp := c.Copy()
	result := make(chan gin.H)
	go func() {
		secrets := &models.ClientSecrets{}
		if err := cCp.ShouldBindJSON(&secrets); err != nil {
			result <- gin.H{"status": false, "error": err.Error()}

			// cCp.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resp := secrets.DeleteClientSecrets()
		if resp["status"] == false {
			result <- gin.H{"status": false, "data": resp}

		} else {
			result <- gin.H{"status": true, "data": resp}

		}

	}()
	resultedData := <-result
	if resultedData["status"] == false {
		c.JSON(http.StatusBadRequest, resultedData)
		return
	} else {

		c.JSON(http.StatusOK, resultedData)

	}
}

var Authenticate = func(c *gin.Context) {
	cCp := c.Copy()
	result := make(chan gin.H)
	go func() {
		account := &models.UserInfo{}

		if err := cCp.ShouldBindJSON(&account); err != nil {
			result <- gin.H{"status": false, "error": err.Error()}

			// cCp.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp := models.Login(account.Email, account.Password)
		result <- gin.H{"data": resp}

	}()
	c.JSON(http.StatusOK, <-result)

}
