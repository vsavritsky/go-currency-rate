package model

type Rate struct {
	CurrencyCode string  `json:"currencyCode"`
	Provider     string  `json:"provider"`
	Value        float64 `json:"value"`
}
