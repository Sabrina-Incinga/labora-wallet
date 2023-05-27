package model

type Customer struct{
	ID            		   int    	`json:"id"`
	FirstName    		   string   `json:"first_name"`
	LasttName    		   string   `json:"last_name"`
	NationalIdentityNumber string   `json:"national_identity_number"`
	NationalIdentityType   string   `json:"national_identity_type"`
	CountryId    		   string   `json:"country_id"`
}

type CustomerDTO struct{
	FirstName    		   string   `json:"first_name"`
	LasttName    		   string   `json:"last_name"`
	NationalIdentityNumber string   `json:"national_identity_number"`
	NationalIdentityType   string   `json:"national_identity_type"`
	CountryId    		   string   `json:"country_id"`
}