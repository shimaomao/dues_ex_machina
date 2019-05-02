package models

import (
	"fmt"
	u "go-contacts/src/utils"
)

type PlatformSymbols struct {
	ID       uint   `gorm:"AUTO_INCREMENT"`
	Symbol   string ` gorm:"primary_key""`
	Platfrom string ` gorm:"primary_key"""`
}

func (symbol *PlatformSymbols) AddSymbols() map[string]interface{} {
	if err := GetDB().Create(symbol).Error; err != nil {
		return u.Message(false, err.Error())

	}
	return u.Message(true, "Successfully Added")

}

func GetAllSymbols(platform string) []PlatformSymbols {
	allSymbols := []PlatformSymbols{}
	err := GetDB().Table("platform_symbols").Where("platfrom = ?", platform).Find(&allSymbols).Error
	if err != nil {
		fmt.Println(err.Error())
	}
	return allSymbols
}
