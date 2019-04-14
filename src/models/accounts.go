package models

import (
	"fmt"
	hermes "go-contacts/src/email/examples/send"
	u "go-contacts/src/utils"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)

/*
JWT claims struct
*/
type Token struct {
	Email string
	jwt.StandardClaims
}

type MarketOrder struct {
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Email      string  `json:"email" validate:"required" `
	SecretName string  `json:"secretName" validate:"required" `
	Price      float64 `json:"price" validate:"required"`
	ApiKey     string  `json:"apiKey" validate:"required"`
	SecretKey  string  `json:"secretKey" validate:"required"`
	Platform   string  `json:"platform" validate:"required"`
	Indicator  string  `json:"indicator" validate:"required"`
	Period     int     `json:"period" validate:"required"`
	Pair       string  `json:"pair" validate:"required"`
}

type ClientSecrets struct {
	CreatedAt  time.Time
	UpdatedAt  time.Time
	SecretName string `form :"secretName" json:"secretName" validate:"required" gorm:"primary_key"`
	ApiKey     string `json:"apiKey" validate:"required"`
	SecretKey  string `json:"secretKey" validate:"required"`
	ID         uint   `gorm:"AUTO_INCREMENT"`
	Email      string `form :"email" json:"email" validate:"required" gorm:"primary_key"`
}

type UserInfo struct {
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          *time.Time
	Email              string `json:"email" validate:"required,email,min=6,contains=@" gorm:"type:varchar(50);primary_key" `
	Password           string `json:"password" validate:"required,min=8"`
	Username           string `json:"username" gorm:"type:varchar(100)"`
	ID                 uint   `gorm:"AUTO_INCREMENT"`
	SecretKey          string `json:"secretkey"`
	ApiKey             string `json:"apikey"`
	EmailVerified      bool   `gorm:"default:'false'"`
	OperationSecretKey string `gorm:"default:'-'"`
}

type EmailVerification struct {
	Token string `form:"token"`
	Email string `form:"email"`
}

var validate *validator.Validate

func generateJWT(Email string) *Token {
	timeToExpire := time.Now().Add(time.Hour).Unix()
	claims := &Token{
		Email,
		jwt.StandardClaims{
			ExpiresAt: timeToExpire,
			Issuer:    "verifyEmail",
		},
	}
	return claims
}

//Validate incoming user details...
func (account *UserInfo) Validate() (map[string]interface{}, bool) {

	validate = validator.New()
	validateErr := validate.Struct(account)
	//fmt.Println(validateErr)
	if validateErr != nil {
		fmt.Println("In rejection")

		if _, ok := validateErr.(*validator.InvalidValidationError); ok {
			fmt.Println(validateErr)
			return u.Message(false, validateErr.Error()), false

		}
		return u.Message(false, validateErr.Error()), false

	}

	//Email must be unique
	temp := &UserInfo{}

	//check for errors and duplicate emails
	err := GetDB().Table("user_infos").Where("email = ?", account.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}
	if temp.Email != "" {
		return u.Message(false, "Email address already in use by another user."), false
	}

	return u.Message(false, "Requirement passed"), true
}

func (client *ClientSecrets) DeleteClientSecrets() map[string]interface{} {
	validate = validator.New()
	validateErr := validate.Struct(client)
	if validateErr != nil {
		fmt.Println("In rejection")

		if _, ok := validateErr.(*validator.InvalidValidationError); ok {
			fmt.Println(validateErr)
			return u.Message(false, validateErr.Error())

		}
		return u.Message(false, validateErr.Error())

	}

	err := GetDB().Table("client_secrets").Where("email = ? AND secret_name = ?", client.Email, client.SecretName).Delete(client).Error
	if err != nil {
		return u.Message(false, err.Error())
	}
	return u.Message(true, "secrets deleted successfully")

}

func (account *UserInfo) Create() map[string]interface{} {

	if resp, ok := account.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)
	claims := generateJWT(account.Email)

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)

	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))

	account.OperationSecretKey = tokenString

	if err := GetDB().Create(account).Error; err != nil {
		return u.Message(false, err.Error())

	}

	// if account.ID <= 0 {
	// 	return u.Message(false, "Failed to create account. ")
	// }
	go hermes.SendEmailVerification(account.Email, account.OperationSecretKey)

	account.Password = "" //delete password

	account.OperationSecretKey = ""
	response := u.Message(true, "Account has been created")
	response["account"] = account
	return response
}

func CheckOperationSecretKey(email string, token string) map[string]interface{} {
	isFound, account := GetUser(email)
	if isFound == false {
		return u.Message(false, "account not found")

	}
	if account.EmailVerified == false && account.OperationSecretKey == token {
		account.OperationSecretKey = ""
		account.EmailVerified = true
	} else {
		return u.Message(false, "account email has already been verified")

	}

	err := GetDB().Table("user_infos").Where("email = ?", email).Update(UserInfo{EmailVerified: true, OperationSecretKey: "-"}).Error
	if err != nil {
		return u.Message(false, err.Error())

	}
	return u.Message(true, "account email successfully verified")

}

func (client *ClientSecrets) GetAllClientSecrets() (bool, []ClientSecrets) {
	secrets := []ClientSecrets{}
	err := GetDB().Table("client_secrets").Where("email = ? ", client.Email).Find(&secrets).Error
	if err != nil {
		return false, secrets

	}
	return true, secrets

}

func Login(email, password string) (bool, map[string]interface{}) {

	// account := UserInfo{}
	isFound, account := GetUser(email)
	if isFound == false {
		return false, u.Message(false, "account not found")

	}
	fmt.Println("email stored password : ", account.Password)
	fmt.Println("req password : ", password)

	err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return false, u.Message(false, "Invalid login credentials. Please try again")
	}
	// temp.account = account
	//Worked! Logged In
	account.Password = ""

	//Create JWT token
	claims := generateJWT(account.Email)

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	// account.Token = tokenString //Store the token in the response

	resp := u.Message(true, "Logged In")
	resp["token"] = tokenString

	return true, resp
}

func GetUser(email string) (bool, UserInfo) {
	respAcc := &UserInfo{}
	outputAcc := *respAcc

	err := GetDB().Table("user_infos").Where("email = ?", email).First(&outputAcc).Error

	if err == gorm.ErrRecordNotFound {
		fmt.Println("Record not found")
		return false, outputAcc
	}
	//copy a new struct
	// outputAcc.Password = ""
	return true, outputAcc
}

func (client *ClientSecrets) AddClientSecrets() map[string]interface{} {
	validate = validator.New()
	validateErr := validate.Struct(client)
	if validateErr != nil {
		fmt.Println("In rejection")

		if _, ok := validateErr.(*validator.InvalidValidationError); ok {
			fmt.Println(validateErr)
			return u.Message(false, validateErr.Error())

		}
		return u.Message(false, validateErr.Error())

	}
	isFound, _ := GetUser(client.Email)
	if isFound == false {
		return u.Message(false, "Account email not found")

	}
	err := GetDB().Table("client_secrets").Create(client).Error
	if err != nil {
		return u.Message(false, err.Error())
	}
	return u.Message(true, "secrets added successfully")

}

func (account *UserInfo) DeleteUser() map[string]interface{} {

	validate = validator.New()
	validateErr := validate.Struct(account)
	//fmt.Println(validateErr)
	if validateErr != nil {
		fmt.Println("In rejection")

		if _, ok := validateErr.(*validator.InvalidValidationError); ok {
			fmt.Println(validateErr)
			return u.Message(false, validateErr.Error())

		}
		return u.Message(false, validateErr.Error())

	}

	isFound, queryResp := GetUser(account.Email)
	if isFound == false {
		return u.Message(false, "Account not found")

	}
	fmt.Println(queryResp)
	err := GetDB().Delete(account).Error
	if err != nil {
		return u.Message(false, err.Error())

	}
	response := u.Message(true, "Account has been deleted")

	//response := make(map[string]interface{})

	return response
}
