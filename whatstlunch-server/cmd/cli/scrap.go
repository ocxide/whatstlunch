package cli

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"io"

	"github.com/PuerkitoBio/goquery"
)

type ScrappedMeal struct {
	Title        string
	Introduction string
	Comensales   uint
	Duration     string
	FoodType     string
	Ingredients  []string
	Preparation  []PreparationStep
}

type PreparationStep struct {
	Order       uint64
	Description string
}

func getById(recipeLink string) (*ScrappedMeal, error) {
	fmt.Printf("Getting meal %s\n", recipeLink)

	rawResponse, err := http.Get(recipeLink)
	if err != nil {
		fmt.Printf("Error making recipe request to url %s: %v", recipeLink, err)
		return nil, err
	}

	defer rawResponse.Body.Close()

	rawHtml, err := io.ReadAll(rawResponse.Body)
	if err != nil {
		fmt.Println("Error reading meal content", err)
		return nil, err
	}

	htmlDoc, err := goquery.NewDocumentFromReader(bytes.NewReader(rawHtml))
	if err != nil {
		fmt.Println("Error parsing meal content", err)
		return nil, err
	}

	title := htmlDoc.Find("h1.titulo.titulo--articulo").Text()
	if !strings.HasPrefix(title, "Receta de") {
		return nil, errors.New("Recipe title not found")
	}

	intro := ""
	htmlDoc.Find("div.intro p").Each(func(_ int, s *goquery.Selection) {
		intro += s.Text()
	})

	comensales := htmlDoc.Find("span.comensales").Text()
	comensales = strings.Replace(comensales, " comensales", "", -1)
	comensales = strings.Replace(comensales, " comensal", "", -1)
	comensalesCount, err := strconv.Atoi(comensales)

	duration := htmlDoc.Find("span.duracion").Text()
	foodType := htmlDoc.Find("span.para").Text()

	ingredients := make([]string, 0)
	htmlDoc.Find("li.ingrediente label").Each(func(_ int, s *goquery.Selection) {
		ingredients = append(ingredients, strings.TrimSpace(s.Text()))
	})

	preparations := make([]PreparationStep, 0)
	htmlDoc.Find("div.apartado").Each(func(_ int, apartado *goquery.Selection) {
		orderStr := apartado.Find("div.orden").Text()
		order, err := strconv.ParseUint(orderStr, 10, 64)

		// If cannot parse, it is not a step "apartado"
		if err != nil {
			return
		}

		descripcion := apartado.Find("p").Text()

		preparations = append(preparations, PreparationStep{
			Order:       order,
			Description: strings.TrimSpace(descripcion),
		})
	})

	return &ScrappedMeal{
		Title:        title,
		Introduction: intro,
		Comensales:   uint(comensalesCount),
		Duration:     duration,
		FoodType:     foodType,
		Ingredients:  ingredients,
		Preparation:  preparations,
	}, nil
}

func GetMeals(url string) ([]ScrappedMeal, error) {
	meals := make([]ScrappedMeal, 0)

	rawResponse, err := http.Get(url)
	if err != nil {
		return meals, err
	}

	defer rawResponse.Body.Close()

	rawHtml, err := io.ReadAll(rawResponse.Body)
	if err != nil {
		fmt.Println("Error reading meal content", err)
		return meals, err
	}

	htmlDoc, err := goquery.NewDocumentFromReader(bytes.NewReader(rawHtml))
	if err != nil {
		fmt.Println("Error parsing meal content", err)
		return meals, err
	}

	wg := new(sync.WaitGroup)

	htmlDoc.Find("a.titulo.titulo--resultado").Each(func(_ int, s *goquery.Selection) {
		recipeLink, _ := s.Attr("href")
		wg.Add(1)

		go func() {
			defer wg.Done()
			meal, err := getById(recipeLink)
			if err == nil {
				meals = append(meals, *meal)
			}
		}()

	})

	wg.Wait()

	return meals, nil
}
