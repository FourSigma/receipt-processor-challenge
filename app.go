package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func NewAPI() API {
	return API{
		svc: NewService(),
	}
}

type API struct {
	svc *Service
}

func (a API) Run() error {
	http.HandleFunc("POST /receipts/process", a.ProcessReceipt)
	http.HandleFunc("GET /receipts/{id}/points", a.GetReceipt)

	return http.ListenAndServe(":8080", nil)
}

func (a API) ProcessReceipt(rw http.ResponseWriter, r *http.Request) {
	body := ReqProcessReceipt{}

	if err := DecodeJSON(r, &body); err != nil {
		return
	}

	resp, err := a.svc.ProcessReceipt(body)
	if err != nil {
		EncodeJSONError(rw, err)
		return
	}

	EncodeJSON(rw, resp, http.StatusOK)

	return
}

func (a API) GetReceipt(rw http.ResponseWriter, r *http.Request) {
	req := ReqGetPoints{
		Id: r.PathValue("id"),
	}

	resp, err := a.svc.GetPoints(req)
	if err != nil {
		log.Println(err)
		return
	}

	EncodeJSON(rw, resp, http.StatusOK)

	return
}

func DecodeJSON(r *http.Request, val any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(val)
}

func EncodeJSON(rw http.ResponseWriter, val any, code int) {
	rw.WriteHeader(code)
	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(val); err != nil {
		log.Println(err)
	}
}

func EncodeJSONError(rw http.ResponseWriter, err error) {
	code := http.StatusInternalServerError
	message := "internal server error"

	switch {
	case errors.Is(err, ErrInvalidInput):
		code = http.StatusBadRequest
		message = ErrInvalidInput.Error()

	case errors.Is(err, ErrNotFound):
		code = http.StatusNotFound
		message = ErrNotFound.Error()
	}

	http.Error(rw, message, code)
}
