package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cdleo/api-go-financial-accounting/internal/service"
	"github.com/gorilla/mux"
)

type budgetSummaryHandler struct {
	retrieveUC service.BudgetSummaryRetrieve
}

func NewBudgetSummaryHandler(retrieve service.BudgetSummaryRetrieve) Handler {
	return &budgetSummaryHandler{
		retrieveUC: retrieve,
	}
}

func (h *budgetSummaryHandler) RegisterEndpoints(router *mux.Router) {
	router.HandleFunc("/budget-summary", h.retrieveBudgetSummary).Methods("GET")
}

func (h *budgetSummaryHandler) retrieveBudgetSummary(w http.ResponseWriter, r *http.Request) {

	urlMonth := r.URL.Query().Get("month")
	urlYear := r.URL.Query().Get("year")

	if urlMonth == "" || urlYear == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	month, err := strconv.Atoi(urlMonth)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	year, err := strconv.Atoi(urlYear)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value, err := h.retrieveUC.GetBudgetSummaryByDate(month, year)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "unable to retrieve: %v", err)
		return
	}

	if value == nil {
		fmt.Fprintf(w, "unable to retrieve budget for month [%d] year [%d]", month, year)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(value); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "unable to encode response: %v", err)
		return
	}

}
