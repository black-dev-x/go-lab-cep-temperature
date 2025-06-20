package cep

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type ErrorResponse struct {
	Erro bool `json:"erro"`
}

const NotFound = "can not find zipcode"
const Invalid = "invalid zipcode"

func Get(cep string) (*CEP, error) {

	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	httpResp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode == 404 {
		return nil, fmt.Errorf(NotFound)
	}
	if httpResp.StatusCode == 400 {
		return nil, fmt.Errorf(Invalid)
	}
	if httpResp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code from zipcode service: %d", httpResp.StatusCode)
	}

	var check map[string]interface{}
	json.NewDecoder(httpResp.Body).Decode(&check)

	if check["erro"] != nil {
		return nil, fmt.Errorf(NotFound)
	}

	data, _ := json.Marshal(check)

	var response CEP
	err = json.Unmarshal(data, &response)

	return &response, err
}
