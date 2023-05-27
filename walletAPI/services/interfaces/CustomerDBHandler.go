package interfaces

import (
	"database/sql"
	"github.com/labora-wallet/walletAPI/model"
)

type CustomerDBHandler interface {
	CreateCustomer(customer model.CustomerDTO, tx *sql.Tx) (int64, error)
	GetCustomerByIdentityNumber(nationalIdentityNumber, nationalIdentityType, countryId string) (*model.Customer, error)
	GetCustomerById(id int64) (*model.Customer, error)
	UpdateCustomer(dto model.CustomerDTO, id int64, tx *sql.Tx) (int64, error)
	DeleteCustomer(id int64, tx *sql.Tx) (int64, error)
}
