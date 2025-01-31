package main

import (
	"errors"
	"fmt"
	"regexp"
	"sync"
	"time"
)

type RecepitStore struct {
	mu    sync.RWMutex
	store map[string]Receipt
}

func (s *RecepitStore) StoreReceipt(id string, r Receipt) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.store[id] = r
}

func (s *RecepitStore) GetReceipt(id string) (Receipt, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	r, ok := s.store[id]
	if !ok {
		return Receipt{}, fmt.Errorf("receipt not found")
	}

	return r, nil
}

type Service struct {
	store *RecepitStore
}

type ReqProcessReceipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []struct {
		ShortDescription string `json:"shortDescription"`
		Price            string `json:"price"`
	} `json:"items"`
	Total string `json:"total"`
}

var (
	reReceiptRetailer             = regexp.MustCompile(`^[\\w\\s\\-&]+$`)
	reReceiptTotal                = regexp.MustCompile(`^\\d+\\.\\d{2}$`)
	reReceiptItemShortDescription = regexp.MustCompile(`^[\\w\\s\\-]+$`)
	reReceiptItemPrice            = regexp.MustCompile(`^\\d+\\.\\d{2}$`)
)

func (r ReqProcessReceipt) IsValid() error {
	var err error

	if r.Retailer == "" {
		err = errors.Join(err, errors.New("receipt: retailer cannot be empty"))
	}

	if !reReceiptRetailer.MatchString(r.Retailer) {
		err = errors.Join(err, errors.New("receipt: retailer must be alphanumeric"))
	}

	if r.PurchaseDate == "" {
		err = errors.Join(err, errors.New("receipt: purchase date cannot be empty"))
	}

	_, derr := time.Parse(time.DateOnly, r.PurchaseDate)
	if derr != nil {
		err = errors.Join(err, fmt.Errorf("receipt: purchase date must be in the format of YYYY-MM-DD - %w", derr))
	}

	if r.PurchaseTime == "" {
		err = errors.Join(err, errors.New("receipt: purchase time cannot be empty"))
	}

	_, terr := time.Parse("15:04", r.PurchaseTime)
	if terr != nil {
		err = errors.Join(err, fmt.Errorf("receipt: purchase time must be in the format of HH:MM - %w", terr))
	}

	if len(r.Items) == 0 {
		err = errors.Join(err, errors.New("receipt: items cannot be empty"))
	}

	if r.Total == "" {
		err = errors.Join(err, errors.New("receipt: total cannot be empty"))
	}

	if !reReceiptTotal.MatchString(r.Total) {
		err = errors.Join(err, errors.New("receipt: total must be in the format of 0.00"))
	}

	for _, item := range r.Items {
		if item.ShortDescription == "" {
			err = errors.Join(err, errors.New("item: short description cannot be empty"))
		}

		if item.Price == "" {
			err = errors.Join(err, errors.New("item: price cannot be empty"))
		}

		if !reReceiptItemShortDescription.MatchString(item.ShortDescription) {
			err = errors.Join(err, errors.New("item: short description must be alphanumeric"))
		}

		if !reReceiptItemPrice.MatchString(item.Price) {
			err = errors.Join(err, errors.New("item: price must be in the format of 0.00"))
		}

	}
	return err
}

// Converts validated request to a Receipt struct
func ConvertReqToReceiptTwo(r ReqProcessReceipt) ReceiptTwo {
	time.RFC1123
}

type RespProcessReceipt struct {
	Id string `json:"id"`
}

func (s Service) ProcessReceipt(req ReqProcessReceipt) (*RespProcessReceipt, error) {
	if err := req.IsValid(); err != nil {
		return nil, fmt.Errorf("error validating request: %w", err)
	}

	s.store.StoreReceipt(receipt.Id, receipt)

	return &RespProcessReceipt{Id: receipt.Id}, nil
}

type ReqGetPoints struct {
	Id string `json:"id"`
}

func (r ReqGetPoints) IsValid() error {
	if r.Id == "" {
		return fmt.Errorf("id cannot be empty")
	}
	return nil
}

func (s Service) GetPoints(req ReqGetPoints) (int64, error) {
	if err := req.IsValid(); err != nil {
		return 0, fmt.Errorf("error validating request: %w", err)
	}
	r, err := s.store.GetReceipt(req.Id)
	if err != nil {
		return 0, fmt.Errorf("error getting recepit: %w", err)
	}

	return r.Points, nil
}
