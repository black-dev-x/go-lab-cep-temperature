package main

import (
	"encoding/json"
	"net/http"

	"github.com/black-dev-x/go-lab-cep-temperature/config"
)

type CepInput struct {
	Cep string `json:"cep"`
}

func main() {
	config.Load()
	http.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		var input CepInput
		json.NewDecoder(r.Body).Decode(&input)
		cep := input.Cep
		length := len(cep)
		if length != 8 {
			http.Error(w, "invalid zipcode", 422)
			return
		}
	})
	http.ListenAndServe(":3000", nil)
}
