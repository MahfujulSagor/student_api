package db

import "github/mahfujulsagor/student_api/internal/types"

type DB interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentByID(id int64) (types.Student, error)
	GetStudents(limit int, offset int) ([]types.Student, error)
	UpdateStudentByID(id int64, name string, email string, age int) error
	DeleteStudentByID(id int64) error
}
