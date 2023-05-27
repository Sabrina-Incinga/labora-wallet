package services

import (
	"database/sql"
	"fmt"

	"github.com/labora-wallet/walletAPI/model"
)

type PostgresCustomerDBHandler struct {
	Db *sql.DB
}

func (p *PostgresCustomerDBHandler) CreateCustomer(customer model.CustomerDTO, tx *sql.Tx) (int64, error) {
	var id int64
	var row *sql.Row

	existentCustomer, err := p.GetCustomerByIdentityNumber(customer.NationalIdentityNumber, customer.NationalIdentityType, customer.CountryId)
	if err != nil {
		return id, err
	} else if existentCustomer != nil {
		return existentCustomer.ID, nil
	}

	query := createCustomerQuery(customer)

	if tx != nil {
		row = tx.QueryRow(query)
	}else{
		row = p.Db.QueryRow(query)
	}

	err = row.Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (p *PostgresCustomerDBHandler) GetCustomerByIdentityNumber(nationalIdentityNumber, nationalIdentityType, countryId string) (*model.Customer, error) {
	row := p.Db.QueryRow(`SELECT
							id 
							, first_name
							, last_name
							, national_identity_number
							, national_identity_type
							, country_id
							FROM public.customer WHERE national_identity_number=$1 AND national_identity_type=$2 AND country_id=$3;`, nationalIdentityNumber, nationalIdentityType, countryId)

	var customer model.Customer

	err := row.Scan(&customer.ID, &customer.FirstName, &customer.LastName, &customer.NationalIdentityNumber, &customer.NationalIdentityType, &customer.CountryId)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &customer, nil
}

func (p *PostgresCustomerDBHandler) GetCustomerById(id int64) (*model.Customer, error) {
	row := p.Db.QueryRow(`SELECT
							id 
							, first_name
							, last_name
							, national_identity_number
							, national_identity_type
							, country_id
							FROM public.customer WHERE id=$1;`, id)

	var customer model.Customer

	err := row.Scan(&customer.ID, &customer.FirstName, &customer.LastName, &customer.NationalIdentityNumber, &customer.NationalIdentityType, &customer.CountryId)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &customer, nil
}

func (p *PostgresCustomerDBHandler) UpdateCustomer(dto model.CustomerDTO, id int64, tx *sql.Tx) (int64, error) {
	var rowsAffected int64
	var err error
	var response sql.Result
	query := updateCustomerQuery(dto, id)

	if tx != nil {
		response, err = tx.Exec(query)
	}else{
		response, err = p.Db.Exec(query)
	}

	if err != nil {
		return rowsAffected, err
	}

	rowsAffected, err = response.RowsAffected()
	if err != nil {
		return rowsAffected, err
	}

	return rowsAffected, nil
}

func (p *PostgresCustomerDBHandler) DeleteCustomer(id int64, tx *sql.Tx) (int64, error) {
	var rowsAffected int64
	var err error
	var response sql.Result
	query := deleteCustomerQuery(id)

	if tx != nil {
		response, err = tx.Exec(query)
	}else{
		response, err = p.Db.Exec(query)
	}

	if err != nil {
		return rowsAffected, err
	}

	rowsAffected, err = response.RowsAffected()
	if err != nil {
		return rowsAffected, err
	}

	return rowsAffected, nil
}

func createCustomerQuery(customer model.CustomerDTO)string{
	query := fmt.Sprintf(`INSERT INTO public.customer(
							first_name
							, last_name
							, national_identity_number
							, national_identity_type
							, country_id)
							VALUES ('%s', '%s', '%s', '%s', '%s') RETURNING id;`, customer.FirstName, customer.LastName, customer.NationalIdentityNumber, customer.NationalIdentityType, customer.CountryId)

	return query
}

func updateCustomerQuery(customer model.CustomerDTO, id int64) string {
	query := fmt.Sprintf(`UPDATE public.customer 
						first_name='%s'
						, last_name='%s'
						, national_identity_number='%s'
						, national_identity_type='%s'
						, country_id='%s'
						WHERE id=%d;`, customer.FirstName, customer.LastName, customer.NationalIdentityNumber, customer.NationalIdentityType, customer.CountryId, id)

	return query
}

func deleteCustomerQuery(id int64) string {
	query := fmt.Sprintf(`DELETE FROM public.customer
						WHERE id=%d;`, id)

	return query
}