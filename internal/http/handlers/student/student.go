package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"github/mahfujulsagor/student_api/internal/db"
	"github/mahfujulsagor/student_api/internal/logger"
	"github/mahfujulsagor/student_api/internal/types"
	"github/mahfujulsagor/student_api/internal/utils/response"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func New(db db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info.Println("Create student handler called")

		var student types.Student

		//? Decode JSON body
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			logger.Error.Println("Empty body", err)
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			logger.Error.Println("Error decoding student:", err)
			return
		}
		defer r.Body.Close()

		//? Request validation
		if err := validator.New().Struct(student); err != nil {
			//? Typecast the error to ValidationErrors
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(err.(validator.ValidationErrors)))
			logger.Error.Println("Validation error:", err)
			return
		}

		//* Create the student in DB
		id, err := db.CreateStudent(student.Name, student.Email, student.Age)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			logger.Error.Println("Error creating student:", err)
			return
		}

		logger.Info.Println("Student created with ID:", id)

		//? Send response
		response.WriteJson(w, http.StatusCreated, map[string]string{
			"success": "OK",
			"message": fmt.Sprintf("Student created with ID %d", id),
		})
	}
}
