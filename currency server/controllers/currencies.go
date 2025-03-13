package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"odyssey/m/v2/repositories"
	"odyssey/m/v2/usecases/currencies"
	"time"

	"gorm.io/gorm"
)

type CurrenciesController struct {
	usecase *currencies.CurrencyUseCase
}

func NewCurrencyController(db *gorm.DB) *CurrenciesController {
	currencyRepo, err := repositories.NewCurrencyRepo(db)
	if err != nil {
		log.Fatal(err)
	}
	ratesRepo, err := repositories.NewRateRepo(db)
	if err != nil {
		log.Fatal(err)
	}

	currencyUseCase, err := currencies.NewCurrencyCase(currencyRepo, ratesRepo)
	if err != nil {
		log.Fatal(err)
	}

	return &CurrenciesController{usecase: currencyUseCase}
}

func (c *CurrenciesController) GetRates(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query().Get("startDate")
	if len(start) == 0 {
		c.errorResponse(w, "startDate is missing", http.StatusUnprocessableEntity)
	}

	end := r.URL.Query().Get("endDate")
	if len(end) == 0 {
		c.errorResponse(w, "endDate is missing", http.StatusUnprocessableEntity)
	}

	startDate, err := time.Parse("2006-01-02", start)
	if err != nil {
		c.errorResponse(w, "Bad response"+err.Error(), http.StatusUnprocessableEntity)
	}
	endDate, err := time.Parse("2006-01-02", end)
	if err != nil {
		c.errorResponse(w, "Bad response"+err.Error(), http.StatusUnprocessableEntity)
	}

	result, err := c.usecase.GetRates(startDate, endDate)
	if err != nil {
		c.errorResponse(w, "Bad response"+err.Error(), http.StatusUnprocessableEntity)
	}
	jsonResp, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

func (c *CurrenciesController) Get(w http.ResponseWriter, r *http.Request) {
	data, err := c.usecase.Get()
	if err != nil {
		c.errorResponse(w, "Bad request "+err.Error(), http.StatusUnprocessableEntity)
		return
	}

	jsonResp, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResp)
	w.WriteHeader(http.StatusOK)
}

func (c *CurrenciesController) errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}
