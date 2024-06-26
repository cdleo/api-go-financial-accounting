package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cdleo/api-go-financial-accounting/internal/entity"
	"github.com/gorilla/mux"
)

type accountHandler struct {
	createUC   entity.AccountCreate
	retrieveUC entity.AccountRetrieve
	updateUC   entity.AccountUpdate
}

func NewAccountHandler(create entity.AccountCreate, retrieve entity.AccountRetrieve, update entity.AccountUpdate) Handler {
	return &accountHandler{
		createUC:   create,
		retrieveUC: retrieve,
		updateUC:   update,
	}
}

func (h *accountHandler) RegisterEndpoints(router *mux.Router) {
	router.HandleFunc("/account", h.createAccount).Methods("POST")
	router.HandleFunc("/account", h.retrieveAccount).Methods("GET")
	router.HandleFunc("/accounts", h.retrieveAccounts).Methods("GET")
	router.HandleFunc("/account", h.updateAccount).Methods("PUT")
}

func (h *accountHandler) createAccount(w http.ResponseWriter, r *http.Request) {

	var request entity.Account

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "input invalid: %v", err)
		return
	}

	err := h.createUC.CreateAccount(&request)
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

func (h *accountHandler) retrieveAccount(w http.ResponseWriter, r *http.Request) {

	accountId := r.URL.Query().Get("id")

	if accountId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value, err := h.retrieveUC.GetAccountByID(accountId)
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

func (h *accountHandler) retrieveAccounts(w http.ResponseWriter, r *http.Request) {

	accounts, err := h.retrieveUC.GetAccounts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "unable to retrieve: %v", err)
		return
	}

	if len(accounts) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(accounts); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "unable to encode response: %v", err)
		return
	}
}

func (h *accountHandler) updateAccount(w http.ResponseWriter, r *http.Request) {

	var request entity.Account

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "input invalid: %v", err)
		return
	}

	err := h.updateUC.UpdateAccount(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "unable to update: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)

}
