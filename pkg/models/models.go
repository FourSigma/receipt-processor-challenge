package models

import (
	"errors"
	"time"
)

var (
	ErrInvalidInput = errors.New("The receipt is invalid.")
	ErrNotFound     = errors.New("No receipt found for that ID.")
)

type Item struct {
	ShortDescription string
	Price            float64
}

type Receipt struct {
	Id          string
	Retailer    string
	Items       []Item
	PurchasedAt time.Time
	Total       float64
	Points      int64
}
