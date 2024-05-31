package main

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

type Meal struct {
	Title       string
	Intro       string
	Comensales  uint
	Duration    string
	FoodType    string
	Preparation []Preparation
}

type Preparation struct {
	Order       string
	Description string
}

func getById(recipeLink string) (*Meal, error) {
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
		ingredients = append(ingredients, s.Text())
	})

	preparations := make([]Preparation, 0)
	htmlDoc.Find("div.apartado").Each(func(_ int, apartado *goquery.Selection) {
		orden := apartado.Find("div.orden").Text()
		descripcion := apartado.Find("p").Text()

		preparations = append(preparations, Preparation{
			Order:       orden,
			Description: descripcion,
		})
	})

	return &Meal{
		Title:       title,
		Intro:       intro,
		Comensales:  uint(comensalesCount),
		Duration:    duration,
		FoodType:    foodType,
		Preparation: preparations,
	}, nil
}

func GetMeals(url string) ([]Meal, error) {
	meals := make([]Meal, 0)

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
