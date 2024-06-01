package cli

import (
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/ocxide/whatstlunch/cmd/database"
)

func Load(url string) {
	db, err := database.Setup()
	if err != nil {
		fmt.Println("Could not setup database", err)
	}

	defer db.Close()

	meals, err := GetMeals(url)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Saving %d meals into 'meals.db'\n", len(meals))

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
