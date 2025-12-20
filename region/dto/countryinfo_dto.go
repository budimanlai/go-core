package dto

import ()

type CreateCountryinfoReq struct {
	IsoAlpha2          string  `json:"iso_alpha2" validate:"omitempty,max=2"`
	IsoAlpha3          *string `json:"iso_alpha3" validate:"omitempty"`
	IsoNumeric         *int    `json:"iso_numeric" validate:"omitempty"`
	FipsCode           *string `json:"fips_code" validate:"omitempty"`
	Name               *string `json:"name" validate:"omitempty"`
	Capital            *string `json:"capital" validate:"omitempty"`
	Areainsqkm         *int64  `json:"areainsqkm" validate:"omitempty"`
	Population         *int    `json:"population" validate:"omitempty"`
	Continent          *string `json:"continent" validate:"omitempty"`
	Tld                *string `json:"tld" validate:"omitempty"`
	CurrencyCode       *string `json:"currency_code" validate:"omitempty"`
	CurrencyName       *string `json:"currency_name" validate:"omitempty"`
	Phone              *string `json:"phone" validate:"omitempty"`
	PostalCodeFormat   *string `json:"postal_code_format" validate:"omitempty"`
	PostalCodeRegex    *string `json:"postal_code_regex" validate:"omitempty"`
	Languages          *string `json:"languages" validate:"omitempty"`
	GeonameId          *int    `json:"geonameId" validate:"omitempty"`
	Neighbours         *string `json:"neighbours" validate:"omitempty"`
	EquivalentFipsCode *string `json:"equivalent_fips_code" validate:"omitempty"`
	CurrencySymbol     *string `json:"currency_symbol" validate:"omitempty"`
	Status             *string `json:"status" validate:"omitempty"`
}

type UpdateCountryinfoReq struct {
	IsoAlpha2          string  `json:"iso_alpha2"`
	IsoAlpha3          *string `json:"iso_alpha3"`
	IsoNumeric         *int    `json:"iso_numeric"`
	FipsCode           *string `json:"fips_code"`
	Name               *string `json:"name"`
	Capital            *string `json:"capital"`
	Areainsqkm         *int64  `json:"areainsqkm"`
	Population         *int    `json:"population"`
	Continent          *string `json:"continent"`
	Tld                *string `json:"tld"`
	CurrencyCode       *string `json:"currency_code"`
	CurrencyName       *string `json:"currency_name"`
	Phone              *string `json:"phone"`
	PostalCodeFormat   *string `json:"postal_code_format"`
	PostalCodeRegex    *string `json:"postal_code_regex"`
	Languages          *string `json:"languages"`
	GeonameId          *int    `json:"geonameId"`
	Neighbours         *string `json:"neighbours"`
	EquivalentFipsCode *string `json:"equivalent_fips_code"`
	CurrencySymbol     *string `json:"currency_symbol"`
	Status             *string `json:"status"`
}
