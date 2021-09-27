package util

//Constants for all supported Currency
const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
	NP  = "NP"
)

//IsSupportedCurrency returns turns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD, NP:
		return true
	}
	return false
}
