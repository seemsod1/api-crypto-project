package subscribe

import (
	"api-crypto-project/internal/lib/api/response"
	"api-crypto-project/internal/storage"
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"strconv"
)

type Request struct {
	MAIL string `json:"mail" validate:"required,email"`
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type MAILSaver interface {
	SaveMail(mailToSave string) (int64, error)
}

func New(mailSaver MAILSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.mail.subscribe.New"

		req, err := RequestRead(w, r)
		if err != nil {
			return
		}

		id, err := mailSaver.SaveMail(req.MAIL)
		if errors.Is(err, storage.MailExists) {
			log.Print(op + ":mail already exists")
			render.JSON(w, r, response.Error("mail already exists"))
			return

		}
		if err != nil {
			log.Print(op + ":failed to add mail")
			render.JSON(w, r, response.Error(op+":failed to add mail"))
			return
		}

		log.Print("mail added" + "id:" + strconv.FormatInt(id, 10))

		render.JSON(w, r, Response{
			Status: "OK",
		})
	}
}

func RequestRead(w http.ResponseWriter, r *http.Request) (Request, error) {
	var req Request
	const op = "handlers.mail.subscribe.New.RequestRead"
	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		log.Print(op + ":failed to decode request body")
		render.JSON(w, r, response.Error(op+": failed to decode request body"))
		return req, fmt.Errorf(op + ": failed to decode request body")
	}
	log.Print("request body decoded")

	if err := validator.New().Struct(req); err != nil {
		log.Print(op + ":invalid request")
		render.JSON(w, r, response.Error(op+":invalid request"))
		return req, fmt.Errorf(op + ":invalid request")
	}
	return req, nil
}
