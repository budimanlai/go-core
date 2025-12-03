package dto

type CountryResponse struct {
	CountryName string `json:"country_name"  copier:"Name"`
	CountryCode string `json:"country_code" copier:"IsoAlpha2"`
}
