package main

import (
	"errors"
	"regexp"
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
	if i.ShortDescription == "" {
		return errors.New("Item short description cannot be empty")
	}

	if i.Price == "" {
		return errors.New("Item price cannot be empty")
	}

	if !reItemShortDescription.MatchString(i.ShortDescription) {
		return errors.New("Item short description must be alphanumeric")
	}

	if !reItemPrice.MatchString(i.Price) {
		return errors.New("Item price must be in the format of 0.00")
	}

	return nil
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
	if r.Retailer == "" {
		return errors.New("Receipt retailer cannot be empty")
	}

	if r.PurchaseDate == "" {
		return errors.New("Receipt purchase date cannot be empty")
	}

	if r.PurchaseTime == "" {
		return errors.New("Receipt purchase time cannot be empty")
	}

	if len(r.Items) == 0 {
		return errors.New("Receipt items cannot be empty")
	}

	if r.Total == "" {
		return errors.New("Receipt total cannot be empty")
	}
	return nil
}
