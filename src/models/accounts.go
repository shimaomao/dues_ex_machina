package models

import (
	"fmt"
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

//a struct to rep user account

type UserInfo struct {
	Email              string    `json:"email" validate:"required,email,min=6,contains=@" gorm:"type:varchar(50);primary_key" `
	Password           string    `json:"password" validate:"required,min=8"`
	Username           string    `json:"username" validate:"required" gorm:"type:varchar(100)"`
	ID                 uint      `gorm:"AUTO_INCREMENT"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"update_at"`
	SecretKey          string    `json:"secretkey"`
	ApiKey             string    `json:"apikey"`
	EmailVerified      bool      `gorm:"default:'false'"`
	OperationSecretKey string    `gorm:"default:'-'"`
}

var validate *validator.Validate

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

func (account *UserInfo) Create() map[string]interface{} {

	if resp, ok := account.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)
	// verified := Account{Email: account.Email, Password: account.Password, Username: account.Username}
	// temp := &UserInfo{}
	// temp.Email = account.Email
	// temp.Password = account.Password
	// temp.Username = account.Username
	GetDB().Create(account)

	if account.ID <= 0 {
		return u.Message(false, "Failed to create account, connection error.")
	}
	//Create new JWT token for the newly registered account
	// claims := make(jwt.MapClaims)
	// claims["exp"] = time.Now().Add(time.Hour * 10)
	// tk := &Token{Email: account.Email, StandardClaims: claims}
	// secretInfo := jwt.MapClaims{"secret": "hehe", "nbf": time.Now()}
	// token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), secretInfo)
	token := jwt.New(jwt.SigningMethodRS512)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 12)
	claims["iat"] = time.Now()
	token.Claims = claims
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	// account.Token = tokenString

	account.Password = "" //delete password

	response := u.Message(true, "Account has been created")
	response["account"] = account
	response["token"] = tokenString
	return response
}

func Login(email, password string) map[string]interface{} {

	account := &UserInfo{}
	err := GetDB().Table("user_infos").Where("email = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return u.Message(false, "Invalid login credentials. Please try again")
	}
	// temp.account = account
	//Worked! Logged In
	account.Password = ""

	//Create JWT token
	tk := &Token{Email: account.Email}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	// account.Token = tokenString //Store the token in the response

	resp := u.Message(true, "Logged In")
	resp["token"] = tokenString

	return resp
}

func GetUser(u uint) *UserInfo {

	acc := &UserInfo{}
	GetDB().Table("user_infos").Where("id = ?", u).First(acc)
	if acc.Email == "" { //User not found!
		return nil
	}

	acc.Password = ""
	return acc
}
