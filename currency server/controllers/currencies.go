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

type GetRatesBody struct {
	StartDate  string   `json:"startDate"`
	EndDate    string   `json:"endDate"`
	Currencies []string `json:"currencies"`
}

func (c *CurrenciesController) GetRates(w http.ResponseWriter, r *http.Request) {
	var body GetRatesBody
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&body)
	if err != nil {
		c.errorResponse(w, "Bad request "+err.Error(), http.StatusUnprocessableEntity)
		return
	}

	startDate, err := time.Parse("2006-01-02", body.StartDate)
	if err != nil {
		c.errorResponse(w, "Bad response"+err.Error(), http.StatusUnprocessableEntity)
	}
	endDate, err := time.Parse("2006-01-02", body.EndDate)
	if err != nil {
		c.errorResponse(w, "Bad response"+err.Error(), http.StatusUnprocessableEntity)
		return
	}

	result, err := c.usecase.GetRates(startDate, endDate, body.Currencies)
	if err != nil {
		c.errorResponse(w, "Bad response"+err.Error(), http.StatusUnprocessableEntity)
		return
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

func (c *CurrenciesController) Sync(w http.ResponseWriter, r *http.Request) {
	var body GetRatesBody
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&body)
	if err != nil {
		c.errorResponse(w, "Bad request "+err.Error(), http.StatusUnprocessableEntity)
		return
	}

	startDate, err := time.Parse("2006-01-02", body.StartDate)
	if err != nil {
		c.errorResponse(w, "Bad response"+err.Error(), http.StatusUnprocessableEntity)
	}
	endDate, err := time.Parse("2006-01-02", body.EndDate)
	if err != nil {
		c.errorResponse(w, "Bad response"+err.Error(), http.StatusUnprocessableEntity)
		return
	}
	err = c.usecase.Sync(startDate, endDate, body.Currencies)
	if err != nil {
		c.errorResponse(w, "Bad response"+err.Error(), http.StatusUnprocessableEntity)
		return
	}

	result, err := c.usecase.GetRates(startDate, endDate, body.Currencies)
	if err != nil {
		c.errorResponse(w, "Bad response"+err.Error(), http.StatusUnprocessableEntity)
		return
	}
	jsonResp, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

func (c *CurrenciesController) errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}
