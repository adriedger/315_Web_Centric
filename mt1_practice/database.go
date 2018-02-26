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
