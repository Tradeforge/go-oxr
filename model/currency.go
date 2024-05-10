package model

import (
	"github.com/Rhymond/go-money"
	"github.com/shopspring/decimal"
)

type GetCurrenciesParams struct {
	ShowAlternative bool `query:"show_alternative,omitempty"`
	ShowInactive    bool `query:"show_inactive,omitempty"`
}

type GetCurrenciesResponse map[money.Currency]string

type GetLatestRatesParams struct {
	Base    string   `query:"base,omitempty"`
	Symbols []string `query:"symbols,omitempty"`
}

type GetLatestRatesResponseMarshallable struct {
	Base       string             `json:"base"`
	Rates      map[string]float64 `json:"rates"`
	License    string             `json:"license"`
	Disclaimer string             `json:"disclaimer"`
	Timestamp  int                `json:"timestamp"`
}

type GetLatestRatesResponse struct {
	Base       money.Currency             `json:"base"`
	Rates      map[money.Currency]float64 `json:"rates"`
	License    string                     `json:"license"`
	Disclaimer string                     `json:"disclaimer"`
	Timestamp  int                        `json:"timestamp"`
}

type ConvertCurrencyParams struct {
	From  string          `path:"from"`
	To    string          `path:"to"`
	Value decimal.Decimal `path:"value"`
}

type ConvertCurrencyResponse struct {
	Meta struct {
		Timestamp int     `json:"timestamp"`
		Rate      float64 `json:"rate"`
	} `json:"meta"`

	Disclaimer string          `json:"disclaimer"`
	License    string          `json:"license"`
	Response   decimal.Decimal `json:"response"`
}
