package controllers

import (
	"crypto/rsa"
	"encoding/json"
	"log"
	"net/http"
	"odyssey/m/v2/usecases/verify"
)

type VerifyController struct {
	usecase *verify.VerifyUseCase
}

func NewVerifyController(key *rsa.PublicKey) *VerifyController {
	verifyUseCase, err := verify.NewVerifyUseCase(key)
	if err != nil {
		log.Fatal(err)
	}

	return &VerifyController{usecase: verifyUseCase}
}

func (c *VerifyController) Verify(w http.ResponseWriter, r *http.Request) {
	data := r.URL.Query().Get("data")
	if len(data) == 0 {
		c.errorResponse(w, "data is missing", http.StatusUnprocessableEntity)
	}
	signature := r.URL.Query().Get("signature")
	if len(data) == 0 {
		c.errorResponse(w, "signature is missing", http.StatusUnprocessableEntity)
	}

	err := c.usecase.Check(data, []byte(signature))
	if err != nil {
		c.errorResponse(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (c *VerifyController) errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}
