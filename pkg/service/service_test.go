package service

import (
	"errors"
	"testing"

	"github.com/FourSigma/receipt-processor-challenge/pkg/models"
)

func TestServiceProcessRecepit(t *testing.T) {
	t.Run("ProcessRecepit: happy path", func(t *testing.T) {
		service := NewService()

		t.Run("ProcessRecepit: Example 1 points calculation", func(t *testing.T) {
			req := ReqProcessReceipt{
				Retailer:     "M&M Corner Market",
				PurchaseDate: "2022-03-20",
				PurchaseTime: "14:33",
				Total:        "9.00",
				Items: []struct {
					ShortDescription string `json:"shortDescription"`
					Price            string `json:"price"`
				}{
					{ShortDescription: "Gatorade", Price: "2.25"},
					{ShortDescription: "Gatorade", Price: "2.25"},
					{ShortDescription: "Gatorade", Price: "2.25"},
					{ShortDescription: "Gatorade", Price: "2.25"},
				},
			}

			resp, err := service.ProcessReceipt(req)
			if err != nil {
				t.Errorf("got %v, want nil", err)
			}

			respPoints, err := service.GetPoints(ReqGetPoints{Id: resp.Id})
			if err != nil {
				t.Errorf("got %v, want nil", err)
			}

			if respPoints.Points != 109 {
				t.Errorf("got %v, want 109", respPoints.Points)
			}
		})

		t.Run("ProcessRecepit: Example 2 points calculation", func(t *testing.T) {
			req := ReqProcessReceipt{
				Retailer:     "Target",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Total:        "35.35",
				Items: []struct {
					ShortDescription string `json:"shortDescription"`
					Price            string `json:"price"`
				}{
					{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
					{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
					{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
					{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
					{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
				},
			}

			resp, err := service.ProcessReceipt(req)
			if err != nil {
				t.Errorf("got %v, want nil", err)
			}

			respPoints, err := service.GetPoints(ReqGetPoints{Id: resp.Id})
			if err != nil {
				t.Errorf("got %v, want nil", err)
			}

			if respPoints.Points != 28 {
				t.Errorf("got %v, want 28", respPoints.Points)
			}
		})
	})

	t.Run("ProcessRecepit: sad input valdiation path", func(t *testing.T) {
		service := NewService()

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
				if !errors.Is(err, models.ErrInvalidInput) {
					t.Errorf("Error not wrapped with ErrrInvalidInput: %v", err)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("got %v, want %v", err, tt.wantErr)
				}
			})
		}
	})
}

func TestServiceGetPoints(t *testing.T) {
	t.Run("GetPoints: happy validation path", func(t *testing.T) {
		service := NewService()

		service.store.StoreReceipt(models.Receipt{
			Id: "c4c34d52-bd98-4b81-80e1-fd4939dbe7fb",
		})

		req := ReqGetPoints{Id: "c4c34d52-bd98-4b81-80e1-fd4939dbe7fb"}

		_, err := service.GetPoints(req)
		if err != nil {
			t.Errorf("got %v, want nil", err)
		}
	})

	t.Run("GetPoints: sad valdiation path", func(t *testing.T) {
		service := NewService()
		tests := []struct {
			name    string
			inputFn ReqGetPoints
			wantErr error
		}{
			{
				name:    "Validation: should return error if Id is empty",
				inputFn: ReqGetPoints{},
				wantErr: ErrIdEmpty,
			},
			{
				name:    "Validation: should return error if Id is invalid",
				inputFn: ReqGetPoints{Id: "#^#$&*^fjerkej dlfjdlfj"},
				wantErr: ErrIdInvalid,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := service.GetPoints(tt.inputFn)
				if !errors.Is(err, models.ErrInvalidInput) {
					t.Errorf("Error not wrapped with ErrrInvalidInput: %v", err)
				}

				if !errors.Is(err, tt.wantErr) {
					t.Errorf("got %v, want %v", err, tt.wantErr)
				}
			})
		}
	})
}
