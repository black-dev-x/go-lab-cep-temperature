package main

import (
	"encoding/json"
	"net/http"

	"github.com/black-dev-x/go-lab-cep-temperature/cep"
	"github.com/black-dev-x/go-lab-cep-temperature/config"
	"github.com/black-dev-x/go-lab-cep-temperature/temperature"
	"github.com/black-dev-x/go-lab-cep-temperature/weather"
)

type CepInput struct {
	Cep string `json:"cep"`
}

func main() {
	config.Load()
	http.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		var input CepInput
		json.NewDecoder(r.Body).Decode(&input)
		cepResponse, err := cep.Get(input.Cep)
		if err != nil {
			if err.Error() == cep.NotFound {
				http.Error(w, cep.NotFound, 404)
			} else if err.Error() == cep.Invalid {
				http.Error(w, cep.Invalid, 422)
			} else {
				http.Error(w, err.Error(), 500)
			}
			return
		}
		weather, err := weather.Get(cepResponse.Localidade)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		temp := temperature.New(cepResponse.Localidade, weather.Current.TempC)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(temp)
	})
	http.ListenAndServe(":4000", nil)
}
