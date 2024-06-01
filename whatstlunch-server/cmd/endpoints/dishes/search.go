package dishes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ocxide/whatstlunch/cmd/database"
)

type DishSelect struct {
	Title        string `db:"title"`
	Introduction string `db:"introduction"`
	Duration     string `db:"duration"`
	FoodType     string `db:"food_type"`
	Ingredients  string `db:"ingredients"`
	Preparation  string `db:"preparation"`
}

type DishFound struct {
	Title        string   `json:"title"`
	Introduction string   `json:"introduction"`
	Duration     string   `json:"duration"`
	FoodType     string   `json:"foodType"`
	Ingredients  []string `json:"ingredients"`
	Preparation  []string `json:"preparation"`
}

func Search(res http.ResponseWriter, req *http.Request) {
	ingredients, has := req.URL.Query()["ingredient"]

	if !has {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	args := make([]any, len(ingredients))
	filter := "ingredients LIKE "

	for i, ingredient := range ingredients {
		filter = filter + "?"
		args = append(args, "%"+ingredient+"%")

		if i < len(ingredients)-1 {
			filter += "AND ingredients LIKE "
		}
	}

	dishes := []DishSelect{}
	err = db.Select(
		&dishes,
		`SELECT
			m.title,
			m.introduction,
			m.duration,
			m.food_type,
			GROUP_CONCAT(i.description) as ingredients,
			GROUP_CONCAT(p.description) as preparation
		FROM meals m
		LEFT JOIN ingredients i ON m.id = i.meal_id
		LEFT JOIN preparations p ON m.id = p.meal_id
		GROUP BY m.id`,

		ingredients[0],
	)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	if len(dishes) == 0 {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	parsedDishes := make([]DishFound, len(dishes))
	for i, dish := range dishes {
		parsedDishes[i] = DishFound{
			Title:        dish.Title,
			Introduction: dish.Introduction,
			Duration:     dish.Duration,
			FoodType:     dish.FoodType,
			Ingredients:  strings.Split(dish.Ingredients, ","),
			Preparation:  strings.Split(dish.Preparation, ","),
		}
	}

	res.WriteHeader(http.StatusOK)
	res.Header().Add("Content-Type", "application/json")

	json.NewEncoder(res).Encode(parsedDishes)
}
