package interfaces

import (
	"database/sql"

	"github.com/labora-wallet/walletAPI/model"
	"github.com/labora-wallet/walletAPI/model/dtos"
)

type CustomerDBHandler interface {
	CreateCustomer(customer dtos.CustomerDTO, tx *sql.Tx) (int64, error)
	GetCustomerByIdentityNumber(nationalIdentityNumber, nationalIdentityType, countryId string) (*model.Customer, error)
	GetCustomerById(id int64) (*model.Customer, error)
	UpdateCustomer(dto dtos.CustomerDTO, id int64, tx *sql.Tx) (int64, error)
	DeleteCustomer(id int64, tx *sql.Tx) (int64, error)
}
