package dishes

import (
	"encoding/json"
	"net/http"

	"github.com/ocxide/whatstlunch/cmd/database"
)

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

	dishes := []DishFound{}
	db.Select(dishes, "SELECT title, introduction, duration, food_type, ingredients, preparation FROM dishes WHERE "+filter, args...)

	if len(dishes) == 0 {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Header().Add("Content-Type", "application/json")

	json.NewEncoder(res).Encode(dishes)
}
