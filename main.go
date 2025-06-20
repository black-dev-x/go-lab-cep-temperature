package main

import (
	"encoding/json"
	"net/http"

	"github.com/black-dev-x/go-lab-cloud-run/cep"
	"github.com/black-dev-x/go-lab-cloud-run/config"
	"github.com/black-dev-x/go-lab-cloud-run/temperature"
	"github.com/black-dev-x/go-lab-cloud-run/weather"
)

func main() {
	config.Load()
	http.HandleFunc("GET /{cep}", func(w http.ResponseWriter, r *http.Request) {
		cepInput := r.PathValue("cep")
		cepResponse, err := cep.Get(cepInput)
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
		temp := temperature.New(weather.Current.TempC)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(temp)
	})
	http.ListenAndServe(":8080", nil)
}
