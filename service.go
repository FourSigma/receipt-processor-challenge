package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
)

type RecepitStore struct {
	mu    sync.RWMutex
	store map[string]Receipt
}

func (s *RecepitStore) StoreReceipt(r Receipt) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.store[r.Id] = r
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

func NewService() *Service {
	return &Service{
		store: &RecepitStore{store: map[string]Receipt{}},
	}
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
	reReceiptRetailer             = regexp.MustCompile("^[\\w\\s\\-&]+$")
	reReceiptTotal                = regexp.MustCompile("^\\d+\\.\\d{2}$")
	reReceiptItemShortDescription = regexp.MustCompile("^[\\w\\s\\-]+$")
	reReceiptItemPrice            = regexp.MustCompile("^\\d+\\.\\d{2}$")
)

var (
	ErrRetailerEmpty               = errors.New("retailer cannot be empty")
	ErrRetailerInvalid             = errors.New("retailer must be alphanumeric")
	ErrPurchaseDateEmpty           = errors.New("purchase date cannot be empty")
	ErrPurchaseDateInvalid         = errors.New("purchase date must be in the format of YYYY-MM-DD")
	ErrPurchaseTimeEmpty           = errors.New("purchase time cannot be empty")
	ErrPurchaseTimeInvalid         = errors.New("purchase time must be in the format of HH:MM")
	ErrItemsEmpty                  = errors.New("items cannot be empty")
	ErrTotalEmpty                  = errors.New("total cannot be empty")
	ErrTotalInvalid                = errors.New("total must be in the format of 0.00")
	ErrItemShortDescriptionEmpty   = errors.New("item short description cannot be empty")
	ErrItemShortDescriptionInvalid = errors.New("item short description must be alphanumeric")
	ErrItemPriceEmpty              = errors.New("item price cannot be empty")
	ErrItemPriceInvalid            = errors.New("item price must be in the format of 0.00")
)

func (r ReqProcessReceipt) IsValid() error {
	var err error

	if r.Retailer == "" {
		err = errors.Join(err, ErrRetailerEmpty)
	}

	if !reReceiptRetailer.MatchString(r.Retailer) {
		err = errors.Join(err, ErrRetailerInvalid)
	}

	if r.PurchaseDate == "" {
		err = errors.Join(err, ErrPurchaseDateEmpty)
	}

	_, derr := time.Parse(time.DateOnly, r.PurchaseDate)
	if derr != nil {
		err = errors.Join(err, ErrPurchaseDateInvalid)
	}

	if r.PurchaseTime == "" {
		err = errors.Join(err, ErrPurchaseTimeEmpty)
	}

	_, terr := time.Parse("15:04", r.PurchaseTime)
	if terr != nil {
		err = errors.Join(err, fmt.Errorf("%w: %w", ErrPurchaseTimeInvalid, terr))
	}

	if len(r.Items) == 0 {
		err = errors.Join(err, ErrItemsEmpty)
	}

	if r.Total == "" {
		err = errors.Join(err, ErrTotalEmpty)
	}

	if !reReceiptTotal.MatchString(r.Total) {
		err = errors.Join(err, ErrTotalInvalid)
	}

	for _, item := range r.Items {
		if item.ShortDescription == "" {
			err = errors.Join(err, ErrItemShortDescriptionEmpty)
		}

		if item.Price == "" {
			err = errors.Join(err, ErrItemPriceEmpty)
		}

		if !reReceiptItemShortDescription.MatchString(item.ShortDescription) {
			err = errors.Join(err, ErrItemShortDescriptionInvalid)
		}

		if !reReceiptItemPrice.MatchString(item.Price) {
			err = errors.Join(err, ErrItemPriceInvalid)
		}

	}
	return err
}

// Converts validated request to a Receipt struct
func ConvertReqToReceiptTwo(req ReqProcessReceipt) (Receipt, error) {
	receipt := Receipt{
		Id:       uuid.NewString(),
		Retailer: req.Retailer,
		Points:   0,
	}

	total, err := strconv.ParseFloat(req.Total, 64)
	if err != nil {
		return Receipt{}, fmt.Errorf("error parsing total: %w", err)
	}

	receipt.Total = total

	const timeFormat = "2006-01-02 15:04"
	receipt.PurchasedAt, err = time.Parse(timeFormat, req.PurchaseDate+" "+req.PurchaseTime)
	if err != nil {
		return Receipt{}, fmt.Errorf("error parsing purchase date and time: %w", err)
	}

	for _, item := range req.Items {
		price, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			return Receipt{}, fmt.Errorf("error parsing item price: %w", err)
		}

		receipt.Items = append(receipt.Items, Item{
			ShortDescription: item.ShortDescription,
			Price:            price,
		})
	}

	return receipt, nil
}

type RespProcessReceipt struct {
	Id string `json:"id"`
}

func (s Service) ProcessReceipt(req ReqProcessReceipt) (*RespProcessReceipt, error) {
	if err := req.IsValid(); err != nil {
		return nil, fmt.Errorf("invalid request - %w %w", ErrInvalidInput, err)
	}

	receipt, err := ConvertReqToReceiptTwo(req)
	if err != nil {
		return nil, fmt.Errorf("error converting request to receipt: %w", err)
	}

	rules := []RuleHandlerFn{
		RuleAlphanumeric,
		RuleRoundDollar,
		RuleMultipleOfQuarter,
		RuleItemPair,
		RuleItemDescription,
		RuleOddDay,
		RuleTimeOfPurchase,
	}

	receipt.Points = CalculatePoints(receipt, rules...)

	s.store.StoreReceipt(receipt)

	return &RespProcessReceipt{Id: receipt.Id}, nil
}

var (
	ErrIdEmpty   = errors.New("id cannot be empty")
	ErrIdInvalid = errors.New("id must must be a non-whitespace character")

	reReceiptId = regexp.MustCompile("^\\S+$")
)

type ReqGetPoints struct {
	Id string `json:"id"`
}

func (r ReqGetPoints) IsValid() error {
	var err error
	if r.Id == "" {
		err = errors.Join(err, ErrIdEmpty)
	}

	if !reReceiptId.MatchString(r.Id) {
		err = errors.Join(err, ErrIdInvalid)
	}

	return err
}

type RespGetPoints struct {
	Points int64 `json:"points"`
}

func (s Service) GetPoints(req ReqGetPoints) (*RespGetPoints, error) {
	if err := req.IsValid(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidInput, err)
	}

	r, err := s.store.GetReceipt(req.Id)
	if err != nil {
		return nil, fmt.Errorf("error getting recepit: %w", err)
	}

	resp := &RespGetPoints{
		Points: r.Points,
	}

	return resp, nil
}
