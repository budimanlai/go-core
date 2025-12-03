package entity

type Country struct {
	IsoAlpha2          string
	IsoAlpha3          string
	IsoNumeric         int
	FipsCode           string
	Name               string
	Capital            string
	AreaInSqKm         int64
	Population         int
	Continent          string
	Tld                string
	CurrencyCode       string
	CurrencyName       string
	Phone              string
	PostalCodeFormat   string
	PostalCodeRegex    string
	Languages          string
	GeonameId          int
	Neighbours         string
	EquivalentFipsCode string
	CurrencySymbol     string
	Status             string
}
