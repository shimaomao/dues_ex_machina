package models

import (
	"fmt"
	u "go-contacts/src/utils"
	"time"

	"github.com/jinzhu/gorm"
	"gopkg.in/go-playground/validator.v9"
)

type MarketOrder struct {
	CreatedAt  time.Time
	UpdatedAt  time.Time
	ID         uint            `gorm:"primary_key"`
	Email      string          `json:"email" validate:"required" `
	SecretName string          `json:"secretName" validate:"required" `
	Quantity   float64 `json:"quantity" validate:"required" gorm:"type:numeric"`
	ApiKey     string          `json:"apiKey" validate:"required"`
	SecretKey  string          `json:"secretKey" validate:"required"`
	Platform   string          `json:"platform" validate:"required"`
	Indicator  string          `json:"indicator" validate:"required"`
	Period     int             `json:"period" validate:"required"`
	Symbol     string          `json:"symbol" validate:"required"`
	Side       string          `json:"side" validate:"required"`
	Threshold  float64 `json:"threshold" gorm:"type:numeric" validate:"required"`
}

func (marketOrder *MarketOrder) AddMarketOrder() map[string]interface{} {
	validate = validator.New()
	validateErr := validate.Struct(marketOrder)
	if validateErr != nil {

		if _, ok := validateErr.(*validator.InvalidValidationError); ok {
			fmt.Println(validateErr)
			return u.Message(false, validateErr.Error())

		}
		return u.Message(false, validateErr.Error())

	}
	if err := GetDB().Create(marketOrder).Error; err != nil {
		return u.Message(false, err.Error())

	}
	response := u.Message(true, "MarketOrder has been created")
	response["marketOrder"] = marketOrder
	return response
}

func (account *UserInfo) GetAllMarketOrder() (bool, []MarketOrder) {
	marketOrder := []MarketOrder{}
	fmt.Println("received account email ", account.Email)
	err := GetDB().Table("market_orders").Where("email = ?", account.Email).Find(&marketOrder).Error
	if err == gorm.ErrRecordNotFound {
		fmt.Println("Record not found")
		return false, marketOrder
	}
	return true, marketOrder

}

func GetAllMarketOrderBySideAndIndicator(side string, indicator string)(bool, []MarketOrder){
	marketOrder := []MarketOrder{}
	
	err := GetDB().Table("market_orders").Where("indicator = ? AND side = ? ", indicator, side).Find(&marketOrder).Error
	if err == gorm.ErrRecordNotFound {
		fmt.Println("Record not found")
		return false, marketOrder
	}
		return true, marketOrder

}

func CheckIndicator(indicator string, side string, symbol string) (bool, []MarketOrder) {
	marketOrder := []MarketOrder{}
	err := GetDB().Table("market_orders").Where("indicator = ? AND side = ? AND symbol = ?", indicator, side, symbol).Find(&marketOrder).Error
	if err == gorm.ErrRecordNotFound {
		fmt.Println("Record not found")
		return false, marketOrder
	}
	return true, marketOrder
}
