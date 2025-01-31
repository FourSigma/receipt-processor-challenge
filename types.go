package main

import (
	"errors"
	"regexp"
	"time"
)

var (
	ErrInvalidInput = errors.New("Please verify input.")
	ErrNotFound     = errors.New("No receipt found for that ID.")
)

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

var (
	reItemShortDescription = regexp.MustCompile(`^[\\w\\s\\-]+$`)
	reItemPrice            = regexp.MustCompile(`^\\d+\\.\\d{2}$`)
)

func (i Item) IsValid() error {
	var err error

	if i.ShortDescription == "" {
		err = errors.Join(err, errors.New("item: short description cannot be empty"))
	}

	if i.Price == "" {
		err = errors.Join(err, errors.New("item: price cannot be empty"))
	}

	if !reItemShortDescription.MatchString(i.ShortDescription) {
		err = errors.Join(err, errors.New("item: short description must be alphanumeric"))
	}

	if !reItemPrice.MatchString(i.Price) {
		err = errors.Join(err, errors.New("item: price must be in the format of 0.00"))
	}

	return err
}

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

var (
	reReceiptRetailer = regexp.MustCompile(`^[\\w\\s\\-&]+$`)
	reReceiptTotal    = regexp.MustCompile(`^\\d+\\.\\d{2}$`)
)

func (r Receipt) IsValid() error {
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

	_, terr := time.Parse("2006-01-02", r.PurchaseDate)
	if terr != nil {
		err = errors.Join(err, errors.New("receipt: purchase date must be in the format of YYYY-MM-DD"))
	}

	if r.PurchaseTime == "" {
		err = errors.Join(err, errors.New("receipt: purchase time cannot be empty"))
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
		if ierr := item.IsValid(); err != nil {
			err = errors.Join(err, ierr)
		}
	}

	return err
}
