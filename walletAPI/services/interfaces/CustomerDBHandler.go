package interfaces

import (
	"github.com/labora-wallet/walletAPI/model"
)

type CustomerDBHandler interface{
	CreateCustomer(customer model.CustomerDTO) (int64, error)
	GetCustomerByIdentityNumber(nationalIdentityNumber, nationalIdentityType, countryId string) (*model.Customer, error)
	GetCustomerById(id int) (*model.Customer, error)
	UpdateCustomer(dto model.CustomerDTO, id int) (int64, error)
	DeleteCustomer(id int) (int64, error)
}