package client

import (
	"context"
	"fmt"
	"net/http"

	"go.tradeforge.dev/oxr/internal/client"
	"go.tradeforge.dev/oxr/model"

	"github.com/Rhymond/go-money"
)

const (
	GetCurrenciesPath   = "/currencies.json"
	GetLatestRatesPath  = "/latest.json"
	ConvertCurrencyPath = "/convert/:value/:from/:to"
)

type CurrencyClient struct {
	*client.Client
}

func (cc *CurrencyClient) GetCurrencies(ctx context.Context, params *model.GetCurrenciesParams, opts ...model.RequestOption) (model.GetCurrenciesResponse, error) {
	aux := map[string]string{}
	err := cc.Call(ctx, http.MethodGet, GetCurrenciesPath, params, &aux)
	res := map[money.Currency]string{}

	for currencyCode, v := range aux {
		currency := money.GetCurrency(currencyCode)
		if currency == nil {
			continue // skip invalid or unknown currencies
		}
		res[*money.GetCurrency(currencyCode)] = v
	}
	return res, err
}

func (cc *CurrencyClient) GetLatestRates(ctx context.Context, params *model.GetLatestRatesParams, opts ...model.RequestOption) (*model.GetLatestRatesResponse, error) {
	aux := model.GetLatestRatesResponseMarshallable{}
	err := cc.Call(ctx, http.MethodGet, GetLatestRatesPath, params, &aux)
	if err != nil {
		return nil, err
	}

	baseCurrency := money.GetCurrency(aux.Base)
	if baseCurrency == nil {
		return nil, fmt.Errorf("invalid base currency: %s", aux.Base)
	}
	rates := map[money.Currency]float64{}
	for currencyCode, rate := range aux.Rates {
		currency := money.GetCurrency(currencyCode)
		if currency == nil {
			continue // skip invalid or unknown currencies
		}
		rates[*currency] = rate
	}

	return &model.GetLatestRatesResponse{
		Base:       *baseCurrency,
		Rates:      rates,
		License:    aux.License,
		Disclaimer: aux.Disclaimer,
		Timestamp:  aux.Timestamp,
	}, nil
}

func (cc *CurrencyClient) ConvertCurrency(ctx context.Context, params *model.ConvertCurrencyParams, opts ...model.RequestOption) (*model.ConvertCurrencyResponse, error) {
	res := &model.ConvertCurrencyResponse{}
	err := cc.Call(ctx, http.MethodGet, ConvertCurrencyPath, params, res)
	return res, err
}
