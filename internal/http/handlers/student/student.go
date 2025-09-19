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
	"strconv"

	"github.com/go-playground/validator/v10"
)

// ? New handles the creation of a new student
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

// ? GetByID handles the retrieval of a student by ID
func GetByID(db db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info.Println("Get by ID handler called")

		//? Extract ID from URL
		id := r.PathValue("id")
		if id == "" {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("missing student ID in URL")))
			logger.Error.Println("Missing student ID in URL")
			return
		}

		//? Convert ID to int64
		studentID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid student ID")))
			logger.Error.Println("Error parsing student ID:", err)
			return
		}

		//? Retrieve student from DB
		student, err := db.GetStudentByID(studentID)
		if err != nil {
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(err))
			logger.Error.Println("Error retrieving student:", err)
			return
		}

		//? Send response
		response.WriteJson(w, http.StatusOK, student)
	}
}

// ? GetList handles the retrieval of a list of students with pagination
// ? Max limit is 50 to prevent abuse
func GetList(db db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info.Println("Get list handler called")

		//? Extract pagination parameters from query
		query := r.URL.Query()
		limitStr := query.Get("limit")
		offsetStr := query.Get("offset")

		//? Set default values if not provided
		if limitStr == "" {
			limitStr = "10"
		}
		if offsetStr == "" {
			offsetStr = "0"
		}

		//? Convert limit and offset to integers
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid limit value")))
			logger.Error.Println("Invalid limit value:", err)
			return
		}
		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid offset value")))
			logger.Error.Println("Invalid offset value:", err)
			return
		}

		//? Enforce a hard upper bound on limit
		const maxLimit int = 50
		if limit > maxLimit {
			limit = maxLimit
		}

		//? Retrieve students from DB
		students, err := db.GetStudents(limit, offset)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			logger.Error.Println("Error retrieving students:", err)
			return
		}

		//? Handle empty results gracefully
		if len(students) == 0 {
			response.WriteJson(w, http.StatusOK, []types.Student{})
			return
		}

		//? Send response
		response.WriteJson(w, http.StatusOK, students)
	}
}

// ? UpdateByID handles the updating of a student by ID
func UpdateByID(db db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info.Println("Update by ID handler called")
		//? Extract ID from URL
		id := r.PathValue("id")
		if id == "" {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("missing student ID in URL")))
			logger.Error.Println("Missing student ID in URL")
			return
		}

		//? Convert ID to int64
		studentID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid student ID")))
			logger.Error.Println("Error parsing student ID:", err)
			return
		}

		//? Check if student exists
		_, err = db.GetStudentByID(studentID)
		if err != nil {
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(err))
			logger.Error.Println("Error retrieving student:", err)
			return
		}

		var student types.Student

		//? Decode JSON body
		err = json.NewDecoder(r.Body).Decode(&student)
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

		//* Update the student in DB
		err = db.UpdateStudentByID(studentID, student.Name, student.Email, student.Age)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			logger.Error.Println("Error updating student with ID:", studentID)
			return
		}

		logger.Info.Println("Student updated with ID:", studentID)
		//? Send response
		response.WriteJson(w, http.StatusOK, map[string]string{
			"success": "ok",
			"message": fmt.Sprintf("Student updated with ID %d", studentID),
		})
	}
}

// ? DeleteByID handles the deletion of a student by ID
func DeleteByID(db db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info.Println("Delete by ID handler called")
		//? Extract ID from URL
		id := r.PathValue("id")
		if id == "" {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("missing student ID in URL")))
			logger.Error.Println("Missing student ID in URL")
			return
		}

		//? Convert ID to int64
		studentID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid student ID")))
			logger.Error.Println("Error parsing student ID:", err)
			return
		}

		//? Check if student exists
		_, err = db.GetStudentByID(studentID)
		if err != nil {
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(err))
			logger.Error.Println("Error retrieving student:", err)
			return
		}

		//* Delete the student from DB
		err = db.DeleteStudentByID(studentID)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			logger.Error.Println("Error deleting student with ID:", studentID)
			return
		}

		logger.Info.Println("Student deleted with ID:", studentID)
		//? Send response
		response.WriteJson(w, http.StatusOK, map[string]string{
			"success": "ok",
			"message": fmt.Sprintf("Student deleted with ID %d", studentID),
		})
	}
}
