package model

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
	"github.com/pkg/errors"
	// "DTS_IT_Perbankan_Back_End/Digitalent-Kominfo_Implementation-MVC-Golang/app/constant"
	"DTS_IT_Perbankan_Back_End/Digitalent-Kominfo_Implementation-MVC-Golang/app/utils"
)

type Account struct {
	DB *gorm.DB
	ID int `gorm:"primary_key" json:"-"`
	IdAccount string `json:"id_account, omitempty"`
	Name string `json:"name"`
	Password string `json:"password, omitempty"`
	AccountNumber int `json:"account_number"`
	Saldo int `json:"saldo"`
	Transaction []Transaction `gorm:"ForeignKey:IdAccountRefer" json:"transaction"`
}

type Auth struct {
	Name string `json:"name"`
	Password string `json:"password"`
}

type Transaction struct {
	ID int `gorm:"primary_key" json:"-"`
	IdAccountRefer int `json:"-"`
	IdTransaction string `json:"id_transaction"`
	TransactionType int `json:"transaction_type"`
	TransactionDescription string `json:"transaction_description"`
	Sender int `json:"sender"`
	Recipient int `json:"recipient"`
	Timestamp int `json:"timestamp"`
}

func Login(auth Auth) (bool, error, string) {
	var account Account
	if err := account.DB.Where(&Account{Name: auth.Name}).First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound{
			return false, errors.Errorf("Account not found"), ""
		}
	}

	err := utils.HashComparator([]byte(account.Password),[]byte(auth.Password))
	if err != nil{
		return false, errors.Errorf("Incorrect Password"),""
	} else {
		sign := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
			"name":auth.Name,
			"account_number":account.AccountNumber,
		})

		token ,err := sign.SignedString([]byte("secret"))
		if err != nil {
			return false,err,""
		}
		return true,nil,token
	}
}

func InsertNewAccount(account Account) (bool, error) {
	account.AccountNumber = utils.RangeIn(111111,999999)
	account.Saldo = 0
	account.IdAccount = fmt.Sprintf("id-%id", utils.RangeIn(111,999))

	result := account.DB.Create(&account)
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func (account Account) GetAccountDetail(idAccount int) (bool,error, []Transaction,Account){
	var transaction []Transaction

	if err := account.DB.Where("sender = ? OR recipient = ?",idAccount,idAccount).
		Find(&transaction).Error;err!=nil{
		if err == gorm.ErrRecordNotFound{
			return false,errors.Errorf("Account not found"), []Transaction{},Account{}
		} else {
			return false, errors.Errorf("invalid prepare statement :%+v\n", err), []Transaction{},Account{}
		}
	}

	if err := account.DB.Where(&Account{AccountNumber: idAccount}).Find(&account).Error;err != nil{
		if err == gorm.ErrRecordNotFound{
			return false,errors.Errorf("Account not found"), []Transaction{},Account{}
		} else {
			return false, errors.Errorf("invalid prepare statement :%+v\n", err), []Transaction{},Account{}
		}
	}

	return true,nil, transaction, account
}
