package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/heytechie/Students-api-go/internal/storage"
	"github.com/heytechie/Students-api-go/internal/types"
	"github.com/heytechie/Students-api-go/internal/utils/response"
)

func New(storage storage.Storage) http.HandlerFunc {
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

		id, err := storage.CreateStudent(student.Name, student.Email, student.Age)
		if err != nil {
			response.WriteJSONResponse(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		slog.Info("User created Successfully", slog.String("userId: ", fmt.Sprint(id)))

		response.WriteJSONResponse(w, http.StatusCreated, map[string]string{
			"success": "ok",
			"id":      fmt.Sprintf("%d", id),
		})
	}
}

func GetStudentById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the student ID from the URL
		id := r.PathValue("id")
		slog.Info("getting student by id", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJSONResponse(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid student ID")))
			return
		}
		student, err := storage.GetStudentById(intId)
		if err != nil {
			response.WriteJSONResponse(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJSONResponse(w, http.StatusOK, student)
	}
}

func GetAllStudents(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Getting all students")
		students, err := storage.GetAllStudents()
		if err != nil {
			response.WriteJSONResponse(w, http.StatusInternalServerError, err)
			return
		}

		response.WriteJSONResponse(w, http.StatusOK, students)
	}
}
