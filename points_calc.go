package main

import (
	"math"
	"strings"
	"time"
	"unicode"
)

type RuleHandlerFn func(Receipt) int64

func CalculatePoints(r Receipt, fns ...RuleHandlerFn) int64 {
	var points int64

	for _, fn := range fns {
		points = points + fn(r)
	}

	return points
}

// One point for every alphanumeric character in the retailer name.
func RuleAlphanumeric(r Receipt) int64 {
	var points int64

	for _, v := range r.Retailer {
		if !unicode.IsLetter(v) && !unicode.IsNumber(v) {
			continue
		}
		points = points + 1
	}

	return points
}

// 50 points if the total is a round dollar amount with no cents
func RuleRoundDollar(r Receipt) int64 {
	if r.Total == float64(int64(r.Total)) {
		return 50
	}
	return 0
}

// 25 points if the total is a multiple of 0.25
func RuleMultipleOfQuarter(r Receipt) int64 {
	if math.Mod(r.Total*4, 1) == 0 {
		return 25
	}
	return 0
}

// 5 points for every two item on the receipt
func RuleItemPair(r Receipt) int64 {
	return int64(len(r.Items)/2) * 5
}

// If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer.
// The result is the number of points earned.
func RuleItemDescription(r Receipt) int64 {
	var points int64

	for _, item := range r.Items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			points = points + int64(math.Ceil(item.Price*0.2))
		}
	}

	return points
}

// 6 points if the day in the purchase date is odd
func RuleOddDay(r Receipt) int64 {
	if r.PurchasedAt.Day()%2 != 0 {
		return 6
	}
	return 0
}

// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
func RuleTimeOfPurchase(r Receipt) int64 {
	t := r.PurchasedAt

	twoPM := time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		14, // hour (2 PM = 14:00)
		0,
		0,
		0,
		t.Location(),
	)

	fourPM := time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		16, // hour (4 PM = 16:00)
		0,
		0,
		0,
		t.Location(),
	)

	if t.After(twoPM) && t.Before(fourPM) {
		return 10
	}
	return 0
}
