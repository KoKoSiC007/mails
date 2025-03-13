package controllers

import (
	"crypto/rsa"
	"encoding/json"
	"log"
	"net/http"
	"odyssey/m/v2/usecases/verify"
	"strconv"
)

type VerifyController struct {
	usecase *verify.VerifyUseCase
}

func NewVerifyController(pubKey *rsa.PublicKey, privKey *rsa.PrivateKey) *VerifyController {
	verifyUseCase, err := verify.NewVerifyUseCase(pubKey, privKey)
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
		return
	}

	err := c.usecase.Check(data, []byte(signature))
	if err != nil {
		c.errorResponse(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (c *VerifyController) GetPublicKey(w http.ResponseWriter, r *http.Request) {
	pemData, err := c.usecase.GetKey()
	if err != nil {
		c.errorResponse(w, "Bad request "+err.Error(), http.StatusUnprocessableEntity)
		return
	}
	w.Header().Set("Content-Type", "application/x-x509-ca-cert")
	w.Header().Set("Content-Disposition", "attachment; filename=certificate.pem")
	w.Header().Set("Content-Length", strconv.Itoa(len(*pemData)))

	// Отправляем данные
	w.Write(*pemData)
}

func (c *VerifyController) GetTestMessage(w http.ResponseWriter, r *http.Request) {
	result, err := c.usecase.TestMessage()
	if err != nil {
		c.errorResponse(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	jsonResp, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func (c *VerifyController) errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}
