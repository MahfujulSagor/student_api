package db

type DB interface {
	CreateStudent(name string, email string, age int) (int64, error)
}
