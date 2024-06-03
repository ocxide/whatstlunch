package infer

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"net/http"
	"strings"

	"github.com/disintegration/imaging"

	"github.com/ocxide/whatstlunch/cmd/config"
)

type CompletionResponse struct {
	Response string `json:"response"`
}

func resizeImage(file io.ReadCloser) (*image.NRGBA, error) {
	defer file.Close()

	decoded, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	dist := imaging.Resize(decoded, 512, 0, imaging.Lanczos)

	return dist, nil
}

type InferIngredients struct {
	Config config.AiConfig
}

func (handler *InferIngredients) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	file, header, err := req.FormFile("image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()
	format, err := imaging.JPEG, nil
	switch header.Header.Get("Content-Type") {
	case "image/jpeg":
		format = imaging.JPEG
	case "image/png":
		format = imaging.PNG
	default:
		err = imaging.ErrUnsupportedFormat
	}

	if err != nil {
		fmt.Fprint(w, "Not supported image format, ", header.Header.Get("Content-Type"), "\n")
		fmt.Fprint(w, "Supported image formats: jpeg, png\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resizedImg, err := resizeImage(file)
	rezizedImgBytes := new(bytes.Buffer)

	err = imaging.Encode(rezizedImgBytes, resizedImg, format)
	if err != nil {
		fmt.Fprint(w, "Error encoding image")
		fmt.Printf("Error encoding image\n")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	base64Img := base64.StdEncoding.EncodeToString(rezizedImgBytes.Bytes())
	content := strings.NewReader(`{
			"model": "` + handler.Config.Model + `",
			"prompt": "Crea una lista de los ingredientes (frutas, verduras, especias, carnes, etc) de lo que puedas ver en la imagen. solo nombres, simple, todo en español. Use dashes to list them.",
			"stream": false,
			"images": ["` + base64Img + `"]
	}`)

	response, err := http.Post(handler.Config.ApiUrl+"/generate", "application/json", content)
	if err != nil {
		fmt.Fprint(w, "Error infering ingredients - error contection LLM\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	completion := CompletionResponse{}

	err = decoder.Decode(&completion)
	if err != nil {
		fmt.Fprint(w, "Error decoding response")
		return
	}

	items := strings.Split(completion.Response, "-")
	if len(items) > 0 && strings.TrimSpace(items[0]) == "" {
		items = items[1:]
	}

	ingredients := make([]string, 0)
	for _, ingredient := range items {
		ingredient = strings.TrimSpace(ingredient)
		ingredient = strings.ToLower(ingredient)

		ingredients = append(ingredients, ingredient)
	}

	err = json.NewEncoder(w).Encode(ingredients)
	if err != nil {
		fmt.Fprint(w, "Error encoding response")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}
