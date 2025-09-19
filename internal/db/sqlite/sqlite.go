package sqlite

import (
	"database/sql"
	"fmt"
	"github/mahfujulsagor/student_api/internal/config"
	"github/mahfujulsagor/student_api/internal/types"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	DB *sql.DB
}

// ? New initializes the SQLite database and returns a SQLite instance
func New(cfg *config.Config) (*SQLite, error) {
	db, err := sql.Open("sqlite3", cfg.DBPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT,
		age INTEGER
	)`)
	if err != nil {
		return nil, err
	}

	return &SQLite{
		DB: db,
	}, nil
}

// ? CreateStudent inserts a new student into the database and returns the inserted ID
func (s *SQLite) CreateStudent(name string, email string, age int) (int64, error) {
	//? Prepare statement to prevent SQL injection
	stmt, err := s.DB.Prepare("INSERT INTO students(name, email, age) VALUES(?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	//? Execute the statement
	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	//? Get the last inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// ? GetStudentByID retrieves a student by ID from the database
func (s *SQLite) GetStudentByID(id int64) (types.Student, error) {
	//? Prepare statement to prevent SQL injection
	stmt, err := s.DB.Prepare("SELECT id, name, email, age FROM students WHERE id = ? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	//? Execute the query
	var student types.Student
	err = stmt.QueryRow(id).Scan(&student.ID, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("student with ID %d not found", id)
		}
		return types.Student{}, fmt.Errorf("query error: %w", err)
	}

	return student, nil
}
