package main

import (
	"errors"
	"testing"
)

func TestService(t *testing.T) {
	service := NewService()

	t.Run("ProcessRecepit: happy input validation path", func(t *testing.T) {
		req := ReqProcessReceipt{
			Retailer:     "Target",
			PurchaseDate: "2021-01-02",
			PurchaseTime: "11:00",
			Total:        "100.23",
			Items: []struct {
				ShortDescription string `json:"shortDescription"`
				Price            string `json:"price"`
			}{
				{ShortDescription: "Pepsi - 12-oz", Price: "1.25"},
			},
		}

		_, err := service.ProcessReceipt(req)
		if err != nil {
			t.Errorf("got %v, want nil", err)
		}
	})

	t.Run("ProcessRecepit: sad input valdiation path", func(t *testing.T) {
		tests := []struct {
			name    string
			inputFn ReqProcessReceipt
			wantErr error
		}{
			{
				name:    "Validation: should return error if retailer is empty",
				inputFn: ReqProcessReceipt{},
				wantErr: ErrRetailerEmpty,
			},
			{
				name:    "Validation: should return error if purchase date is empty",
				inputFn: ReqProcessReceipt{},
				wantErr: ErrPurchaseDateEmpty,
			},
			{
				name:    "Validation: should return error if purchase time is empty",
				inputFn: ReqProcessReceipt{},
				wantErr: ErrPurchaseTimeEmpty,
			},
			{
				name:    "Validation: should return error if len(items) is 0",
				inputFn: ReqProcessReceipt{},
				wantErr: ErrItemsEmpty,
			},
			{
				name:    "Validation: should return error if total is empty",
				inputFn: ReqProcessReceipt{},
				wantErr: ErrTotalEmpty,
			},
			{
				name: "Validation: should return error if item.shortDescription is empty",
				inputFn: ReqProcessReceipt{Items: []struct {
					ShortDescription string `json:"shortDescription"`
					Price            string `json:"price"`
				}{
					{},
				}},
				wantErr: ErrItemShortDescriptionEmpty,
			},
			{
				name: "Validation: should return error if item.Price is empty",
				inputFn: ReqProcessReceipt{Items: []struct {
					ShortDescription string `json:"shortDescription"`
					Price            string `json:"price"`
				}{
					{},
				}},
				wantErr: ErrItemPriceEmpty,
			},
			{
				name:    "Validation: should return error when purchase date is invalid",
				inputFn: ReqProcessReceipt{PurchaseDate: "01-01-2021"},
				wantErr: ErrPurchaseDateInvalid,
			},
			{
				name:    "Validation: should return error when purchase time is invalid",
				inputFn: ReqProcessReceipt{PurchaseDate: "1:00PM"},
				wantErr: ErrPurchaseTimeInvalid,
			},

			{
				name:    "Validation: should return error when retailer is not alphanumeric",
				inputFn: ReqProcessReceipt{Retailer: "Target&$^#/[]"},
				wantErr: ErrRetailerInvalid,
			},
			{
				name:    "Validation: should return error when total is invalid",
				inputFn: ReqProcessReceipt{Total: "343.343"},
				wantErr: ErrTotalInvalid,
			},
			{
				name: "Validation: should return error if item.shortDescription is invalid",
				inputFn: ReqProcessReceipt{Items: []struct {
					ShortDescription string `json:"shortDescription"`
					Price            string `json:"price"`
				}{
					{ShortDescription: "item1!@#"},
				}},
				wantErr: ErrItemShortDescriptionInvalid,
			},
			{
				name: "Validation: should return error if item.Price is invalid",
				inputFn: ReqProcessReceipt{Items: []struct {
					ShortDescription string `json:"shortDescription"`
					Price            string `json:"price"`
				}{
					{Price: "343.343"},
				}},
				wantErr: ErrItemPriceInvalid,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := service.ProcessReceipt(tt.inputFn)
				if !errors.Is(err, ErrInvalidInput) {
					t.Errorf("Error not wrapped with ErrrInvalidInput: %v", err)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("got %v, want %v", err, tt.wantErr)
				}
			})
		}
	})
}
