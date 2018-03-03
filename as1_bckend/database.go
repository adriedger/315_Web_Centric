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
	q := `INSERT INTO class VALUES(:class_name, :creator_key)`
	_, err := db.NamedExec(q, class)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) JoinClass(enrollment Enrollment) error {
	q := `INSERT INTO enrollment(username, class_name) VALUES(:username, :class_name)`
	_, err := db.NamedExec(q, enrollment)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) GetClasses() ([]Class, error) {
	classes := []Class{}
	q := `SELECT * FROM class`
	err := db.Select(&classes, q)
	if err != nil {
		return []Class{}, err
	}
	return classes, nil
}

func (db *Database) AddQuestion(question Question) error {
	//key attempt
	classes := []Class{}
	q := `SELECT * FROM class WHERE class_name = $1 AND creator_key = $2`
	err := db.Select(&classes, q, question.ClassName, question.KeyAttempt)
	if err != nil {
		return err
	}
	if len(classes) < 1 {
		return fmt.Errorf("database -> keys do not match")
	}
	//add question
	q = `INSERT INTO questions(question, answer, class_name) VALUES(:question, :answer, :class_name)`
	_, err = db.NamedExec(q, question)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) AddResponse(response Response) error {
	//add response
	q := `INSERT INTO responses(response, enroll_id, question_id, class_name, username, question)
		SELECT :response, e.enroll_id, q.question_id, :class_name, :username, :question
		FROM enrollment AS e
		INNER JOIN questions AS q ON e.class_name = :class_name AND e.username = :username
		AND q.class_name = :class_name AND q.question = :question`
	res, err := db.NamedExec(q, response)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("database -> given parametes do not exist")
	}
	return nil
}

func (db *Database) ModifyResponse(response Response) error {
	//modify response
	q := `UPDATE responses SET response = :response
		 WHERE class_name = :class_name AND username = :username AND question = :question`
	res, err := db.NamedExec(q, response)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("database -> question with those params dosent exist")
	}
	return nil
}

func (db *Database) DeleteQuestion(question Question) error {
	//check if key_attempt matches creator key
	classes := []Class{}
	q := `SELECT * FROM class WHERE class_name = $1 AND creator_key = $2`
	err := db.Select(&classes, q, question.ClassName, question.KeyAttempt)
	if err != nil {
		return err
	}
	if len(classes) < 1 {
		return fmt.Errorf("database -> keys do not match")
	}
	//delete question
	q = `DELETE FROM questions WHERE question = :question AND class_name = :class_name`
	_, err = db.NamedExec(q, question)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) GetQuestions(class_name string) ([]Question, error) {
	questions := []Question{}
	q := `SELECT question, class_name, answer FROM questions WHERE class_name = $1`
	err := db.Select(&questions, q, class_name)
	if err != nil {
		return []Question{}, err
	}
	/*
		if len(questions) < 1 {
			return []Question{}, fmt.Errorf("database -> class does not have any questions")
		}
	*/
	return questions, nil
}

func (db *Database) GetResponses(question Question) ([]Response, error) {
	//check key attempt
	classes := []Class{}
	q := `SELECT * FROM class WHERE class_name = $1 AND creator_key = $2`
	err := db.Select(&classes, q, question.ClassName, question.KeyAttempt)
	if err != nil {
		return []Response{}, err
	}
	if len(classes) < 1 {
		return []Response{}, fmt.Errorf("database -> keys do not match")
	}
	//get responses
	responses := []Response{}
	q = `SELECT class_name, username, question, response FROM responses WHERE question = $1 and class_name = $2`
	err = db.Select(&responses, q, question.Question, question.ClassName)
	if err != nil {
		return []Response{}, err
	}
	/*
		if len(responses) < 1 {
			return []Response{}, fmt.Errorf("database -> question does not have any responses")
		}
	*/
	return responses, nil
}
