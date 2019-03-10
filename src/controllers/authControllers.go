package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	hermes "go-contacts/src/email/examples/send"
	"go-contacts/src/models"
	"net/http"
)

func CreateAccount(c *gin.Context) {
	cCp := c.Copy()
	result := make(chan gin.H)

	go func() {
		account := &models.Account{}
		if err := cCp.ShouldBindJSON(&account); err != nil {
			cCp.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// u1 := uuid.Must(uuid.NewV4())
		u1 := uuid.NewV4()
		if len(u1.String()) < 16 {
			result <- gin.H{"status": "uuid creation has failed"}
			return
		}
		resp := account.Create()
		if resp["status"] == false {
			result <- gin.H{"status": "acc creation has failed", "data": resp}
			return
		}
		result <- gin.H{"status": "new acc created successfully", "data": resp}
		// cCp.JSON(http.StatusOK, gin.H{"status": "you are logged in", "data": resp})
	}()
	go hermes.SendEmailVerification()
	c.JSON(http.StatusOK, <-result)

}

// var CreateAccount = func(w http.ResponseWriter, r *http.Request) {

// 	account := &models.Account{}
// 	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
// 	if err != nil {
// 		u.Respond(w, u.Message(false, "Invalid request"))
// 		return
// 	}

// 	resp := account.Create() //Create account
// 	u.Respond(w, resp)
// }

var Authenticate = func(c *gin.Context) {
	cCp := c.Copy()
	result := make(chan gin.H)
	go func() {
		account := &models.Account{}

		if err := cCp.ShouldBindJSON(&account); err != nil {
			cCp.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp := models.Login(account.Email, account.Password)
		result <- gin.H{"status": "you are logged in", "data": resp}

	}()
	c.JSON(http.StatusOK, <-result)

}
