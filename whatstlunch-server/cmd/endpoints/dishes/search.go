package dishes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

func parseRequirement(requireStr string) (float64, bool, error) {
	if requireStr == "" {
		return 0.0, false, nil
	}

	require := 0.0

	requireI, err := strconv.Atoi(requireStr)
	isCount := err == nil

	if isCount {
		require = float64(requireI)
		return require, true, nil
	}

	require, err = strconv.ParseFloat(requireStr, 64)

	if err != nil {
		return 0.0, false, err
	}

	if require > 1.0 {
		return float64(int(require)), true, nil
	}

	return require, false, nil
}

func Search(res http.ResponseWriter, req *http.Request) {
	ingredients, has := req.URL.Query()["ingredient"]

	if !has {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	for i, ingredient := range ingredients {
		ingredients[i] = strings.ToLower(strings.TrimSpace(ingredient))
	}

	requireStr := req.URL.Query().Get("require")
	require, isCount, err := parseRequirement(requireStr)

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	args := make([]any, 0, len(ingredients))
	filter := "WHERE ingredients LIKE "

	for i, ingredient := range ingredients {
		filter = filter + "?"
		args = append(args, "%"+ingredient+"%")

		if i < len(ingredients)-1 {
			filter += " OR ingredients LIKE "
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
			(
				SELECT group_concat(idescription, ';') FROM
				(SELECT description as idescription FROM ingredients i WHERE i.meal_id = m.id)
			) as ingredients,
			(
				SELECT GROUP_CONCAT(pdescription, ';') FROM
				(SELECT description as pdescription FROM preparations p WHERE p.meal_id = m.id)
			) as preparation
		FROM meals m `+filter,
		args...,
	)

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	parsedDishes := make([]DishFound, len(dishes))
	for i, dish := range dishes {
		parsedDishes[i] = DishFound{
			Title:        dish.Title,
			Introduction: dish.Introduction,
			Duration:     dish.Duration,
			FoodType:     dish.FoodType,
			Ingredients:  strings.Split(dish.Ingredients, ";"),
			Preparation:  strings.Split(dish.Preparation, ";"),
		}
	}

	minIngredients := 0
	if isCount {
		minIngredients = int(require)
	} else {
		minIngredients = int(float64(len(ingredients)) * require)
	}

	parsedDishes = filterDishesByRequire(parsedDishes, ingredients, minIngredients)

	res.WriteHeader(http.StatusOK)
	res.Header().Add("Content-Type", "application/json")

	json.NewEncoder(res).Encode(parsedDishes)
}

func filterDishesByRequire(dishes []DishFound, searchIngredients []string, require int) []DishFound {
	filtered := make([]DishFound, 0, len(dishes))
	for _, dish := range dishes {

		if !matches(dish.Ingredients, searchIngredients, require) {
			continue
		}

		filtered = append(filtered, dish)
	}

	return filtered
}

func matches(ingredients []string, searchIngredients []string, mustMatch int) bool {
	matched := 0

	for _, search := range searchIngredients {
		anyMatch(&matched, ingredients, search)

		if matched >= mustMatch {
			return true
		}
	}

	return false
}

func anyMatch(matched *int, ingredients []string, search string) {
	for _, ingredient := range ingredients {
		if strings.Contains(ingredient, search) {
			*matched += 1
			return
		}
	}
}
