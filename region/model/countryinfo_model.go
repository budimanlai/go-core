package model

import ()

type CountryinfoModel struct {
	IsoAlpha2          string  `gorm:"column:iso_alpha2;primaryKey"`
	IsoAlpha3          *string `gorm:"column:iso_alpha3"`
	IsoNumeric         *int    `gorm:"column:iso_numeric"`
	FipsCode           *string `gorm:"column:fips_code"`
	Name               *string `gorm:"column:name"`
	Capital            *string `gorm:"column:capital"`
	Areainsqkm         *int64  `gorm:"column:areainsqkm"`
	Population         *int    `gorm:"column:population"`
	Continent          *string `gorm:"column:continent"`
	Tld                *string `gorm:"column:tld"`
	CurrencyCode       *string `gorm:"column:currency_code"`
	CurrencyName       *string `gorm:"column:currency_name"`
	Phone              *string `gorm:"column:phone"`
	PostalCodeFormat   *string `gorm:"column:postal_code_format"`
	PostalCodeRegex    *string `gorm:"column:postal_code_regex"`
	Languages          *string `gorm:"column:languages"`
	GeonameId          *int    `gorm:"column:geonameId"`
	Neighbours         *string `gorm:"column:neighbours"`
	EquivalentFipsCode *string `gorm:"column:equivalent_fips_code"`
	CurrencySymbol     *string `gorm:"column:currency_symbol"`
	Status             *string `gorm:"column:status"`
}

func (CountryinfoModel) TableName() string {
	return "countryinfo"
}
