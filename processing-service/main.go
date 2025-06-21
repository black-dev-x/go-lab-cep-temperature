package main

import (
	"context"
	"encoding/json"
	"net"
	"net/http"

	"github.com/black-dev-x/go-lab-cep-temperature/cep"
	"github.com/black-dev-x/go-lab-cep-temperature/config"
	"github.com/black-dev-x/go-lab-cep-temperature/temperature"
	"github.com/black-dev-x/go-lab-cep-temperature/weather"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
)

type CepInput struct {
	Cep string `json:"cep"`
}

func HTTPHandler() http.Handler {
	mux := http.NewServeMux()
	handleFunc := func(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
		handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))
		mux.Handle(pattern, handler)
	}
	handleFunc("POST /", CEP)
	handler := otelhttp.NewHandler(mux, "/")
	return handler
}

func main() {
	config.Load()
	config.SetupOTEL("processing-service")

	srv := &http.Server{
		Addr:        ":4000",
		BaseContext: func(_ net.Listener) context.Context { return config.Context },
		Handler:     HTTPHandler(),
	}
	println("Server is running on port 4000")
	serverError := make(chan error, 1)
	go func() {
		serverError <- srv.ListenAndServe()
	}()
	error := <-serverError
	println(error.Error())
	println("Server stopped.")
}

func CEP(w http.ResponseWriter, r *http.Request) {
	ctx, span := config.OTEL.Tracer.Start(r.Context(), "Second Cep Handler")
	defer span.End()
	config.OTEL.Logger.InfoContext(ctx, "2nd CEP request received", "method", r.Method, "url", r.URL.String())
	var input CepInput

	cepAttribute := attribute.String("cep", input.Cep)
	span.SetAttributes(cepAttribute)

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
}
