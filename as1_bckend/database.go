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
}

func yo() {
	fmt.Println("yo")
}
