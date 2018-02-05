package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var Connect string = "dbname=as1_bckend user=postgres host=localhost port=5432 sslmode=disable"

type Database struct {
	*sqlx.DB
}

// OpenDatabase attempts to open the database specified by DataSource
// and return a handle to it
func OpenDatabase() (*Database, error) {
	db := Database{}
	var err error

	db.DB, err = sqlx.Open("postgres", Connect)
	if err != nil {
		return nil, fmt.Errorf("Open (%v): %v", Connect, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Ping: %v", err)
	}

	return &db, nil
}

func (db *Database) AddClass(class Class) error {
	//	fmt.Printf("%+v\n", class)
	q := `INSERT INTO class VALUES(:class_id, :class_name, :creator_key)`
	_, err := db.NamedExec(q, class)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) GetClass(class_id string) (Class, error) {
	classes := []Class{}
	q := `SELECT * FROM class WHERE class_id = $1`
	err := db.Select(&classes, q, class_id)
	if err != nil {
		return Class{}, err
	}
	if len(classes) < 1 {
		return Class{}, fmt.Errorf("database -> class does not exist")
	}
	return classes[0], nil
	//add functionality to return all students enrolled in the class
}

func (db *Database) JoinClass(enrollment Enrollment) error {
	//check if class exists, get class id
	classes := []Class{}
	q := `SELECT * FROM class WHERE class_id = $1`
	err := db.Select(&classes, q, enrollment.ClassID)
	if err != nil {
		return err
	}
	if len(classes) < 1 {
		return fmt.Errorf("database -> class does not exist")
	}
	//add class id and username to enrollment
	q = `INSERT INTO enrollment VALUES(:enroll_id, :username, :class_id)`
	_, err = db.NamedExec(q, enrollment)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) AddQuestion(question Question) error {
	q := `INSERT INTO questions VALUES(:question, :class_id, :answer)`
	_, err := db.NamedExec(q, question)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) AddResponse(response Response) error {
	//check if response already exists for given enroll id
	responses := []Response{}
	q := `SELECT * FROM responses WHERE question = $1 AND enroll_id = $2`
	err := db.Select(&responses, q, response.Question, response.EnrollID)
	if err != nil {
		return err
	}
	if len(responses) > 0 {
		return fmt.Errorf("database -> response already exists for give question and enroll id")
	}
	//add response
	q = `INSERT INTO responses VALUES(:enroll_id, :response, :question)`
	_, err = db.NamedExec(q, response)
	if err != nil {
		return err
	}
	return nil
}
