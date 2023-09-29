package rate

import (
	"api-crypto-project/internal/http-server/handlers/mail"
	"api-crypto-project/internal/lib/api/response"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log"
	"net/http"
	"strings"
)

const op = "handlers.rate.request.New"

func New() http.HandlerFunc {
	return Rate
}
func Rate(w http.ResponseWriter, r *http.Request) {

	resp, err := RequestProcessing(w, r)
	if err != nil {
		return
	}

	err = ResponseProcessing(w, r, resp)
	if err != nil {
		return
	}
}

func RequestProcessing(w http.ResponseWriter, r *http.Request) (*http.Response, error) {
	var resp *http.Response

	req, err := RequestRead(w, r)
	if err != nil {
		return resp, err
	}

	resp, err = RequestSend(w, r, req)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
func RequestRead(w http.ResponseWriter, r *http.Request) (mail.Request, error) {
	var req mail.Request
	err := render.DecodeJSON(r.Body, &req)
	if errors.Is(err, io.EOF) {
		log.Print(op + ": request body is empty")
		render.JSON(w, r, response.Error("empty request"))
		return req, fmt.Errorf(op + ": empty request")
	}
	if err != nil {
		log.Print(op + ":failed to decode request body")
		render.JSON(w, r, response.Error("failed to decode request"))
		return req, fmt.Errorf(op + ":failed to decode request")
	}
	log.Print("request body decoded")

	if err := validator.New().Struct(req); err != nil {
		log.Print(op + ":invalid request")
		render.JSON(w, r, response.Error("invalid request"))
		return req, fmt.Errorf(op + ":invalid request")
	}
	return req, nil
}
func RequestSend(w http.ResponseWriter, r *http.Request, request mail.Request) (*http.Response, error) {
	resp, err := http.Get(mail.GenerateRequest(request))
	if err != nil {
		log.Print(op + ":Failed to get price")
		render.JSON(w, r, response.Error("Failed to get price"))
		return resp, fmt.Errorf(op + ":failed to get price")
	}
	return resp, nil
}

func ResponseProcessing(w http.ResponseWriter, r *http.Request, resp *http.Response) error {
	respBody, err := ResponseRead(w, r, resp)
	if err != nil {
		return err
	}

	if err = ResponseSend(w, r, respBody); err != nil {
		return err
	}
	return nil
}
func ResponseRead(w http.ResponseWriter, r *http.Request, resp *http.Response) ([]byte, error) {
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print(op + ":failed to read response")
		render.JSON(w, r, response.Error("failed to read response"))
		return respBody, fmt.Errorf(op + ":failed to read response")
	}
	return respBody, nil
}
func ResponseSend(w http.ResponseWriter, r *http.Request, respBody []byte) error {
	crypto, err := UnmarshalCredits(respBody)
	price := strings.Split(crypto.Credits.PRICE, ".")
	crypto.Credits.PRICE = price[0]
	if err != nil {
		return err
	}
	render.JSON(w, r, mail.Response{
		Status:  "OK",
		Message: crypto,
	})
	return nil
}
func UnmarshalCredits(respBody []byte) (mail.Crypto, error) {
	var crypto mail.Crypto
	if err := json.Unmarshal(respBody, &crypto); err != nil {
		return crypto, fmt.Errorf(op + ":failed to create response")
	}
	return crypto, nil
}
