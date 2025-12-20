package entity

import ()

type Countryinfo struct {
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
