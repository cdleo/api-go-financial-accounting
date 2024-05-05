package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cdleo/api-go-financial-accounting/internal/entity"
	"github.com/gorilla/mux"
	"github.com/swaggest/openapi-go"
	"github.com/swaggest/rest"
	"github.com/swaggest/rest/gorillamux"
	"github.com/swaggest/rest/nethttp"
	"github.com/swaggest/rest/request"
)

type budgetSummaryReq struct {
}

type budgetSummaryHandler struct {
	retrieveUC entity.BudgetSummaryRetrieve
	dec        nethttp.RequestDecoder
}

func NewBudgetSummaryHandler(retrieve entity.BudgetSummaryRetrieve) Handler {
	decoderFactory := request.NewDecoderFactory()
	decoderFactory.ApplyDefaults = true
	decoderFactory.SetDecoderFunc(rest.ParamInPath, gorillamux.PathToURLValues)

	return &budgetSummaryHandler{
		retrieveUC: retrieve,
		dec:        decoderFactory.MakeDecoder(http.MethodGet, budgetSummaryReq{}, nil),
	}
}

func (h *budgetSummaryHandler) RegisterEndpoints(router *mux.Router) {
	router.Handle("/budget-summary", h).Methods("GET")
}

// SetupOpenAPIOperation declares OpenAPI schema for the handler.
func (h *budgetSummaryHandler) SetupOpenAPIOperation(oc openapi.OperationContext) error {
	oc.SetTags("My Tag")
	oc.SetSummary("My Summary")
	oc.SetDescription("This endpoint aggregates request in structured way.")

	oc.AddReqStructure(budgetSummaryReq{})
	oc.AddRespStructure(entity.BudgetSummaryAccount{})
	oc.AddRespStructure(nil, openapi.WithContentType("text/html"), openapi.WithHTTPStatus(http.StatusNoContent))
	oc.AddRespStructure(nil, openapi.WithContentType("text/html"), openapi.WithHTTPStatus(http.StatusBadRequest))
	oc.AddRespStructure(nil, openapi.WithContentType("text/html"), openapi.WithHTTPStatus(http.StatusInternalServerError))
	return nil
}

func (h *budgetSummaryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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

	value, err := h.retrieveUC.GetBudgetSummaryByDate(context.TODO(), month, year)
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
