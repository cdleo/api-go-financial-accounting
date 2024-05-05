package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cdleo/api-go-financial-accounting/internal/entity"
	"github.com/gorilla/mux"
)

type transferenceHandler struct {
	transferUC entity.MakeTransference
}

func NewTransferenceHandler(transfer entity.MakeTransference) Handler {
	return &transferenceHandler{
		transferUC: transfer,
	}
}

func (h *transferenceHandler) RegisterEndpoints(router *mux.Router) {
	router.HandleFunc("/transfer", h.makeTransference).Methods("POST")
}

func (h *transferenceHandler) makeTransference(w http.ResponseWriter, r *http.Request) {

	var request entity.Transfer

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "input invalid: %v", err)
		return
	}

	err := h.transferUC.MakeTransference(context.TODO(), request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "unable to create: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
