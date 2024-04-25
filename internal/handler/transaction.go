package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cdleo/api-go-financial-accounting/internal/entity"
	"github.com/cdleo/api-go-financial-accounting/internal/service"
	"github.com/gorilla/mux"
)

type transactionHandler struct {
	createUC service.MakeTransaction
}

func NewTransactionHandler(create service.MakeTransaction) Handler {
	return &transactionHandler{
		createUC: create,
	}
}

func (h *transactionHandler) RegisterEndpoints(router *mux.Router) {
	router.HandleFunc("/account-record", h.addTrxToAccount).Methods("POST")
}

func (h *transactionHandler) addTrxToAccount(w http.ResponseWriter, r *http.Request) {

	var request entity.Transaction

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "input invalid: %v", err)
		return
	}

	err := h.createUC.MakeTransaction(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "unable to create: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
