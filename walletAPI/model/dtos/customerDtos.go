package dtos

type CustomerDTO struct {
	FirstName              string `json:"first_name"`
	LastName               string `json:"last_name"`
	NationalIdentityNumber string `json:"national_identity_number"`
	NationalIdentityType   string `json:"national_identity_type"`
	CountryId              string `json:"country_id"`
}
