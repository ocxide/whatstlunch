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

func InferIngredients(w http.ResponseWriter, req *http.Request) {
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
			"model": "llava:7b",
			"prompt": "Crea una lista de los ingredientes (frutas, verduras, especias, carnes, etc) de lo que puedas ver en la imagen. Incluye solo nombres, todo en español.\n- ",
			"stream": false,
			"images": ["` + base64Img + `"]
	}`)

	response, err := http.Post("http://127.0.0.1:11434/api/generate", "application/json", content)
	if err != nil {
		fmt.Fprint(w, "Error infering ingredients - error contection LLM", err)
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

	ingredients := make([]string, 0)
	for _, ingredient := range strings.Split(completion.Response, "\n- ") {
		ingredients = append(ingredients, ingredient)
	}

	err = json.NewEncoder(w).Encode(ingredients)
	if err != nil {
		fmt.Fprint(w, "Error encoding response")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
