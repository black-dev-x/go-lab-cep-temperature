package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

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

		json, _ := json.Marshal(input)

		req, _ := http.NewRequest("POST", os.Getenv("PROCESSING_SERVICE_URL"), bytes.NewBuffer(json))
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, "failed to send request", 500)
			return
		}
		defer resp.Body.Close()

		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})
	http.ListenAndServe(":3000", nil)
}
