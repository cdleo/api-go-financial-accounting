package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cdleo/api-go-financial-accounting/internal/entity"
	"github.com/cdleo/api-go-financial-accounting/internal/service"
	"github.com/gorilla/mux"
)

type transferenceHandler struct {
	transferUC service.MakeTransference
}

func NewTransferenceHandler(transfer service.MakeTransference) Handler {
	return &transferenceHandler{
		transferUC: transfer,
	}
}

func (h *transferenceHandler) RegisterEndpoints(router *mux.Router) {
	router.HandleFunc("/transfer", h.makeTransference).Methods("POST")
}

func (h *transferenceHandler) makeTransference(w http.ResponseWriter, r *http.Request) {

	var request entity.Transference

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "input invalid: %v", err)
		return
	}

	err := h.transferUC.MakeTransference(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "unable to create: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
