package util

var (
	SupportedCurrencies = []string{"USD", "EUR", "CAD"}
)

func IsSupportedCurrency(currency string) bool {
	for _, cur := range SupportedCurrencies {
		if cur == currency {
			return true
		}
	}
	return false
}
