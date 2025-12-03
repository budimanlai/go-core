package models

type Country struct {
	IsoAlpha2          string `gorm:"column:iso_alpha2;type:varchar(2);primaryKey;not null;default:''" json:"iso_alpha2"`
	IsoAlpha3          string `gorm:"column:iso_alpha3;type:varchar(3)" json:"iso_alpha3"`
	IsoNumeric         *int   `gorm:"column:iso_numeric;type:int(11)" json:"iso_numeric"`
	FipsCode           string `gorm:"column:fips_code;type:varchar(3)" json:"fips_code"`
	Name               string `gorm:"column:name;type:varchar(200);index" json:"name"`
	Capital            string `gorm:"column:capital;type:varchar(200)" json:"capital"`
	AreaInSqKm         *int64 `gorm:"column:areainsqkm;type:bigint(20)" json:"areainsqkm"`
	Population         *int   `gorm:"column:population;type:int(11)" json:"population"`
	Continent          string `gorm:"column:continent;type:char(2)" json:"continent"`
	Tld                string `gorm:"column:tld;type:varchar(4)" json:"tld"`
	CurrencyCode       string `gorm:"column:currency_code;type:char(3)" json:"currency_code"`
	CurrencyName       string `gorm:"column:currency_name;type:varchar(32)" json:"currency_name"`
	Phone              string `gorm:"column:phone;type:varchar(16)" json:"phone"`
	PostalCodeFormat   string `gorm:"column:postal_code_format;type:varchar(64)" json:"postal_code_format"`
	PostalCodeRegex    string `gorm:"column:postal_code_regex;type:varchar(256)" json:"postal_code_regex"`
	Languages          string `gorm:"column:languages;type:varchar(200)" json:"languages"`
	GeonameId          *int   `gorm:"column:geonameId;type:int(11)" json:"geoname_id"`
	Neighbours         string `gorm:"column:neighbours;type:varchar(64)" json:"neighbours"`
	EquivalentFipsCode string `gorm:"column:equivalent_fips_code;type:varchar(3)" json:"equivalent_fips_code"`
	CurrencySymbol     string `gorm:"column:currency_symbol;type:varchar(3)" json:"currency_symbol"`
	Status             string `gorm:"column:status;type:varchar(15);default:'inactive'" json:"status"`
}

func (Country) TableName() string {
	return "countryinfo"
}
