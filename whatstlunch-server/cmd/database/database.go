package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var schema = `
PRAGMA foreign_keys;

CREATE TABLE IF NOT EXISTS "preparations" (
	"id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
	"order"	INTEGER NOT NULL,
	"description"	TEXT NOT NULL,
	"meal_id"	INTEGER NOT NULL,
	FOREIGN KEY(meal_id) REFERENCES meals(id)
);

CREATE TABLE IF NOT EXISTS "ingredients" (
	"id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,	
	"description"	TEXT NOT NULL,
	"meal_id"	INTEGER NOT NULL,
	FOREIGN KEY(meal_id) REFERENCES meals(id)
);

CREATE TABLE IF NOT EXISTS "meals" (
	"id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
	"title"	TEXT NOT NULL,
	"introduction"	TEXT,
	"comensales"	INTEGER,
	"duration"	TEXT,
	"food_type"	TEXT
);`

func Connect() (*sqlx.DB, error) {
	return sqlx.Connect("sqlite3", "meals.db")
}

// Creates the db connection and initializes the tables from zero (deletes the previous data)
func Setup() (*sqlx.DB, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}

	tx, err := db.Begin()

	if err != nil {
		return nil, err
	}

	tx.Exec("DELETE FROM preparations")
	tx.Exec("DROP TABLE IF EXISTS preparations")

	tx.Exec("DELETE FROM ingredients")
	tx.Exec("DROP TABLE IF EXISTS ingredients")

	tx.Exec("DELETE FROM meals")
	tx.Exec("DROP TABLE IF EXISTS meals")

	tx.Exec(schema)

	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return db, nil
}
