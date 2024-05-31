package main

import (
	"fmt"
	"sync"

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

func Load(url string) {
	db, err := sqlx.Connect("sqlite3", "meals.db")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	meals, err := GetMeals(url)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Saving %d meals into 'meals.db'\n", len(meals))

	tx, err := db.Begin()

	if err != nil {
		panic(err)
	}

	// Drop tables if they exist, only aplicable for this specific `load` command
	// Other commands may want to just append data rather than replacing the existing data
	tx.Exec("DELETE FROM preparations")
	tx.Exec("DROP TABLE IF EXISTS preparations")

	tx.Exec("DELETE FROM ingredients")
	tx.Exec("DROP TABLE IF EXISTS ingredients")

	tx.Exec("DELETE FROM meals")
	tx.Exec("DROP TABLE IF EXISTS meals")
	// 

	tx.Exec(schema)

	err = tx.Commit()
	if err != nil {
		fmt.Println("Could not setup database", err)
		return
	}

	wg := new(sync.WaitGroup)

	for _, meal := range meals {
		wg.Add(1)

		go func(meal ScrappedMeal) {
			defer wg.Done()
			SaveMeal(meal, db)
		}(meal)
	}

	wg.Wait()
}

func SaveMeal(meal ScrappedMeal, db *sqlx.DB) {
	fmt.Printf("Saving meal %s\n", meal.Title)

	result, err := db.Exec(
		"INSERT INTO meals (title, introduction, comensales, duration, food_type) VALUES (?, ?, ?, ?, ?) RETURNING id",
		meal.Title, meal.Introduction, meal.Comensales, meal.Duration, meal.FoodType,
	)

	if err != nil {
		return
	}

	meal_id, err := result.LastInsertId()

	if err != nil {
		fmt.Printf("Error saving meal %s: %v\n", meal.Title, err)
		return
	}

	for _, ingredient := range meal.Ingredients {
		_, err := db.Exec(
			"INSERT INTO ingredients (description, meal_id) VALUES (?, ?)",
			ingredient, meal_id,
		)

		if err != nil {
			fmt.Printf("Error saving ingredient %s: %v\n", ingredient, err)
			continue
		}
	}

	for _, preparation := range meal.Preparation {
		_, err := db.Exec(
			"INSERT INTO preparations (`order`, description, meal_id) VALUES (?, ?, ?)",
			preparation.Order, preparation.Description, meal_id,
		)

		if err != nil {
			fmt.Printf(
				"Error saving preparation \"%s\" with order %d: %v\n",
				truncate(preparation.Description, 20), preparation.Order, err,
			)
			continue
		}
	}
}

func truncate(str string, max int) string {
	if len(str) <= max {
		return str
	}
	return str[:max] + "..."
}
