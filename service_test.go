package main

import (
	"errors"
	"testing"
)

func TestService(t *testing.T) {
	service := NewService()

	t.Run("ProcessRecepit", func(t *testing.T) {
		tests := []struct {
			name    string
			inputFn ReqProcessReceipt
			wantErr error
		}{
			{
				name:    "Validation: should return error if retailer is empty",
				inputFn: ReqProcessReceipt{},
				wantErr: ErrRetailerAlphanumeric,
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
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := service.ProcessReceipt(tt.inputFn)
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("got %v, want %v", err, tt.wantErr)
				}
			})
		}
	})
}
