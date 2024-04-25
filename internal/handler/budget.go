package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cdleo/api-go-financial-accounting/internal/entity"
	"github.com/gorilla/mux"
)

type budgetHandler struct {
	createUC   entity.BudgetCreate
	retrieveUC entity.BudgetRetrieve
	updateUC   entity.BudgetUpdate
}

func NewBudgetHandler(create entity.BudgetCreate, retrieve entity.BudgetRetrieve, update entity.BudgetUpdate) Handler {
	return &budgetHandler{
		createUC:   create,
		retrieveUC: retrieve,
		updateUC:   update,
	}
}

func (h *budgetHandler) RegisterEndpoints(router *mux.Router) {
	router.HandleFunc("/budget", h.createBudget).Methods("POST")
	router.HandleFunc("/budget", h.retrieveBudgets).Methods("GET")
	router.HandleFunc("/budget/{id}", h.retrieveBudget).Methods("GET")
	router.HandleFunc("/budget/{id}", h.updateBudget).Methods("PUT")
}

func (h *budgetHandler) createBudget(w http.ResponseWriter, r *http.Request) {

	var request entity.Budget

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "input invalid: %v", err)
		return
	}

	err := h.createUC.CreateBudget(&request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "unable to create: %v", err)
		return
	}

	response := map[string]interface{}{
		"id": request.ID,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "unable to encode response: %v", err)
	}
}

func (h *budgetHandler) retrieveBudget(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id := strings.TrimSpace(vars["id"])

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value, err := h.retrieveUC.GetBudgetById(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "unable to retrieve: %v", err)
		return
	}

	if value == nil {
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

func (h *budgetHandler) retrieveBudgets(w http.ResponseWriter, r *http.Request) {

	results, err := h.retrieveUC.GetBudgets()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "unable to retrieve: %v", err)
		return
	}

	if len(results) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	budgets := make([]entity.BudgetInfo, 0)
	for _, budget := range results {
		budgets = append(budgets, entity.BudgetInfo{
			ID:          budget.ID,
			Description: budget.Description,
		})
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(budgets); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "unable to encode response: %v", err)
		return
	}
}

func (h *budgetHandler) updateBudget(w http.ResponseWriter, r *http.Request) {

	var request entity.Budget

	vars := mux.Vars(r)

	id := strings.TrimSpace(vars["id"])

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "input invalid: %v", err)
		return
	}

	request.ID = id

	err := h.updateUC.UpdateBudget(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "unable to create: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
