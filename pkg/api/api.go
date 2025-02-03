package api

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FourSigma/receipt-processor-challenge/pkg/models"
	"github.com/FourSigma/receipt-processor-challenge/pkg/service"
)

func New() API {
	return API{
		svc: service.NewService(),
	}
}

type API struct {
	svc *service.Service
}

func (a API) Run() {
	// Routes
	mux := http.NewServeMux()
	mux.HandleFunc("POST /receipts/process", a.ProcessReceipt)
	mux.HandleFunc("GET /receipts/{id}/points", a.GetReceipt)

	// Server setup and shutdown
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {

		log.Println("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server on :8080 - %s", err)
		}

	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Starting to server shutdown...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	log.Println("Server gracefully stopped...")

}

func (a API) ProcessReceipt(rw http.ResponseWriter, r *http.Request) {
	body := service.ReqProcessReceipt{}

	if err := DecodeJSON(r, &body); err != nil {
		return
	}

	resp, err := a.svc.ProcessReceipt(r.Context(), body)
	if err != nil {
		EncodeJSONError(rw, err)
		return
	}

	EncodeJSON(rw, resp, http.StatusOK)
}

func (a API) GetReceipt(rw http.ResponseWriter, r *http.Request) {
	req := service.ReqGetPoints{
		Id: r.PathValue("id"),
	}

	resp, err := a.svc.GetPoints(r.Context(), req)
	if err != nil {
		EncodeJSONError(rw, err)
		return
	}

	EncodeJSON(rw, resp, http.StatusOK)
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
	case errors.Is(err, models.ErrInvalidInput):
		code = http.StatusBadRequest
		message = models.ErrInvalidInput.Error()

	case errors.Is(err, models.ErrNotFound):
		code = http.StatusNotFound
		message = models.ErrNotFound.Error()
	}

	http.Error(rw, message, code)
}
