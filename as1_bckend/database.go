package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
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
	fmt.Println("here1")
	fmt.Printf("%+v\n", class)
	fmt.Println("here2")
	q := `INSERT INTO class VALUES(:class_id, :class_name, :creator_key)`
	_, err := db.NamedExec(q, class)
	if err != nil {
		return err
		os.Exit(1)
	}
	return nil
}

func yo() {
	fmt.Println("yo")
}
