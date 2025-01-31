package main

import "fmt"

type Service struct{}

type RequestProcessReceipt struct {
	Receipt
}

type ResponseProcessReceipt struct {
	Id string `json:"id"`
}

func (s Service) ProcessReceipt(req RequestProcessReceipt) (ResponseProcessReceipt, error) {
	if err := req.IsValid(); err != nil {
		return fmt.Errorf("error validating request: %w", err)
	}
	return ResponseProcessReceipt{}, nil
}

func (s Service) GetPoints(id string) (int64, error) {
	return 0, nil
}

type PointsHandler func(Receipt) int64

func CalculatePoints(r Receipt, fns ...PointsHandler) int64 {
	var points int64

	for _, fn := range fns {
		points = points + fn(r)
	}

	return points
}
