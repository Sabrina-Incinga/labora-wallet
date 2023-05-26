package models

type Customer struct{
	ID            int    			`json:"id"`
	FirstName    string     		`json:"first_name"`
	CreationDate string 			`json:"creation_date"`
	NationalIdentityNumber string   `json:"national_identity_number"`
	NationalIdentityType   string   `json:"national_identity_type"`
	CountryId    string     		`json:"country_id"`
}