package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/black-dev-x/go-lab-cep-temperature/config"
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
	config.SetupOTEL("input-service")

	srv := &http.Server{
		Addr:        ":3000",
		Handler:     HTTPHandler(),
		BaseContext: func(_ net.Listener) context.Context { return config.Context },
	}
	println("Server is running on port 3000")
	serverError := make(chan error, 1)
	go func() {
		serverError <- srv.ListenAndServe()
	}()
	error := <-serverError
	println(error.Error())
	println("Server stopped.")
}

func CEP(w http.ResponseWriter, r *http.Request) {
	ctx, span := config.OTEL.Tracer.Start(r.Context(), "First Cep Handler")
	defer span.End()
	config.OTEL.Logger.InfoContext(ctx, "CEP request received", "method", r.Method, "url", r.URL.String())
	var input CepInput
	json.NewDecoder(r.Body).Decode(&input)
	cep := input.Cep
	time.Sleep(1 * time.Second)

	cepAttribute := attribute.String("cep", cep)
	span.SetAttributes(cepAttribute)

	length := len(cep)
	if length != 8 {
		http.Error(w, "invalid zipcode", 422)
		return
	}

	json, _ := json.Marshal(input)

	client := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	req, _ := http.NewRequest("POST", os.Getenv("PROCESSING_SERVICE_URL"), bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "failed to send request", 500)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
