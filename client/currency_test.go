package client

import (
	"context"
	"github.com/Rhymond/go-money"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"

	"go.tradeforge.dev/oxr/model"
)

func TestCurrencyClient_GetCurrencies(t *testing.T) {
	type args struct {
		ctx    context.Context
		params *model.GetCurrenciesParams
		opts   []model.RequestOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "pass",
			args: args{
				ctx:    context.Background(),
				params: &model.GetCurrenciesParams{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc := NewClient(OpenExchangeRatesAPIURL, "") // This endpoint is free.
			got, err := cc.GetCurrencies(tt.args.ctx, tt.args.params, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCurrencies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Greater(t, len(got), 0)
		})
	}
}

func TestCurrencyClient_GetLatestRates(t *testing.T) {
	testAPIKey := os.Getenv("TEST_OXR_API_KEY")

	type args struct {
		ctx    context.Context
		params *model.GetLatestRatesParams
		opts   []model.RequestOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "pass",
			args: args{
				ctx:    context.Background(),
				params: &model.GetLatestRatesParams{},
			},
		},
		{
			name: "pass/with-base-currency",
			args: args{
				ctx: context.Background(),
				params: &model.GetLatestRatesParams{
					Base: money.CZK,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc := NewClient(OpenExchangeRatesAPIURL, testAPIKey)
			got, err := cc.GetLatestRates(tt.args.ctx, tt.args.params, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLatestRates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Greater(t, len(got.Rates), 0)
			if tt.args.params.Base != "" {
				assert.Equal(t, tt.args.params.Base, got.Base.Code)
			}
		})
	}
}

func TestCurrencyClient_ConvertCurrency(t *testing.T) {
	testAPIKey := os.Getenv("TEST_OXR_API_KEY")

	type args struct {
		ctx    context.Context
		params *model.ConvertCurrencyParams
		opts   []model.RequestOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "pass",
			args: args{
				ctx: context.Background(),
				params: &model.ConvertCurrencyParams{
					From:  money.CZK,
					To:    money.USD,
					Value: decimal.NewFromInt(1),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc := NewClient(OpenExchangeRatesAPIURL, testAPIKey)
			got, err := cc.ConvertCurrency(tt.args.ctx, tt.args.params, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertCurrency() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Greater(t, got.Response, 0)
		})
	}
}
