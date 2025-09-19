package db

import "github/mahfujulsagor/student_api/internal/types"

type DB interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentByID(id int64) (types.Student, error)
}
