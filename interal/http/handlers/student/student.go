package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/heytechie/Students-api-go/interal/types"
	"github.com/heytechie/Students-api-go/interal/utils/response"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJSONResponse(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err != nil {
			response.WriteJSONResponse(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		//req validation
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJSONResponse(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		response.WriteJSONResponse(w, http.StatusCreated, map[string]string{
			"success": "ok",
		})
	}
}
