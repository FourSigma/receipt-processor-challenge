package main

import (
	"encoding/json"
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
	http.HandleFunc("GET /receipts/:id/points", a.GetReceipt)

	return nil
}

func (a API) DecodeJSON(r *http.Request, val any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(val)
}

func (a API) EncodeJSON(rw http.ResponseWriter, val any) error {
	rw.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(val)
}

func (a API) ProcessReceipt(rw http.ResponseWriter, r *http.Request) {
	body := ReqProcessReceipt{}

	if err := a.DecodeJSON(r, &body); err != nil {
		return
	}

	resp, err := a.svc.ProcessReceipt(body)
	if err != nil {
		return
	}

	if err := a.EncodeJSON(rw, resp); err != nil {
		return
	}

	return
}

func (a API) GetReceipt(rw http.ResponseWriter, r *http.Request) {
	req := ReqGetPoints{
		Id: r.PathValue("id"),
	}

	resp, err := a.svc.GetPoints(req)
	if err != nil {
		return
	}

	if err := a.EncodeJSON(rw, resp); err != nil {
		return
	}

	return
}
