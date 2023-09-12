package responses

type Data struct {
	ExchangeRate string `json:"exchange_rate"`
}

type CurrencyResponse struct {
	Data []Data `json:"data"`
}
